import React, { useMemo } from "react";
import { Layout, Menu, Select, Space, Typography } from "antd";
import { Link, Outlet, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import i18n from "../i18n";
import { useAuth, useOrgUnits } from "./AuthProvider";
import { useConfig } from "./ConfigProvider";

const { Sider, Header, Content } = Layout;
const { Text } = Typography;

export default function AppLayout() {
  const loc = useLocation();
  const key = loc.pathname.startsWith("/okr") ? "/okr" : loc.pathname;
  const { t } = useTranslation();
  const { profile, setActiveOrgUnit, activeOrgUnit, hasPermission } = useAuth();
  const { isFeatureEnabled } = useConfig();
  const units = useOrgUnits();

  const items = useMemo(
    () =>
      [
        { key: "/dashboard", label: <Link to="/dashboard">{t("nav.dashboard")}</Link>, permission: "dashboard.view" },
        { key: "/okr", label: <Link to="/okr">{t("nav.okr")}</Link>, permission: "okr.view" },
        { key: "/projects", label: <Link to="/projects">{t("nav.projects")}</Link>, permission: "projects.view" },
        { key: "/meetings", label: <Link to="/meetings">{t("nav.meetings")}</Link>, permission: "meetings.view" },
        { key: "/reports", label: <Link to="/reports">{t("nav.reports")}</Link>, permission: "reports.view" },
        { key: "/analytics", label: <Link to="/analytics">{t("nav.analytics")}</Link>, permission: "analytics.view", flag: "analytics" },
        { key: "/docs", label: <Link to="/docs">{t("nav.docs")}</Link>, permission: "dashboard.view", flag: "apiDocs" },
      ].filter((item) => hasPermission(item.permission) && (!item.flag || isFeatureEnabled(item.flag, true))),
    [hasPermission, isFeatureEnabled, t]
  );

  const langValue = i18n.language?.startsWith("en") ? "en" : "zh";

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider width={220}>
        <div style={{ color: "#fff", padding: 16, fontWeight: 600 }}>{t("app.name")}</div>
        <Menu theme="dark" mode="inline" selectedKeys={[key]} items={items} />
      </Sider>
      <Layout>
        <Header style={{ background: "#fff", display: "flex", alignItems: "center", justifyContent: "space-between", gap: 12, paddingInline: 16 }}>
          <Space>
            <div>
              <Text type="secondary" style={{ marginRight: 8 }}>
                {t("common.teamSwitch", "Team/Org unit")}
              </Text>
              <Select
                value={activeOrgUnit}
                placeholder={t("common.selectTeam", "Select team")}
                style={{ minWidth: 220 }}
                options={units.map((u) => ({ value: u.id, label: u.name }))}
                onChange={(v) => setActiveOrgUnit(v)}
              />
            </div>
            {profile?.user && <Text>{profile.user.display_name}</Text>}
          </Space>
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
