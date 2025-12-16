import React, { useEffect, useMemo, useState } from "react";
import {
  Button,
  Card,
  Drawer,
  Flex,
  Input,
  InputNumber,
  Progress,
  Select,
  Slider,
  Space,
  Statistic,
  Table,
  Tag,
  Typography,
  message,
} from "antd";
import { useTranslation } from "react-i18next";
import { addLink, fetchMetrics, listObjectives, OKRKeyResult, OKRMetrics, removeLink, updateKRProgress } from "../../api/okr";
import { useOrgUnits } from "../../components/AuthProvider";

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
  const units = useOrgUnits();
  const [objectives, setObjectives] = useState<EnrichedObjective[]>([]);
  const [metrics, setMetrics] = useState<OKRMetrics | null>(null);
  const [filters, setFilters] = useState<{ status?: string; cycle_id?: string; team_id?: string; owner_user_id?: string }>(
    {}
  );
  const [loading, setLoading] = useState(false);
  const [metricsLoading, setMetricsLoading] = useState(false);
  const [selected, setSelected] = useState<EnrichedObjective | null>(null);
  const [progressDrafts, setProgressDrafts] = useState<Record<string, { current?: number; confidence?: number }>>({});

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

  useEffect(() => {
    setMetricsLoading(true);
    fetchMetrics(filters)
      .then(setMetrics)
      .finally(() => setMetricsLoading(false));
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
          <Select
            placeholder={t("okr.filters.team", { defaultValue: "Team" })}
            style={{ width: 200 }}
            allowClear
            onChange={(v) => setFilters((f) => ({ ...f, team_id: v }))}
            options={units.map((u) => ({ label: u.name, value: u.id }))}
          />
          <Input
            placeholder={t("okr.filters.owner", { defaultValue: "Owner user id" })}
            style={{ width: 200 }}
            allowClear
            onChange={(e) => setFilters((f) => ({ ...f, owner_user_id: e.target.value || undefined }))}
          />
          <Button onClick={() => setFilters({})}>{t("common.reset", { defaultValue: "Reset" })}</Button>
        </Space>
      </Card>

      <Card title={t("okr.overview", { defaultValue: "OKR Health" })} loading={metricsLoading}>
        <Space size={32} wrap>
          <Statistic
            title={t("okr.completion", { defaultValue: "Avg. completion" })}
            value={metrics?.completion || 0}
            precision={1}
            suffix="%"
          />
          <Statistic
            title={t("okr.deviation", { defaultValue: "Avg. deviation" })}
            value={metrics?.deviation || 0}
            precision={1}
            suffix="%"
          />
          <div>
            <Typography.Text type="secondary">
              {t("okr.progress", { defaultValue: "Overall progress" })}
            </Typography.Text>
            <Progress
              style={{ width: 260 }}
              percent={Math.round(metrics?.completion || 0)}
              status={(metrics?.completion || 0) >= 80 ? "success" : "active"}
            />
          </div>
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
                              <div>
                                {t("okr.target", { defaultValue: "Target" })}: {kr.target_value ?? t("common.not_set")}
                              </div>
                            </div>
                            <Progress percent={Math.round(progressValue(kr))} size="small" />
                          </Flex>
                          <Space direction="vertical" style={{ marginTop: 8, width: "100%" }}>
                            <Space>
                              <InputNumber
                                placeholder={t("okr.current", { defaultValue: "Current" }) as string}
                                value={progressDrafts[kr.id]?.current ?? kr.current_value}
                                onChange={(val) =>
                                  setProgressDrafts((drafts) => ({
                                    ...drafts,
                                    [kr.id]: { ...drafts[kr.id], current: val ?? undefined },
                                  }))
                                }
                              />
                              <div style={{ flex: 1 }}>
                                <Typography.Text type="secondary">
                                  {t("okr.confidence", { defaultValue: "Confidence" })}
                                </Typography.Text>
                                <Slider
                                  min={0}
                                  max={100}
                                  value={progressDrafts[kr.id]?.confidence ?? kr.confidence ?? 0}
                                  onChange={(val) =>
                                    setProgressDrafts((drafts) => ({
                                      ...drafts,
                                      [kr.id]: { ...drafts[kr.id], confidence: Array.isArray(val) ? val[0] : val },
                                    }))
                                  }
                                />
                              </div>
                            </Space>
                            <Space>
                              <Button size="small" type="primary" onClick={() => handleProgressSave(kr)}>
                                {t("common.save")}
                              </Button>
                              <Button size="small" onClick={() => handleLink(selected.id, kr.id)}>Link</Button>
                              <Button size="small" danger onClick={() => handleUnlink(kr.id)}>Unlink</Button>
                            </Space>
                          </Space>
                        </Card>
                      ))}
                    </Space>
          </Flex>
        )}
      </Drawer>
    </Flex>
  );

  async function handleProgressSave(kr: OKRKeyResult) {
    const draft = progressDrafts[kr.id] || {};
    const current = draft.current ?? kr.current_value ?? 0;
    const confidence = draft.confidence ?? kr.confidence ?? 0;

    await updateKRProgress(kr.id, current, confidence);
    message.success(t("common.saved"));

    setObjectives((prev) =>
      prev.map((obj) => ({
        ...obj,
        keyResults: obj.keyResults.map((item) =>
          item.id === kr.id ? { ...item, current_value: current, confidence } : item
        ),
      }))
    );

    setSelected((prev) =>
      prev
        ? {
            ...prev,
            keyResults: prev.keyResults.map((item) =>
              item.id === kr.id ? { ...item, current_value: current, confidence } : item
            ),
          }
        : prev
    );

    setProgressDrafts((drafts) => ({ ...drafts, [kr.id]: { current, confidence } }));
  }

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
