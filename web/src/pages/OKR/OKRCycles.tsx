import React, { useEffect, useMemo, useState } from "react";
import {
  Button,
  Card,
  Drawer,
  Flex,
  Progress,
  Select,
  Space,
  Table,
  Tag,
  Typography,
} from "antd";
import { useTranslation } from "react-i18next";
import { addLink, listObjectives, OKRKeyResult, removeLink } from "../../api/okr";

type EnrichedObjective = {
  id: string;
  title: string;
  status: string;
  cycle_id: string;
  team_id?: string;
  owner_user_id?: string;
  description?: string;
  keyResults: OKRKeyResult[];
};

const statusColors: Record<string, string> = {
  draft: "default",
  active: "processing",
  closed: "success",
  archived: "warning",
};

export default function Page() {
  const { t } = useTranslation();
  const [objectives, setObjectives] = useState<EnrichedObjective[]>([]);
  const [filters, setFilters] = useState<{ status?: string; cycle_id?: string }>(
    {}
  );
  const [loading, setLoading] = useState(false);
  const [selected, setSelected] = useState<EnrichedObjective | null>(null);

  useEffect(() => {
    setLoading(true);
    listObjectives(filters)
      .then(({ objectives, key_results }) => {
        const grouped: Record<string, OKRKeyResult[]> = {};
        key_results.forEach((kr) => {
          grouped[kr.objective_id] = grouped[kr.objective_id] || [];
          grouped[kr.objective_id].push(kr);
        });
        setObjectives(
          objectives.map((o) => ({
            ...o,
            keyResults: grouped[o.id] || [],
          }))
        );
      })
      .finally(() => setLoading(false));
  }, [filters]);

  const groupedByStatus = useMemo(() => {
    return objectives.reduce<Record<string, EnrichedObjective[]>>((acc, obj) => {
      acc[obj.status] = acc[obj.status] || [];
      acc[obj.status].push(obj);
      return acc;
    }, {});
  }, [objectives]);

  const columns = [
    { title: t("common.title"), dataIndex: "title" },
    {
      title: t("common.status"),
      dataIndex: "status",
      render: (value: string) => <Tag color={statusColors[value] || "default"}>{value}</Tag>,
    },
    {
      title: t("common.progress"),
      render: (_: unknown, record: EnrichedObjective) => {
        const progress = averageProgress(record.keyResults);
        return <Progress percent={Math.round(progress)} size="small" />;
      },
    },
  ];

  return (
    <Flex vertical gap={16}>
      <Card title={t("okr.filters.title", { defaultValue: "Filters" })}>
        <Space wrap>
          <Select
            placeholder={t("common.status")}
            style={{ width: 160 }}
            allowClear
            onChange={(v) => setFilters((f) => ({ ...f, status: v }))}
            options={[
              { label: "Active", value: "active" },
              { label: "Draft", value: "draft" },
              { label: "Closed", value: "closed" },
              { label: "Archived", value: "archived" },
            ]}
          />
          <Select
            placeholder={t("okr.filters.cycle", { defaultValue: "Cycle" })}
            style={{ width: 180 }}
            allowClear
            onChange={(v) => setFilters((f) => ({ ...f, cycle_id: v }))}
            options={objectives.map((o) => ({ label: o.cycle_id, value: o.cycle_id }))}
          />
        </Space>
      </Card>

      <Card title={t("okr.list", { defaultValue: "Objectives" })}>
        <Table
          loading={loading}
          columns={columns}
          dataSource={objectives}
          rowKey="id"
          onRow={(record) => ({
            onClick: () => setSelected(record),
            style: { cursor: "pointer" },
          })}
          pagination={false}
        />
      </Card>

      <Card title={t("okr.board", { defaultValue: "Kanban" })}>
        <Flex gap={12} wrap>
          {Object.entries(groupedByStatus).map(([status, items]) => (
            <Card key={status} title={<Tag color={statusColors[status] || "default"}>{status}</Tag>} style={{ minWidth: 240 }}>
              <Space direction="vertical" style={{ width: "100%" }}>
                {items.map((item) => (
                  <Card
                    key={item.id}
                    size="small"
                    onClick={() => setSelected(item)}
                    style={{ cursor: "pointer" }}
                  >
                    <Typography.Text strong>{item.title}</Typography.Text>
                    <Progress percent={Math.round(averageProgress(item.keyResults))} size="small" />
                  </Card>
                ))}
              </Space>
            </Card>
          ))}
        </Flex>
      </Card>

      <Drawer
        open={!!selected}
        onClose={() => setSelected(null)}
        width={420}
        title={selected?.title}
      >
        {selected && (
          <Flex vertical gap={12}>
            <Typography.Paragraph>{selected.description}</Typography.Paragraph>
            <Typography.Title level={5}>{t("okr.keyResults", { defaultValue: "Key Results" })}</Typography.Title>
            <Space direction="vertical" style={{ width: "100%" }}>
              {selected.keyResults.map((kr) => (
                <Card key={kr.id} size="small">
                  <Flex justify="space-between">
                    <div>
                      <Typography.Text strong>{kr.title}</Typography.Text>
                      <div>{t("okr.metric", { defaultValue: "Metric" })}: {kr.metric_type}</div>
                    </div>
                    <Progress percent={Math.round(progressValue(kr))} size="small" />
                  </Flex>
                  <Space style={{ marginTop: 8 }}>
                    <Button size="small" onClick={() => handleLink(selected.id, kr.id)}>Link</Button>
                    <Button size="small" danger onClick={() => handleUnlink(kr.id)}>Unlink</Button>
                  </Space>
                </Card>
              ))}
            </Space>
          </Flex>
        )}
      </Drawer>
    </Flex>
  );

  async function handleLink(objectiveId: string, keyResultId?: string) {
    await addLink({
      objective_id: objectiveId,
      key_result_id: keyResultId,
      entity_type: "project",
      entity_id: objectives[0]?.id || objectiveId,
      relation: "supports",
    });
  }

  async function handleUnlink(keyResultId: string) {
    // The in-memory backend responds without requiring a specific link id; here we simply call remove with the KR id.
    await removeLink(keyResultId);
  }
}

function averageProgress(krs: OKRKeyResult[]): number {
  if (!krs.length) return 0;
  const total = krs.reduce((sum, kr) => sum + progressValue(kr), 0);
  return total / krs.length;
}

function progressValue(kr: OKRKeyResult): number {
  if (!kr.target_value || kr.target_value === 0 || kr.current_value === undefined) return 0;
  return Math.min((kr.current_value / kr.target_value) * 100, 100);
}
