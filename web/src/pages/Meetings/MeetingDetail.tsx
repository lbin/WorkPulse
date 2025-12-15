import React, { useEffect, useState } from "react";
import { Button, Card, Col, Descriptions, Form, Input, Row, Space, Table, Tag, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import dayjs from "dayjs";
import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  addMeetingAction,
  addMeetingLink,
  assignMeetingAction,
  convertMeetingAction,
  getMeeting,
  updateMeeting,
  type Meeting,
  type MeetingAction,
  type MeetingLink,
} from "../../api/meetings";

export default function MeetingDetail() {
  const { t } = useTranslation();
  const { id } = useParams();
  const [meeting, setMeeting] = useState<Meeting>();
  const [actions, setActions] = useState<MeetingAction[]>([]);
  const [links, setLinks] = useState<MeetingLink[]>([]);
  const [loading, setLoading] = useState(false);
  const [notesForm] = Form.useForm();
  const [actionForm] = Form.useForm();
  const [linkForm] = Form.useForm();

  const fetchData = () => {
    if (!id) return;
    setLoading(true);
    getMeeting(id)
      .then((res) => {
        setMeeting(res.data.data.meeting);
        setActions(res.data.data.actions);
        setLinks(res.data.data.links);
        notesForm.setFieldsValue({
          notes: res.data.data.meeting?.notes,
          agenda: res.data.data.meeting?.agenda,
          status: res.data.data.meeting?.status,
        });
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchData();
  }, [id]);

  const saveNotes = async () => {
    if (!id) return;
    const values = await notesForm.validateFields();
    await updateMeeting(id, { notes: values.notes, agenda: values.agenda, status: values.status });
    message.success(t("common.saved"));
    fetchData();
  };

  const addAction = async () => {
    if (!id) return;
    const values = await actionForm.validateFields();
    await addMeetingAction(id, {
      title: values.title,
      owner_user_id: values.owner_user_id || undefined,
      due_date: values.due_date || undefined,
    });
    actionForm.resetFields();
    fetchData();
  };

  const updateActionStatus = async (actionId: string, status: string) => {
    await assignMeetingAction(actionId, { status });
    fetchData();
  };

  const convertAction = async (actionId: string) => {
    await convertMeetingAction(actionId);
    message.success(t("meetings.converted"));
    fetchData();
  };

  const addLinkAction = async () => {
    if (!id) return;
    const values = await linkForm.validateFields();
    await addMeetingLink(id, values);
    linkForm.resetFields();
    fetchData();
  };

  const columns: ColumnsType<MeetingAction> = [
    { title: t("meetings.action_title"), dataIndex: "title" },
    {
      title: t("meetings.owner"),
      dataIndex: "owner_user_id",
      render: (v?: string) => v || t("common.not_set"),
    },
    {
      title: t("meetings.due"),
      dataIndex: "due_date",
      render: (v?: string) => (v ? dayjs(v).format("YYYY-MM-DD") : "-"),
    },
    {
      title: t("meetings.status"),
      dataIndex: "status",
      render: (v: string) => <Tag color={v === "done" ? "green" : "blue"}>{v}</Tag>,
    },
    {
      title: t("common.actions"),
      render: (_, record) => (
        <Space>
          <Button size="small" onClick={() => updateActionStatus(record.id, record.status === "done" ? "todo" : "done")}> 
            {record.status === "done" ? t("meetings.reopen") : t("meetings.mark_done")}
          </Button>
          <Button size="small" onClick={() => convertAction(record.id)}>{t("meetings.convert_task")}</Button>
        </Space>
      ),
    },
  ];

  return (
    <Space direction="vertical" style={{ width: "100%" }} size="large">
      <Row gutter={16}>
        <Col span={14}>
          <Card loading={loading} title={t("meetings.overview")}>
            {meeting && (
              <Descriptions column={1} bordered size="small">
                <Descriptions.Item label={t("meetings.title")}>{meeting.title}</Descriptions.Item>
                <Descriptions.Item label={t("meetings.time")}>
                  {dayjs(meeting.scheduled_at).format("YYYY-MM-DD HH:mm")}
                </Descriptions.Item>
                <Descriptions.Item label={t("meetings.agenda")}>{meeting.agenda || "-"}</Descriptions.Item>
                <Descriptions.Item label={t("meetings.team")}>{meeting.team_id || t("common.not_set")}</Descriptions.Item>
              </Descriptions>
            )}
          </Card>
          <Card title={t("meetings.notes_editor")} style={{ marginTop: 16 }}>
            <Form form={notesForm} layout="vertical">
              <Form.Item name="agenda" label={t("meetings.agenda")}> 
                <Input.TextArea rows={3} />
              </Form.Item>
              <Form.Item name="notes" label={t("meetings.notes")}> 
                <Input.TextArea rows={6} />
              </Form.Item>
              <Form.Item name="status" label={t("meetings.status")}> 
                <Input />
              </Form.Item>
              <Button type="primary" onClick={saveNotes}>{t("common.save")}</Button>
            </Form>
          </Card>
        </Col>
        <Col span={10}>
          <Card title={t("meetings.action_items")}>
            <Table size="small" rowKey="id" columns={columns} dataSource={actions} pagination={false} />
            <Form form={actionForm} layout="vertical" style={{ marginTop: 16 }}>
              <Form.Item name="title" label={t("meetings.action_title")} rules={[{ required: true, message: t("common.required") }]}> 
                <Input placeholder={t("meetings.action_placeholder")} />
              </Form.Item>
              <Form.Item name="owner_user_id" label={t("meetings.owner")}> 
                <Input placeholder={t("meetings.owner_placeholder")} />
              </Form.Item>
              <Form.Item name="due_date" label={t("meetings.due")}> 
                <Input placeholder="YYYY-MM-DD" />
              </Form.Item>
              <Button type="primary" onClick={addAction}>{t("meetings.add_action")}</Button>
            </Form>
          </Card>
          <Card title={t("meetings.links")} style={{ marginTop: 16 }}>
            <Table
              size="small"
              rowKey="id"
              pagination={false}
              dataSource={links}
              columns={[
                { title: t("meetings.target_type"), dataIndex: "target_type" },
                { title: t("meetings.target_id"), dataIndex: "target_id" },
                { title: t("meetings.relation"), dataIndex: "relation" },
              ]}
            />
            <Form form={linkForm} layout="vertical" style={{ marginTop: 12 }}>
              <Form.Item name="target_type" label={t("meetings.target_type")} rules={[{ required: true }]}> 
                <Input placeholder="okr_objective / okr_key_result / project / report" />
              </Form.Item>
              <Form.Item name="target_id" label={t("meetings.target_id")} rules={[{ required: true }]}> 
                <Input />
              </Form.Item>
              <Form.Item name="relation" label={t("meetings.relation")}> 
                <Input />
              </Form.Item>
              <Button onClick={addLinkAction} type="dashed">{t("meetings.add_link")}</Button>
            </Form>
          </Card>
        </Col>
      </Row>
    </Space>
  );
}
