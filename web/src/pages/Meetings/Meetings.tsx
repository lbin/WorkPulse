import React, { useEffect, useMemo, useState } from "react";
import { Button, Card, DatePicker, Form, Input, Modal, Space, Table, Tag } from "antd";
import type { ColumnsType } from "antd/es/table";
import { Link } from "react-router-dom";
import dayjs, { Dayjs } from "dayjs";
import { useTranslation } from "react-i18next";
import { createMeeting, listMeetings, type Meeting } from "../../api/meetings";

export default function MeetingsPage() {
  const { t } = useTranslation();
  const [range, setRange] = useState<[Dayjs, Dayjs] | null>(() => {
    const start = dayjs().startOf("week");
    const end = dayjs().endOf("week");
    return [start, end];
  });
  const [teamId, setTeamId] = useState<string | undefined>();
  const [data, setData] = useState<Meeting[]>([]);
  const [loading, setLoading] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);
  const [form] = Form.useForm();

  const query = useMemo(() => {
    return {
      start: range ? range[0].toISOString() : undefined,
      end: range ? range[1].toISOString() : undefined,
      team_id: teamId || undefined,
    };
  }, [range, teamId]);

  useEffect(() => {
    setLoading(true);
    listMeetings(query)
      .then((res) => setData(res.data.data))
      .finally(() => setLoading(false));
  }, [query]);

  const columns: ColumnsType<Meeting> = [
    {
      title: t("meetings.title"),
      dataIndex: "title",
      render: (text, record) => <Link to={`/meetings/${record.id}`}>{text}</Link>,
    },
    {
      title: t("meetings.time"),
      dataIndex: "scheduled_at",
      render: (value: string) => dayjs(value).format("YYYY-MM-DD HH:mm"),
    },
    {
      title: t("meetings.team"),
      dataIndex: "team_id",
      render: (value?: string) => value || t("common.not_set"),
    },
    {
      title: t("meetings.status"),
      dataIndex: "status",
      render: (value: string) => <Tag color="blue">{value}</Tag>,
    },
  ];

  const onCreate = async () => {
    const values = await form.validateFields();
    await createMeeting({
      title: values.title,
      scheduled_at: values.scheduled_at.toISOString(),
      duration_minutes: 60,
      team_id: values.team_id || undefined,
      agenda: values.agenda,
    });
    setCreateOpen(false);
    form.resetFields();
    listMeetings(query).then((res) => setData(res.data.data));
  };

  return (
    <Space direction="vertical" style={{ width: "100%" }} size="large">
      <Card
        title={t("meetings.filters")}
        extra={
          <Button type="primary" onClick={() => setCreateOpen(true)}>
            {t("meetings.create")}
          </Button>
        }
      >
        <Space wrap>
          <DatePicker.RangePicker value={range as any} onChange={(v) => setRange(v as [Dayjs, Dayjs] | null)} showTime allowClear />
          <Input
            style={{ width: 220 }}
            placeholder={t("meetings.team_placeholder")}
            value={teamId}
            onChange={(e) => setTeamId(e.target.value || undefined)}
          />
        </Space>
      </Card>

      <Card title={t("meetings.list")}>
        <Table rowKey="id" loading={loading} columns={columns} dataSource={data} pagination={false} />
      </Card>

      <Modal open={createOpen} onCancel={() => setCreateOpen(false)} onOk={onCreate} title={t("meetings.create")}>
        <Form form={form} layout="vertical">
          <Form.Item name="title" label={t("meetings.title")}
            rules={[{ required: true, message: t("common.required") }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="agenda" label={t("meetings.agenda")}>
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="scheduled_at" label={t("meetings.time")}
            rules={[{ required: true, message: t("common.required") }]}
          >
            <DatePicker showTime style={{ width: "100%" }} />
          </Form.Item>
          <Form.Item name="team_id" label={t("meetings.team")}>
            <Input placeholder={t("meetings.team_placeholder") as string} />
          </Form.Item>
        </Form>
      </Modal>
    </Space>
  );
}
