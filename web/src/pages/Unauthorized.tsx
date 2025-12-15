import React from "react";
import { Button, Result } from "antd";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

export default function Unauthorized() {
  const nav = useNavigate();
  const { t } = useTranslation();
  return (
    <Result
      status="403"
      title={t("common.unauthorizedTitle", "No access")}
      subTitle={t("common.unauthorizedDesc", "You do not have permission to view this page.")}
      extra={
        <Button type="primary" onClick={() => nav("/dashboard")}> 
          {t("common.backHome", "Back to dashboard")}
        </Button>
      }
    />
  );
}
