import React, { useEffect, useState } from "react";
import { Card, Form, Input, Button, List, message } from "antd";
import { useNavigate } from "react-router-dom";
import { createTeam, joinTeam, listOrgTeams, listMyTeams, Team } from "../../api/teams";

// TeamSetup lets users create a new team or join an existing one after login.
export default function TeamSetup() {
  const [teams, setTeams] = useState<Team[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // Load teams and redirect if the user is already a member.
  const loadData = async () => {
    const availableTeams = await listOrgTeams();
    setTeams(availableTeams);
    const memberships = await listMyTeams();
    if (memberships.length > 0) {
      const currentTeam = availableTeams.find(
        (t) => t.id === memberships[0].team_id
      );
      localStorage.setItem("activeTeamId", memberships[0].team_id);
      if (currentTeam) {
        localStorage.setItem("activeTeamName", currentTeam.name);
      }
      navigate("/dashboard", { replace: true });
      return;
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  // Submit a new team then set it as active for the session.
  const onCreateTeam = async (values: any) => {
    setLoading(true);
    try {
      const team = await createTeam(values.name);
      localStorage.setItem("activeTeamId", team.id);
      localStorage.setItem("activeTeamName", team.name);
      message.success("团队创建成功");
      navigate("/dashboard", { replace: true });
    } catch (err: any) {
      message.error(err?.response?.data?.error || "创建团队失败");
    } finally {
      setLoading(false);
    }
  };

  // Join an existing team and mark it active.
  const handleJoin = async (teamId: string) => {
    setLoading(true);
    try {
      await joinTeam(teamId);
      localStorage.setItem("activeTeamId", teamId);
      const team = teams.find((t) => t.id === teamId);
      if (team) {
        localStorage.setItem("activeTeamName", team.name);
      }
      message.success("加入团队成功");
      navigate("/dashboard", { replace: true });
    } catch (err: any) {
      message.error(err?.response?.data?.error || "加入失败");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        display: "flex",
        gap: 16,
        flexWrap: "wrap",
        padding: 24,
        minHeight: "100vh",
        background: "#f5f5f5",
      }}
    >
      <Card title="创建团队" style={{ width: 400 }}>
        <Form layout="vertical" onFinish={onCreateTeam}>
          <Form.Item
            label="团队名称"
            name="name"
            rules={[{ required: true, message: "请输入团队名称" }]}
          >
            <Input placeholder="例如：增长小组" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              创建并进入
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="加入已有团队" style={{ flex: 1, minWidth: 360 }}>
        <List
          bordered
          dataSource={teams}
          locale={{ emptyText: "暂无团队，请先创建" }}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button
                  type="link"
                  onClick={() => handleJoin(item.id)}
                  key={item.id}
                  loading={loading}
                >
                  加入
                </Button>,
              ]}
            >
              <List.Item.Meta title={item.name} description={item.path} />
            </List.Item>
          )}
        />
      </Card>
    </div>
  );
}
