import React, { useEffect, useState } from "react";
import { Button, Card, Form, Input, Typography, message } from "antd";
import { useAuth } from "../components/AuthProvider";
import { createTeam, joinTeam, listMyTeams, TeamSummary } from "../api/auth";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

// Onboarding helps authenticated users create or join a team before entering the app.
export default function Onboarding() {
  const { profile, refresh } = useAuth();
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [teams, setTeams] = useState<TeamSummary[]>(
    profile?.org_units?.map((u) => ({ id: u.id, name: u.name, path: u.path })) || []
  );
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (profile?.org_units?.length) {
      setTeams(profile.org_units.map((u) => ({ id: u.id, name: u.name, path: u.path })));
    } else {
      listMyTeams().then(setTeams).catch(() => setTeams([]));
    }
  }, [profile]);

  useEffect(() => {
    if (teams.length > 0) {
      navigate("/dashboard", { replace: true });
    }
  }, [teams, navigate]);

  const handleCreate = async (values: any) => {
    setLoading(true);
    try {
      await createTeam(values.name, values.parent_team_id || undefined);
      message.success(t("auth.teamCreated"));
      await refresh();
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || t("auth.actionFailed"));
    } finally {
      setLoading(false);
    }
  };

  const handleJoin = async (values: any) => {
    setLoading(true);
    try {
      await joinTeam(values.team_id);
      message.success(t("auth.teamJoined"));
      await refresh();
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || t("auth.actionFailed"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 720, margin: "32px auto" }}>
      <Typography.Title level={3}>{t("auth.onboardingTitle")}</Typography.Title>
      <Typography.Paragraph>{t("auth.onboardingBody")}</Typography.Paragraph>

      <Card title={t("auth.createTeam") } style={{ marginBottom: 16 }}>
        <Form layout="vertical" onFinish={handleCreate}>
          <Form.Item label={t("auth.teamName")} name="name" rules={[{ required: true, message: t("common.required") }]}>
            <Input placeholder={t("auth.teamNamePlaceholder") as string} />
          </Form.Item>
          <Form.Item label={t("auth.parentTeamOptional")} name="parent_team_id">
            <Input placeholder={t("auth.parentTeamPlaceholder") as string} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t("auth.createTeam")}
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title={t("auth.joinTeam") }>
        <Form layout="vertical" onFinish={handleJoin}>
          <Form.Item label={t("auth.teamId") } name="team_id" rules={[{ required: true, message: t("common.required") }]}>
            <Input placeholder={t("auth.teamIdPlaceholder") as string} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t("auth.joinTeam")}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
