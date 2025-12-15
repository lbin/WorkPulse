import React from "react";
import { Layout, Menu, Select } from "antd";
import { Link, Outlet, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import i18n from "../i18n";

const { Sider, Header, Content } = Layout;

export default function AppLayout() {
  const loc = useLocation();
  const key = loc.pathname.startsWith("/okr") ? "/okr" : loc.pathname;
  const { t } = useTranslation();

  const items = [
    { key: "/dashboard", label: <Link to="/dashboard">{t("nav.dashboard")}</Link> },
    { key: "/okr", label: <Link to="/okr">{t("nav.okr")}</Link> },
    { key: "/projects", label: <Link to="/projects">{t("nav.projects")}</Link> },
    { key: "/meetings", label: <Link to="/meetings">{t("nav.meetings")}</Link> },
    { key: "/reports", label: <Link to="/reports">{t("nav.reports")}</Link> },
    { key: "/analytics", label: <Link to="/analytics">{t("nav.analytics")}</Link> },
  ];

  const langValue = i18n.language?.startsWith("en") ? "en" : "zh";

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider width={220}>
        <div style={{ color: "#fff", padding: 16, fontWeight: 600 }}>{t("app.name")}</div>
        <Menu theme="dark" mode="inline" selectedKeys={[key]} items={items} />
      </Sider>
      <Layout>
        <Header style={{ background: "#fff", display: "flex", alignItems: "center", justifyContent: "flex-end", gap: 12 }}>
          <Select
            value={langValue}
            style={{ width: 140 }}
            options={[
              { value: "zh", label: t("common.chinese") },
              { value: "en", label: t("common.english") },
            ]}
            onChange={(v) => i18n.changeLanguage(v)}
          />
        </Header>
        <Content style={{ padding: 16 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
