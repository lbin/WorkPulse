import React, { useEffect, useMemo, useState } from "react";
import dayjs, { Dayjs } from "dayjs";
import {
  Button,
  Card,
  DatePicker,
  Drawer,
  Flex,
  Form,
  Input,
  InputNumber,
  Modal,
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
import {
  addLink,
  createCycle,
  createKeyResult,
  createObjective,
  fetchMetrics,
  listCycles,
  listObjectives,
  OKRCycle,
  OKRKeyResult,
  OKRMetrics,
  removeLink,
  updateKRProgress,
} from "../../api/okr";
import { useOrgUnits } from "../../components/AuthProvider";

const { RangePicker } = DatePicker;
const { TextArea } = Input;

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
  const [cycles, setCycles] = useState<OKRCycle[]>([]);
  const [objectives, setObjectives] = useState<EnrichedObjective[]>([]);
  const [metrics, setMetrics] = useState<OKRMetrics | null>(null);
  const [filters, setFilters] = useState<{ status?: string; cycle_id?: string; team_id?: string; owner_user_id?: string }>(
    {}
  );
  const [loading, setLoading] = useState(false);
  const [metricsLoading, setMetricsLoading] = useState(false);
  const [selected, setSelected] = useState<EnrichedObjective | null>(null);
  const [progressDrafts, setProgressDrafts] = useState<Record<string, { current?: number; confidence?: number }>>({});
  const [cycleModalOpen, setCycleModalOpen] = useState(false);
  const [objectiveModalOpen, setObjectiveModalOpen] = useState(false);
  const [krModalOpen, setKrModalOpen] = useState(false);
  const [cycleForm] = Form.useForm();
  const [objectiveForm] = Form.useForm();
  const [krForm] = Form.useForm();

  useEffect(() => {
    loadCycles();
  }, []);

  useEffect(() => {
    refreshObjectives();
  }, [filters]);

  useEffect(() => {
    refreshMetrics();
  }, [filters]);

  const groupedByStatus = useMemo(() => {
    return objectives.reduce<Record<string, EnrichedObjective[]>>((acc, obj) => {
      acc[obj.status] = acc[obj.status] || [];
      acc[obj.status].push(obj);
      return acc;
    }, {});
  }, [objectives]);

  const objectiveOptions = useMemo(
    () => objectives.map((o) => ({ label: o.title, value: o.id })),
    [objectives]
  );

  const cycleLookup = useMemo(() => {
    const map = new Map<string, OKRCycle>();
    cycles.forEach((c) => map.set(c.id, c));
    return map;
  }, [cycles]);

  const columns = [
    { title: t("common.title"), dataIndex: "title" },
    {
      title: t("okr.filters.cycle", { defaultValue: "Cycle" }),
      dataIndex: "cycle_id",
      render: (value: string) => cycleLookup.get(value)?.name || value,
    },
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

  const cycleColumns = [
    { title: t("common.title"), dataIndex: "name" },
    { title: t("common.type", { defaultValue: "Type" }), dataIndex: "type" },
    {
      title: t("common.status"),
      dataIndex: "status",
      render: (v: string) => <Tag color={statusColors[v] || "default"}>{v}</Tag>,
    },
    {
      title: t("okr.filters.team", { defaultValue: "Team" }),
      dataIndex: "team_id",
      render: (v?: string) => units.find((u) => u.id === v)?.name || v || "-",
    },
    {
      title: t("okr.filters.cycle", { defaultValue: "Timeline" }),
      render: (_: unknown, record: OKRCycle) => (
        <Typography.Text type="secondary">
          {record.start_date || t("common.not_set")} – {record.end_date || t("common.not_set")}
        </Typography.Text>
      ),
    },
  ];

  return (
    <Flex vertical gap={16}>
      <Card
        title={t("okr.filters.title", { defaultValue: "Filters" })}
        extra={
          <Space>
            <Button onClick={() => setCycleModalOpen(true)}>{t("okr.newCycle", { defaultValue: "New cycle" })}</Button>
            <Button type="primary" onClick={() => setObjectiveModalOpen(true)}>
              {t("okr.newObjective", { defaultValue: "New objective" })}
            </Button>
            <Button onClick={() => setKrModalOpen(true)}>{t("okr.newKR", { defaultValue: "New key result" })}</Button>
          </Space>
        }
      >
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
            options={cycles.map((cycle) => ({ label: cycle.name, value: cycle.id }))}
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

      <Card title={t("okr.filters.cycle", { defaultValue: "Cycles" })}>
        <Table
          loading={loading}
          columns={cycleColumns}
          dataSource={cycles}
          rowKey="id"
          pagination={false}
        />
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

      <Drawer open={!!selected} onClose={() => setSelected(null)} width={420} title={selected?.title}>
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

      <Modal
        title={t("okr.newCycle", { defaultValue: "New cycle" })}
        open={cycleModalOpen}
        onCancel={() => setCycleModalOpen(false)}
        onOk={handleCreateCycle}
        okText={t("common.save")}
        destroyOnClose
      >
        <Form form={cycleForm} layout="vertical">
          <Form.Item name="name" label={t("common.title") as string} rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="type" label={t("common.type") as string}>
            <Select
              options={[
                { value: "quarter", label: "Quarter" },
                { value: "year", label: "Year" },
                { value: "month", label: "Month" },
              ]}
            />
          </Form.Item>
          <Form.Item name="status" label={t("common.status") as string}>
            <Select
              options={[
                { label: "Active", value: "active" },
                { label: "Draft", value: "draft" },
                { label: "Closed", value: "closed" },
              ]}
            />
          </Form.Item>
          <Form.Item name="dates" label={t("okr.filters.cycle", { defaultValue: "Timeline" }) as string}>
            <RangePicker style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item name="team_id" label={t("okr.filters.team", { defaultValue: "Team" }) as string}>
            <Select allowClear options={units.map((u) => ({ label: u.name, value: u.id }))} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={t("okr.newObjective", { defaultValue: "New objective" })}
        open={objectiveModalOpen}
        onCancel={() => setObjectiveModalOpen(false)}
        onOk={handleCreateObjective}
        okText={t("common.save")}
        destroyOnClose
      >
        <Form form={objectiveForm} layout="vertical">
          <Form.Item name="cycle_id" label={t("okr.filters.cycle", { defaultValue: "Cycle" }) as string} rules={[{ required: true }]}>
            <Select options={cycles.map((cycle) => ({ label: cycle.name, value: cycle.id }))} />
          </Form.Item>
          <Form.Item name="title" label={t("common.title") as string} rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label={t("common.description", { defaultValue: "Description" }) as string}>
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="team_id" label={t("okr.filters.team", { defaultValue: "Team" }) as string}>
            <Select allowClear options={units.map((u) => ({ label: u.name, value: u.id }))} />
          </Form.Item>
          <Form.Item name="owner_user_id" label={t("okr.filters.owner", { defaultValue: "Owner user id" }) as string}>
            <Input />
          </Form.Item>
          <Form.Item name="status" label={t("common.status") as string}>
            <Select
              options={[
                { label: "Active", value: "active" },
                { label: "Draft", value: "draft" },
                { label: "Closed", value: "closed" },
              ]}
            />
          </Form.Item>
          <Form.Item name="tags" label={t("common.tags", { defaultValue: "Tags" }) as string}>
            <Select mode="tags" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={t("okr.newKR", { defaultValue: "New key result" })}
        open={krModalOpen}
        onCancel={() => setKrModalOpen(false)}
        onOk={handleCreateKeyResult}
        okText={t("common.save")}
        destroyOnClose
      >
        <Form form={krForm} layout="vertical">
          <Form.Item name="objective_id" label={t("okr.keyResult", { defaultValue: "Key Result" }) as string} rules={[{ required: true }]}>
            <Select options={objectiveOptions} />
          </Form.Item>
          <Form.Item name="title" label={t("common.title") as string} rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="metric_type" label={t("okr.metric", { defaultValue: "Metric" }) as string}>
            <Select
              options={[
                { label: "Number", value: "number" },
                { label: "Percent", value: "percent" },
                { label: "Currency", value: "currency" },
              ]}
            />
          </Form.Item>
          <Form.Item name="target_value" label={t("okr.target", { defaultValue: "Target" }) as string}>
            <InputNumber style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item name="unit" label={t("common.unit", { defaultValue: "Unit" }) as string}>
            <Input />
          </Form.Item>
          <Form.Item name="status" label={t("common.status") as string}>
            <Select
              options={[
                { label: "Active", value: "active" },
                { label: "Draft", value: "draft" },
                { label: "Closed", value: "closed" },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </Flex>
  );

  async function loadCycles() {
    const data = await listCycles();
    setCycles(data);
  }

  async function refreshObjectives() {
    setLoading(true);
    try {
      const { objectives, key_results } = await listObjectives(filters);
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
    } finally {
      setLoading(false);
    }
  }

  async function refreshMetrics() {
    setMetricsLoading(true);
    try {
      const m = await fetchMetrics(filters);
      setMetrics(m);
    } finally {
      setMetricsLoading(false);
    }
  }

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
    await removeLink(keyResultId);
  }

  async function handleCreateCycle() {
    const values = await cycleForm.validateFields();
    const dates = values.dates as Dayjs[] | undefined;
    await createCycle({
      name: values.name,
      type: values.type,
      status: values.status,
      team_id: values.team_id,
      start_date: dates?.[0]?.format("YYYY-MM-DD"),
      end_date: dates?.[1]?.format("YYYY-MM-DD"),
    });
    message.success(t("common.saved"));
    setCycleModalOpen(false);
    cycleForm.resetFields();
    await loadCycles();
  }

  async function handleCreateObjective() {
    const values = await objectiveForm.validateFields();
    await createObjective(values);
    message.success(t("common.saved"));
    setObjectiveModalOpen(false);
    objectiveForm.resetFields();
    await refreshObjectives();
    await refreshMetrics();
  }

  async function handleCreateKeyResult() {
    const values = await krForm.validateFields();
    await createKeyResult(values);
    message.success(t("common.saved"));
    setKrModalOpen(false);
    krForm.resetFields();
    await refreshObjectives();
    await refreshMetrics();
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
