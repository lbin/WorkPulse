import React, { useEffect, useMemo, useState } from "react";
import {
  Button,
  Card,
  DatePicker,
  Divider,
  Flex,
  Form,
  Input,
  Select,
  Space,
  Tag,
  Typography,
  message,
} from "antd";
import dayjs, { Dayjs } from "dayjs";
import { useNavigate, useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  Report,
  changeReportStatus,
  createReport,
  getReport,
  updateReport,
} from "../../api/reports";

const typeOptions = ["daily", "weekly", "monthly", "quarterly", "halfyearly", "yearly"];
const statusColor: Record<string, string> = {
  draft: "default",
  submitted: "processing",
  approved: "success",
  rejected: "warning",
  archived: "purple",
};

type StructuredField = { label?: string; value?: string };
type LinkField = { target_type?: string; target_id?: string; relation?: string };

type FormShape = {
  title?: string;
  type: string;
  period: [Dayjs, Dayjs];
  summary?: string;
  content?: { rich_text?: string; fields?: StructuredField[] };
  links?: LinkField[];
};

export default function Page() {
  const { id } = useParams();
  const isNew = !id || id === "new";
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [form] = Form.useForm<FormShape>();
  const [loading, setLoading] = useState(true);
  const [report, setReport] = useState<Report | null>(null);

  useEffect(() => {
    if (isNew) {
      seedDefaults();
      setLoading(false);
      return;
    }
    getReport(id as string)
      .then((res) => {
        const payload = normalizeContent(res.report.content);
        setReport({ ...res.report, links: res.links });
        form.setFieldsValue({
          title: res.report.title || undefined,
          type: res.report.type,
          summary: res.report.summary || undefined,
          period: [dayjs(res.report.period_start), dayjs(res.report.period_end)],
          content: payload,
          links: res.links?.map((l) => ({
            target_type: l.target_type,
            target_id: l.target_id,
            relation: l.relation,
          })),
        });
      })
      .catch(() => message.error(t("reports.loadFailed")))
      .finally(() => setLoading(false));
  }, [id]);

  const statusTag = useMemo(() => {
    if (!report) return null;
    return <Tag color={statusColor[report.status] || "default"}>{report.status}</Tag>;
  }, [report]);

  return (
    <Flex vertical gap={16}>
      <Card title={isNew ? t("reports.new") : t("reports.editReport")}> 
        <Form layout="vertical" form={form} initialValues={defaultValues()}>
          <Flex gap={16} wrap>
            <Form.Item name="title" label={t("reports.titleLabel")} style={{ flex: 1 }}>
              <Input placeholder={t("reports.titlePlaceholder") as string} />
            </Form.Item>
            <Form.Item name="type" label={t("reports.type")}> 
              <Select options={typeOptions.map((v) => ({ label: v, value: v }))} style={{ width: 200 }} />
            </Form.Item>
            <Form.Item name="period" label={t("reports.period")}> 
              <DatePicker.RangePicker />
            </Form.Item>
            <Form.Item label={t("reports.currentStatus")}>
              {statusTag || <Tag>{t("reports.draft")}</Tag>}
            </Form.Item>
          </Flex>

          <Form.Item name="summary" label={t("reports.summary")}> 
            <Input.TextArea rows={3} placeholder={t("reports.summaryPlaceholder") as string} />
          </Form.Item>

          <Divider orientation="left">{t("reports.richText")}</Divider>
          <Form.Item name={["content", "rich_text"]}>
            <Input.TextArea rows={6} placeholder={t("reports.richTextPlaceholder") as string} />
          </Form.Item>

          <Divider orientation="left">{t("reports.structuredFields")}</Divider>
          <Form.List name={["content", "fields"]}>
            {(fields, { add, remove }) => (
              <Flex vertical gap={8}>
                {fields.map((field) => (
                  <Flex gap={8} key={field.key} align="baseline" wrap>
                    <Form.Item
                      {...field}
                      name={[field.name, "label"]}
                      style={{ flex: 1 }}
                    >
                      <Input placeholder={t("reports.fieldName") as string} />
                    </Form.Item>
                    <Form.Item
                      {...field}
                      name={[field.name, "value"]}
                      style={{ flex: 2 }}
                    >
                      <Input placeholder={t("reports.fieldValue") as string} />
                    </Form.Item>
                    <Button danger onClick={() => remove(field.name)}>
                      {t("reports.delete")}
                    </Button>
                  </Flex>
                ))}
                <Button type="dashed" onClick={() => add()}>
                  {t("reports.addField")}
                </Button>
              </Flex>
            )}
          </Form.List>

          <Divider orientation="left">{t("reports.links")}</Divider>
          <Form.List name="links">
            {(fields, { add, remove }) => (
              <Flex vertical gap={8}>
                {fields.map((field) => (
                  <Flex key={field.key} gap={8} wrap align="baseline">
                    <Form.Item
                      {...field}
                      name={[field.name, "target_type"]}
                      style={{ width: 180 }}
                    >
                      <Select
                        placeholder={t("reports.targetType") as string}
                        options={[
                          { label: "Objective", value: "okr_objective" },
                          { label: "Key Result", value: "okr_key_result" },
                          { label: "Project", value: "project" },
                          { label: "Meeting", value: "meeting" },
                        ]}
                      />
                    </Form.Item>
                    <Form.Item
                      {...field}
                      name={[field.name, "target_id"]}
                      style={{ minWidth: 220 }}
                    >
                      <Input placeholder={t("reports.targetId") as string} />
                    </Form.Item>
                    <Form.Item
                      {...field}
                      name={[field.name, "relation"]}
                      style={{ width: 160 }}
                    >
                      <Input placeholder={t("reports.relation") as string} />
                    </Form.Item>
                    <Button danger onClick={() => remove(field.name)}>
                      {t("reports.delete")}
                    </Button>
                  </Flex>
                ))}
                <Button type="dashed" onClick={() => add()}>
                  {t("reports.addLink")}
                </Button>
              </Flex>
            )}
          </Form.List>

          <Divider />
          <Space>
            <Button type="primary" onClick={() => handleSave(false)} loading={loading}>
              {t("reports.saveDraft")}
            </Button>
            <Button onClick={() => handleSave(true)} loading={loading}>
              {t("reports.submit")}
            </Button>
            {!isNew && report?.status === "submitted" && (
              <>
                <Button type="primary" onClick={() => handleReview("approve")}>
                  {t("reports.approve")}
                </Button>
                <Button danger onClick={() => handleReview("reject")}>
                  {t("reports.reject")}
                </Button>
              </>
            )}
          </Space>
        </Form>
      </Card>

      {!isNew && report?.links && (
        <Card title={t("reports.associations")}> 
          <Space direction="vertical" style={{ width: "100%" }}>
            {report.links.map((l) => (
              <Flex key={`${l.target_type}-${l.target_id}`} justify="space-between">
                <div>
                  <Typography.Text strong>{l.target_type}</Typography.Text>
                  <div>{l.target_id}</div>
                </div>
                <Tag>{l.relation || "related"}</Tag>
              </Flex>
            ))}
          </Space>
        </Card>
      )}
    </Flex>
  );

  async function handleSave(shouldSubmit: boolean) {
    try {
      setLoading(true);
      const values = await form.validateFields();
      const payload = serializePayload(values);
      if (isNew) {
        const res = await createReport(payload);
        setReport(res.report);
        message.success(t("reports.saved"));
        if (shouldSubmit) {
          await changeReportStatus(res.report.id, "submit");
          message.success(t("reports.submitSuccess"));
        }
        navigate(`/reports/${res.report.id}/edit`);
        const refreshed = await getReport(res.report.id);
        setReport({ ...refreshed.report, links: refreshed.links });
      } else {
        await updateReport(id as string, payload);
        message.success(t("reports.saved"));
        const targetId = id as string;
        if (shouldSubmit) {
          await changeReportStatus(targetId, "submit");
          message.success(t("reports.submitSuccess"));
        }
        const refreshed = await getReport(targetId);
        setReport({ ...refreshed.report, links: refreshed.links });
      }
    } catch (err) {
      message.error(t("reports.saveFailed"));
    } finally {
      setLoading(false);
    }
  }

  async function handleReview(action: "approve" | "reject") {
    try {
      await changeReportStatus(id as string, action);
      message.success(t("reports.reviewUpdated"));
      const refreshed = await getReport(id as string);
      setReport({ ...refreshed.report, links: refreshed.links });
    } catch (err) {
      message.error(t("reports.saveFailed"));
    }
  }

  function serializePayload(values: FormShape) {
    return {
      title: values.title,
      type: values.type,
      period_start: values.period?.[0]?.format("YYYY-MM-DD"),
      period_end: values.period?.[1]?.format("YYYY-MM-DD"),
      summary: values.summary,
      content: values.content || { rich_text: "", fields: [] },
      links:
        values.links
          ?.filter((l) => l.target_type && l.target_id)
          .map((l) => ({
            target_type: l.target_type,
            target_id: l.target_id,
            relation: l.relation,
          })) || [],
    };
  }

  function normalizeContent(content: any): FormShape["content"] {
    if (!content) return { rich_text: "", fields: [] };
    if (typeof content === "string") {
      try {
        return JSON.parse(content);
      } catch (err) {
        return { rich_text: content, fields: [] };
      }
    }
    return content;
  }

  function defaultValues(): FormShape {
    const today = dayjs();
    return {
      type: "weekly",
      period: [today.startOf("week"), today.endOf("week")],
      content: { rich_text: "", fields: [] },
      links: [],
    } as FormShape;
  }

  function seedDefaults() {
    form.setFieldsValue(defaultValues());
  }
}
