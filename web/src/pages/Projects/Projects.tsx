import React, { useEffect, useMemo, useState } from "react";
import {
  Badge,
  Button,
  Card,
  Col,
  Divider,
  Drawer,
  Empty,
  Form,
  Input,
  Progress,
  Row,
  Select,
  Space,
  Table,
  Tabs,
  Tag,
  Timeline,
  Typography,
} from "antd";
import { useTranslation } from "react-i18next";
import { usePermission } from "../../components/AuthProvider";
import {
  Project,
  Task,
  TaskLink,
  fetchMilestones,
  fetchProjectOKRProgress,
  fetchProjects,
  fetchTaskLinks,
  fetchTasks,
  linkTask,
  unlinkTask,
} from "../../api/projects";

const { Title, Paragraph, Text } = Typography;

const statusColor: Record<string, string> = {
  draft: "default",
  active: "processing",
  paused: "warning",
  closed: "success",
  todo: "default",
  doing: "processing",
  done: "success",
};

export default function ProjectsPage() {
  const { t } = useTranslation();
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<string | null>(null);
  const [milestones, setMilestones] = useState<any[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [okrProgress, setOkrProgress] = useState<{ project_progress: number; key_results: { key_result_id: string; progress: number }[] } | null>(
    null
  );
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [taskLinks, setTaskLinks] = useState<TaskLink[]>([]);
  const [linkForm] = Form.useForm();
  const canManage = usePermission("projects.manage");

  useEffect(() => {
    fetchProjects()
      .then((res) => {
        setProjects(res);
        if (res.length > 0) {
          setSelectedProjectId(res[0].id);
        }
      })
      .catch((err) => console.error(err));
  }, []);

  useEffect(() => {
    if (!selectedProjectId) return;
    Promise.all([
      fetchMilestones(selectedProjectId),
      fetchTasks(selectedProjectId),
      fetchProjectOKRProgress(selectedProjectId),
    ])
      .then(([ms, ts, progress]) => {
        setMilestones(ms);
        setTasks(ts);
        setOkrProgress(progress);
      })
      .catch((err) => console.error(err));
  }, [selectedProjectId]);

  const milestoneTimeline = useMemo(
    () =>
      milestones.map((m: any) => ({
        label: m.due_date ? new Date(m.due_date).toLocaleDateString() : t("common.tbd"),
        children: (
          <Space direction="vertical" size="small">
            <Text strong>{m.title}</Text>
            <Text type="secondary">{m.description}</Text>
            <Progress percent={Math.round(m.progress)} size="small" />
          </Space>
        ),
      })),
    [milestones, t]
  );

  const groupedTasks = useMemo(() => {
    const map: Record<string, Task[]> = { todo: [], doing: [], done: [] };
    tasks.forEach((t) => {
      const bucket = map[t.status] || map.todo;
      bucket.push(t);
    });
    return map;
  }, [tasks]);

  const openTask = async (task: Task) => {
    setSelectedTask(task);
    setDrawerOpen(true);
    try {
      const links = await fetchTaskLinks(task.id);
      setTaskLinks(links);
    } catch (err) {
      console.error(err);
    }
  };

  const handleLinkSubmit = async () => {
    if (!selectedTask) return;
    if (!canManage) return;
    const values = await linkForm.validateFields();
    const link = await linkTask(selectedTask.id, values.entity_type, values.entity_id);
    setTaskLinks((prev) => [link, ...prev]);
    linkForm.resetFields();
  };

  const handleUnlink = async (linkId: string) => {
    if (!selectedTask) return;
    if (!canManage) return;
    await unlinkTask(selectedTask.id, linkId);
    setTaskLinks((prev) => prev.filter((l) => l.id !== linkId));
  };

  const columns = [
    {
      title: t("projects.name", "Name"),
      dataIndex: "name",
      key: "name",
      render: (value: string, record: Project) => (
        <Space direction="vertical" size={0}>
          <Text strong>{value}</Text>
          <Text type="secondary">{record.description}</Text>
        </Space>
      ),
    },
    {
      title: t("common.status", "Status"),
      dataIndex: "status",
      render: (status: string) => <Badge status={statusColor[status] as any} text={status} />,
    },
    {
      title: t("projects.progress", "Progress"),
      dataIndex: "progress",
      render: (value: number) => <Progress percent={Math.round(value)} size="small" />,
    },
    {
      title: t("projects.milestones", "Milestones"),
      dataIndex: "milestone_count",
    },
    {
      title: t("projects.tasks", "Tasks"),
      dataIndex: "task_count",
    },
  ];

  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      <Title level={3}>{t("projects.title", "Projects & Roadmap")}</Title>
      <Tabs
        defaultActiveKey="list"
        items={[
          {
            key: "list",
            label: t("projects.list", "Project list"),
            children: (
              <Card>
                <Table
                  dataSource={projects}
                  columns={columns as any}
                  rowKey="id"
                  onRow={(record) => ({
                    onClick: () => setSelectedProjectId(record.id),
                  })}
                  pagination={false}
                />
              </Card>
            ),
          },
          {
            key: "timeline",
            label: t("projects.timeline", "Gantt / timeline"),
            children: (
              <Row gutter={16}>
                <Col span={16}>
                  <Card title={t("projects.timeline", "Timeline")}>{milestoneTimeline.length > 0 ? <Timeline items={milestoneTimeline} /> : <Empty />}</Card>
                </Col>
                <Col span={8}>
                  <Card title={t("projects.okr", "OKR rollup")}>{
                    okrProgress ? (
                      <Space direction="vertical" style={{ width: "100%" }}>
                        <Paragraph>
                          {t("projects.projectProgress", "Project progress")}:
                          <Progress percent={Math.round(okrProgress.project_progress)} />
                        </Paragraph>
                        {okrProgress.key_results.map((kr) => (
                          <div key={kr.key_result_id}>
                            <Text strong>{kr.key_result_id}</Text>
                            <Progress percent={Math.round(kr.progress)} size="small" />
                          </div>
                        ))}
                      </Space>
                    ) : (
                      <Empty />
                    )
                  }</Card>
                </Col>
              </Row>
            ),
          },
          {
            key: "kanban",
            label: t("projects.kanban", "Task kanban"),
            children: (
              <Row gutter={12}>
                {Object.entries(groupedTasks).map(([status, list]) => (
                  <Col span={8} key={status}>
                    <Card
                      title={
                        <Space>
                          <Badge status={statusColor[status] as any} /> <span>{status.toUpperCase()}</span>
                        </Space>
                      }
                      extra={
                        <Text type="secondary">
                          {list.length} {t("projects.tasks", "Tasks")}
                        </Text>
                      }
                    >
                      <Space direction="vertical" style={{ width: "100%" }}>
                        {list.length === 0 && <Empty />}
                        {list.map((task) => (
                          <Card key={task.id} size="small" hoverable onClick={() => openTask(task)}>
                            <Space direction="vertical" size={4} style={{ width: "100%" }}>
                              <Text strong>{task.title}</Text>
                              <Space wrap>
                                {task.tags?.map((tag: string) => (
                                  <Tag key={tag}>{tag}</Tag>
                                ))}
                              </Space>
                            </Space>
                          </Card>
                        ))}
                      </Space>
                    </Card>
                  </Col>
                ))}
              </Row>
            ),
          },
        ]}
      />

      <Drawer
        width={420}
        title={selectedTask?.title}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        destroyOnClose
      >
        {selectedTask ? (
          <Space direction="vertical" style={{ width: "100%" }} size="middle">
            <Badge status={statusColor[selectedTask.status] as any} text={selectedTask.status} />
            {selectedTask.parent_task_id && (
              <Text type="secondary">{t("projects.subtaskOf", "Subtask of")} {selectedTask.parent_task_id}</Text>
            )}
            <Divider orientation="left">{t("projects.links", "Links to OKRs / reports")}</Divider>
            <Space direction="vertical" style={{ width: "100%" }}>
              {taskLinks.length === 0 && <Text type="secondary">{t("projects.noLinks", "No links yet")}</Text>}
              {taskLinks.map((l) => (
                <Card key={l.id} size="small">
                  <Space style={{ width: "100%" }}>
                    <Text>{l.entity_type}</Text>
                    <Text code>{l.entity_id}</Text>
                    {canManage && (
                      <Button size="small" type="link" danger onClick={() => handleUnlink(l.id)}>
                        {t("common.remove", "Remove")}
                      </Button>
                    )}
                  </Space>
                </Card>
              ))}
            </Space>
            {canManage ? (
              <Form layout="vertical" form={linkForm} onFinish={handleLinkSubmit}>
                <Form.Item name="entity_type" label={t("projects.linkType", "Link type")} rules={[{ required: true }]}>
                  <Select
                    options={[
                      { label: "OKR KR", value: "okr_kr" },
                      { label: "OKR Objective", value: "okr_objective" },
                      { label: "Report", value: "report" },
                    ]}
                  />
                </Form.Item>
                <Form.Item name="entity_id" label={t("projects.linkId", "Linked ID")} rules={[{ required: true }]}>
                  <Input placeholder="UUID" />
                </Form.Item>
                <Button type="primary" htmlType="submit" block>
                  {t("projects.addLink", "Add link")}
                </Button>
              </Form>
            ) : (
              <Text type="secondary">{t("projects.noManagePermission", "Linking requires project edit permission.")}</Text>
            )}
          </Space>
        ) : (
          <Empty />
        )}
      </Drawer>
    </Space>
  );
}
