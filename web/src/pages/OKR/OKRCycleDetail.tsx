import React, { useEffect, useMemo, useState } from "react";
import { Card, Flex, Progress, Space, Statistic, Table, Tag, Typography } from "antd";
import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { fetchMetrics, OKRCoverage, OKRMetrics } from "../../api/okr";

const entityLabels: Record<string, string> = {
  report: "Reports",
  project: "Projects",
  meeting: "Meetings",
};

export default function OKRCycleDetail() {
  const { id } = useParams<{ id: string }>();
  const { t } = useTranslation();
  const [metrics, setMetrics] = useState<OKRMetrics | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    fetchMetrics({ cycle_id: id })
      .then(setMetrics)
      .finally(() => setLoading(false));
  }, [id]);

  const coverageData = useMemo(() => metrics?.coverage || [], [metrics]);

  return (
    <Flex vertical gap={16}>
      <Card
        loading={loading}
        title={t("okr.overview", { defaultValue: "OKR Health" })}
        extra={id && <Tag color="blue">Cycle: {id}</Tag>}
      >
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
              percent={Number((metrics?.completion || 0).toFixed(1))}
              status={(metrics?.completion || 0) >= 80 ? "success" : "active"}
            />
          </div>
        </Space>
      </Card>

      <Card
        loading={loading}
        title={t("okr.coverage", { defaultValue: "Coverage by linked work" })}
      >
        <Table<OKRCoverage>
          rowKey={(row) => row.entity_type}
          pagination={false}
          dataSource={coverageData}
          columns={[
            {
              title: t("common.type", { defaultValue: "Type" }),
              dataIndex: "entity_type",
              render: (v: string) => entityLabels[v] || v,
            },
            {
              title: t("okr.coverageRate", { defaultValue: "Coverage" }),
              render: (_, row) => (
                <Progress percent={Number(row.coverage.toFixed(1))} size="small" />
              ),
            },
            {
              title: t("okr.deviation", { defaultValue: "Avg. deviation" }),
              render: (_, row) => `${row.deviation.toFixed(1)}%`,
            },
            {
              title: t("okr.linked", { defaultValue: "Linked KRs" }),
              render: (_, row) => (
                <Typography.Text>
                  {row.linked_key_results} / {row.total_key_results}
                </Typography.Text>
              ),
            },
          ]}
        />
        {!coverageData.length && (
          <Typography.Text type="secondary">
            {t("okr.noCoverage", { defaultValue: "No linked work items found." })}
          </Typography.Text>
        )}
      </Card>

      <Card
        loading={loading}
        title={t("okr.unlinked", { defaultValue: "Key results without links" })}
      >
        <Table
          rowKey="id"
          pagination={false}
          dataSource={metrics?.unlinked_key_results || []}
          columns={[
            { title: t("common.title"), dataIndex: "title" },
            {
              title: t("okr.progress", { defaultValue: "Progress" }),
              render: () => <Tag color="default">0%</Tag>,
            },
          ]}
        />
        {!metrics?.unlinked_key_results?.length && (
          <Typography.Text type="secondary">
            {t("okr.allLinked", { defaultValue: "All key results are linked." })}
          </Typography.Text>
        )}
      </Card>
    </Flex>
  );
}
