import React, { useState } from "react";
import { Card, Tabs, Form, Input, Button, message } from "antd";
import { useNavigate } from "react-router-dom";
import { login, register } from "../../api/auth";
import { listMyTeams, listOrgTeams } from "../../api/teams";

// AuthPage renders login and registration forms for email/password.
export default function AuthPage() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // After a successful login/registration, check memberships and redirect accordingly.
  const handlePostAuth = async () => {
    const [memberships, teams] = await Promise.all([
      listMyTeams(),
      listOrgTeams(),
    ]);
    if (memberships.length > 0) {
      const [first] = memberships;
      localStorage.setItem("activeTeamId", first.team_id);
      const targetTeam = teams.find((t) => t.id === first.team_id);
      if (targetTeam) {
        localStorage.setItem("activeTeamName", targetTeam.name);
      }
      navigate("/dashboard", { replace: true });
    } else {
      navigate("/team-setup", { replace: true });
    }
  };

  // Common handler for persisting token and routing decisions.
  const handleAuthResult = async (token: string) => {
    localStorage.setItem("token", token);
    await handlePostAuth();
  };

  // Submit login and save token when it succeeds.
  const onLogin = async (values: any) => {
    setLoading(true);
    try {
      const res = await login(values.email, values.password);
      await handleAuthResult(res.token);
    } catch (err: any) {
      message.error(err?.response?.data?.error || "登录失败");
    } finally {
      setLoading(false);
    }
  };

  // Submit registration, save token, then prompt team setup.
  const onRegister = async (values: any) => {
    setLoading(true);
    try {
      const res = await register(
        values.email,
        values.password,
        values.display_name,
        values.org_name
      );
      await handleAuthResult(res.token);
    } catch (err: any) {
      message.error(err?.response?.data?.error || "注册失败");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "100vh",
        background: "#f5f5f5",
      }}
    >
      <Card title="WorkPulse 登录/注册" style={{ width: 420 }}>
        <Tabs defaultActiveKey="login">
          <Tabs.TabPane tab="登录" key="login">
            <Form layout="vertical" onFinish={onLogin}>
              <Form.Item
                label="邮箱"
                name="email"
                rules={[{ required: true, message: "请输入邮箱" }]}
              >
                <Input type="email" placeholder="you@example.com" />
              </Form.Item>
              <Form.Item
                label="密码"
                name="password"
                rules={[{ required: true, message: "请输入密码" }]}
              >
                <Input.Password placeholder="不少于6位" />
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading} block>
                  登录
                </Button>
              </Form.Item>
            </Form>
          </Tabs.TabPane>
          <Tabs.TabPane tab="注册" key="register">
            <Form layout="vertical" onFinish={onRegister}>
              <Form.Item
                label="邮箱"
                name="email"
                rules={[{ required: true, message: "请输入邮箱" }]}
              >
                <Input type="email" />
              </Form.Item>
              <Form.Item
                label="密码"
                name="password"
                rules={[{ required: true, message: "请输入密码" }]}
              >
                <Input.Password placeholder="不少于6位" />
              </Form.Item>
              <Form.Item label="昵称" name="display_name">
                <Input placeholder="用于展示的名称" />
              </Form.Item>
              <Form.Item label="团队名称" name="org_name">
                <Input placeholder="可选，默认使用邮箱前缀" />
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading} block>
                  注册并登录
                </Button>
              </Form.Item>
            </Form>
          </Tabs.TabPane>
        </Tabs>
      </Card>
    </div>
  );
}
