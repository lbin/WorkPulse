import React, { useEffect, useMemo, useState } from "react";
import {
  Button,
  Card,
  DatePicker,
  Drawer,
  Flex,
  Form,
  Input,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from "antd";
import { Link, useNavigate } from "react-router-dom";
import dayjs, { Dayjs } from "dayjs";
import { useTranslation } from "react-i18next";
import { Report, ReportLink, exportReports, getReport, listReports } from "../../api/reports";
import { usePermission } from "../../components/AuthProvider";

const statusColor: Record<string, string> = {
  draft: "default",
  submitted: "processing",
  approved: "success",
  rejected: "warning",
  archived: "purple",
};

export default function Page() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [items, setItems] = useState<Report[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [selectedLinks, setSelectedLinks] = useState<ReportLink[]>([]);
  const canManage = usePermission("reports.manage");

  useEffect(() => {
    load();
  }, []);

  const columns = useMemo(
    () => [
      { title: t("reports.type"), dataIndex: "type" },
      {
        title: t("reports.period"),
        render: (_: unknown, record: Report) => `${record.period_start} ~ ${record.period_end}`,
      },
      {
        title: t("common.status"),
        dataIndex: "status",
        render: (value: string) => <Tag color={statusColor[value] || "default"}>{value}</Tag>,
      },
      { title: t("reports.owner"), dataIndex: "author_user_id" },
      {
        title: t("reports.lastUpdated"),
        dataIndex: "updated_at",
        render: (value: string) => (value ? dayjs(value).format("YYYY-MM-DD") : "-"),
      },
      {
        title: t("common.actions", { defaultValue: "Actions" }),
        render: (_: unknown, record: Report) => (
          <Space>
            {canManage && <Link to={`/reports/${record.id}/edit`}>{t("reports.edit")}</Link>}
            <Button type="link" onClick={() => openDetail(record.id)}>
              {t("reports.viewDetail")}
            </Button>
          </Space>
        ),
      },
    ],
    [t, canManage]
  );

  return (
    <Flex vertical gap={16}>
      <Card title={t("reports.filters")}> 
        <Form form={form} layout="inline" onFinish={load} initialValues={{ status: undefined, type: undefined }}>
          <Form.Item name="type" label={t("reports.type")}>
            <Select
              style={{ width: 140 }}
              allowClear
              options={["daily", "weekly", "monthly", "quarterly", "halfyearly", "yearly"].map((v) => ({
                label: v,
                value: v,
              }))}
            />
          </Form.Item>
          <Form.Item name="status" label={t("common.status")}>
            <Select
              style={{ width: 140 }}
              allowClear
              options={["draft", "submitted", "approved", "rejected"].map((v) => ({ label: v, value: v }))}
            />
          </Form.Item>
          <Form.Item name="author_user_id" label={t("reports.owner")}> 
            <Input placeholder={t("reports.ownerPlaceholder") as string} style={{ width: 200 }} />
          </Form.Item>
          <Form.Item name="period" label={t("reports.period")}> 
            <DatePicker.RangePicker allowEmpty={[true, true]} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {t("reports.search")}
              </Button>
              <Button onClick={() => form.resetFields()}>{t("reports.reset")}</Button>
            </Space>
          </Form.Item>
          <Form.Item style={{ marginLeft: "auto" }}>
            <Space>
              <Button onClick={handleExport}>{t("reports.export")}</Button>
              {canManage && (
                <Button type="primary" onClick={() => navigate("/reports/new/edit")}>
                  {t("reports.new")}
                </Button>
              )}
            </Space>
          </Form.Item>
        </Form>
      </Card>

      <Card title={t("reports.listTitle")}> 
        <Table
          rowKey="id"
          loading={loading}
          dataSource={items}
          columns={columns as any}
          pagination={{
            total,
            onChange: (page, pageSize) => load({ page: page - 1, page_size: pageSize }),
          }}
        />
      </Card>

      <Drawer open={!!selectedId} onClose={() => setSelectedId(null)} width={520} title={t("reports.detail")}> 
        {selectedId && (
          <Flex vertical gap={12}>
            {renderDetail(selectedId)}
            <Card size="small" title={t("reports.links")}> 
              {selectedLinks.length === 0 && (
                <Typography.Text type="secondary">{t("reports.noLinks")}</Typography.Text>
              )}
              <Space direction="vertical" style={{ width: "100%" }}>
                {selectedLinks.map((link) => (
                  <Flex key={`${link.target_type}-${link.target_id}`} justify="space-between">
                    <div>
                      <Typography.Text strong>{link.target_type}</Typography.Text>
                      <div>{link.target_id}</div>
                    </div>
                    <Tag>{link.relation || "related"}</Tag>
                  </Flex>
                ))}
              </Space>
            </Card>
          </Flex>
        )}
      </Drawer>
    </Flex>
  );

  async function load(pagination?: { page?: number; page_size?: number }) {
    const values = form.getFieldsValue();
    const params: any = {
      type: values.type,
      status: values.status,
      author_user_id: values.author_user_id || undefined,
      page: pagination?.page ?? 0,
      page_size: pagination?.page_size ?? 10,
    };
    const period = values.period as [Dayjs, Dayjs] | undefined;
    if (period?.[0]) params.from = period[0].format("YYYY-MM-DD");
    if (period?.[1]) params.to = period[1].format("YYYY-MM-DD");

    setLoading(true);
    try {
      const res = await listReports(params);
      setItems(res.items);
      setTotal(res.total);
    } catch (err) {
      message.error(t("reports.loadFailed"));
    } finally {
      setLoading(false);
    }
  }

  async function openDetail(id: string) {
    setSelectedId(id);
    try {
      const res = await getReport(id);
      setSelectedLinks(res.links || []);
    } catch (err) {
      message.error(t("reports.loadFailed"));
    }
  }

  async function handleExport() {
    try {
      const values = form.getFieldsValue();
      const params: any = {
        type: values.type,
        status: values.status,
        author_user_id: values.author_user_id || undefined,
      };
      const period = values.period as [Dayjs, Dayjs] | undefined;
      if (period?.[0]) params.from = period[0].format("YYYY-MM-DD");
      if (period?.[1]) params.to = period[1].format("YYYY-MM-DD");
      const blob = await exportReports(params);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "reports.csv";
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      message.error(t("reports.exportFailed"));
    }
  }

  function renderDetail(id: string) {
    const report = items.find((r) => r.id === id);
    if (!report) return null;
    return (
      <Card size="small" title={report.title || report.type}>
        <Space direction="vertical">
          <Typography.Text>{report.summary}</Typography.Text>
          <Typography.Text type="secondary">
            {report.period_start} ~ {report.period_end}
          </Typography.Text>
          <Tag color={statusColor[report.status] || "default"}>{report.status}</Tag>
        </Space>
      </Card>
    );
  }
}
