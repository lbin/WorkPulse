import React, { useState } from "react";
import { Button, Card, Form, Input, Tabs, Typography, message } from "antd";
import { useNavigate, useLocation } from "react-router-dom";
import { login, register } from "../api/auth";
import { useAuth } from "../components/AuthProvider";
import { useTranslation } from "react-i18next";

// Login renders a combined login/register experience using email + password.
export default function Login() {
  const { refresh } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [tab, setTab] = useState<string>("login");

  const handleAuth = async (values: any, mode: "login" | "register") => {
    setLoading(true);
    try {
      const action = mode === "login" ? login : register;
      const res = mode === "login" ? await action(values.email, values.password) : await (action as typeof register)(values.email, values.password, values.display_name);
      localStorage.setItem("token", res.token);
      await refresh();
      message.success(mode === "login" ? t("auth.loginSuccess") : t("auth.registerSuccess"));
      const hasTeam = res.context?.org_units?.length > 0;
      const redirectTarget = hasTeam ? (location.state as any)?.from || "/dashboard" : "/onboarding";
      navigate(redirectTarget, { replace: true });
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || t("auth.loginFailed"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "100vh" }}>
      <Card title={t("app.name")}
        style={{ width: 420 }}
      >
        <Tabs activeKey={tab} onChange={(k) => setTab(k)} items={[
          { key: "login", label: t("auth.login"), children: (
            <Form layout="vertical" onFinish={(vals) => handleAuth(vals, "login")}>
              <Form.Item label={t("auth.email")} name="email" rules={[{ required: true, message: t("common.required") }, { type: "email", message: t("auth.emailInvalid") }]}>
                <Input placeholder={t("auth.emailPlaceholder") as string} />
              </Form.Item>
              <Form.Item label={t("auth.password")} name="password" rules={[{ required: true, message: t("common.required") }]}> 
                <Input.Password placeholder={t("auth.passwordPlaceholder") as string} />
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" block loading={loading}>
                  {t("auth.login")}
                </Button>
              </Form.Item>
            </Form>
          ) },
          { key: "register", label: t("auth.register"), children: (
            <Form layout="vertical" onFinish={(vals) => handleAuth(vals, "register")}>
              <Form.Item label={t("auth.email")} name="email" rules={[{ required: true, message: t("common.required") }, { type: "email", message: t("auth.emailInvalid") }]}>
                <Input placeholder={t("auth.emailPlaceholder") as string} />
              </Form.Item>
              <Form.Item label={t("auth.password")} name="password" rules={[{ required: true, message: t("common.required") }]}> 
                <Input.Password placeholder={t("auth.passwordPlaceholder") as string} />
              </Form.Item>
              <Form.Item label={t("auth.displayName")} name="display_name">
                <Input placeholder={t("auth.displayNamePlaceholder") as string} />
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" block loading={loading}>
                  {t("auth.register")}
                </Button>
              </Form.Item>
            </Form>
          ) }
        ]}
        />
        <Typography.Paragraph type="secondary" style={{ marginTop: 8 }}>
          {t("auth.credentialHint")}
        </Typography.Paragraph>
      </Card>
    </div>
  );
}
