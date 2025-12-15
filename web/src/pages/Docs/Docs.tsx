import React from "react";
import { Alert, Card, Space, Typography } from "antd";
import { useTranslation } from "react-i18next";
import { useConfig } from "../../components/ConfigProvider";

export default function Docs() {
  const { t } = useTranslation();
  const { config } = useConfig();

  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      <Typography.Title level={3}>{t("docs.title", "API Documentation")}</Typography.Title>
      <Alert
        type="info"
        showIcon
        message={t("docs.runtime", "Docs are generated from the active router and versioned gateway.")}
      />
      <Card title={t("docs.openapi", "OpenAPI endpoint")}> 
        <Typography.Paragraph copyable={{ text: config?.docs?.openapi || "" }}>
          {config?.docs?.openapi || t("docs.unavailable", "OpenAPI URL unavailable")}
        </Typography.Paragraph>
      </Card>
      <Card title={t("docs.graphql", "GraphQL SDL endpoint")}> 
        <Typography.Paragraph copyable={{ text: config?.docs?.graphql_sdl || "" }}>
          {config?.docs?.graphql_sdl || t("docs.unavailable", "GraphQL SDL URL unavailable")}
        </Typography.Paragraph>
      </Card>
    </Space>
  );
}
