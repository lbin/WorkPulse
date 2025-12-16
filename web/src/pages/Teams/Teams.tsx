import React, { useEffect, useState } from "react";
import { Button, Card, Flex, Form, Input, Space, Table, Tag, Typography, message } from "antd";
import { useTranslation } from "react-i18next";
import { createTeam, joinTeam, listMyTeams, TeamSummary } from "../../api/auth";
import { useAuth } from "../../components/AuthProvider";

export default function TeamsPage() {
  const { t } = useTranslation();
  const { activeOrgUnit, setActiveOrgUnit, refresh } = useAuth();
  const [teams, setTeams] = useState<TeamSummary[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadTeams();
  }, []);

  const columns = [
    { title: t("common.title", { defaultValue: "Title" }), dataIndex: "name" },
    { title: "ID", dataIndex: "id" },
    { title: t("common.actions", { defaultValue: "Actions" }), render: (_: unknown, record: TeamSummary) => renderActions(record) },
  ];

  const renderActions = (team: TeamSummary) => {
    const isActive = activeOrgUnit === team.id;
    return (
      <Space>
        {isActive ? (
          <Tag color="green">{t("teams.active", { defaultValue: "Active" })}</Tag>
        ) : (
          <Button size="small" onClick={() => handleActivate(team.id)}>
            {t("teams.setActive", { defaultValue: "Set active" })}
          </Button>
        )}
      </Space>
    );
  };

  const handleActivate = async (id: string) => {
    setActiveOrgUnit(id);
    message.success(t("teams.switched", { defaultValue: "Switched team" }));
  };

  const loadTeams = async () => {
    setLoading(true);
    try {
      const data = await listMyTeams();
      setTeams(data);
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || t("teams.loadFailed", { defaultValue: "Failed to load teams" }));
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (values: any) => {
    setLoading(true);
    try {
      await createTeam(values.name, values.parent_team_id || undefined);
      message.success(t("teams.created", { defaultValue: "Team created" }));
      await refresh();
      await loadTeams();
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || t("teams.actionFailed", { defaultValue: "Action failed" }));
    } finally {
      setLoading(false);
    }
  };

  const handleJoin = async (values: any) => {
    setLoading(true);
    try {
      await joinTeam(values.team_id);
      message.success(t("teams.joined", { defaultValue: "Joined team" }));
      await refresh();
      await loadTeams();
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || t("teams.actionFailed", { defaultValue: "Action failed" }));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Flex vertical gap={16}>
      <Card
        title={t("teams.title", { defaultValue: "Team management" })}
        extra={<Typography.Text type="secondary">{t("teams.subtitle", { defaultValue: "Create or join teams for your org" })}</Typography.Text>}
      >
        <Table
          rowKey="id"
          loading={loading}
          dataSource={teams}
          columns={columns}
          pagination={false}
        />
        {!teams.length && (
          <Typography.Text type="secondary">{t("teams.empty", { defaultValue: "No teams yet. Create one below." })}</Typography.Text>
        )}
      </Card>

      <Card title={t("teams.create", { defaultValue: "Create team" })}>
        <Form layout="vertical" onFinish={handleCreate}>
          <Form.Item
            label={t("teams.teamName", { defaultValue: "Team name" })}
            name="name"
            rules={[{ required: true, message: t("common.required") }]}
          >
            <Input placeholder={t("auth.teamNamePlaceholder") as string} />
          </Form.Item>
          <Form.Item label={t("teams.parentTeam", { defaultValue: "Parent team (optional)" })} name="parent_team_id">
            <Input placeholder={t("auth.parentTeamPlaceholder") as string} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t("teams.create", { defaultValue: "Create team" })}
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title={t("teams.join", { defaultValue: "Join team" })}>
        <Form layout="vertical" onFinish={handleJoin}>
          <Form.Item
            label={t("auth.teamId")}
            name="team_id"
            rules={[{ required: true, message: t("common.required") }]}
          >
            <Input placeholder={t("auth.teamIdPlaceholder") as string} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t("teams.join", { defaultValue: "Join team" })}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </Flex>
  );
}
