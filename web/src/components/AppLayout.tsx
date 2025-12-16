import React from "react";
import { Layout, Menu } from "antd";
import { Link, Outlet, useLocation } from "react-router-dom";

const { Sider, Header, Content } = Layout;

const items = [
  { key: "/dashboard", label: <Link to="/dashboard">今日</Link> },
  { key: "/okr", label: <Link to="/okr">OKR</Link> },
  { key: "/projects", label: <Link to="/projects">项目</Link> },
  { key: "/meetings", label: <Link to="/meetings">纪要</Link> },
  { key: "/reports", label: <Link to="/reports">报告</Link> },
  { key: "/analytics", label: <Link to="/analytics">分析</Link> },
];

export default function AppLayout() {
  const loc = useLocation();
  const key = loc.pathname.startsWith("/okr") ? "/okr" : loc.pathname;
  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider width={220}>
        <div style={{ color: "#fff", padding: 16, fontWeight: 600 }}>WorkPulse</div>
        <Menu theme="dark" mode="inline" selectedKeys={[key]} items={items} />
      </Sider>
      <Layout>
        <Header style={{ background: "#fff" }} />
        <Content style={{ padding: 16 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
