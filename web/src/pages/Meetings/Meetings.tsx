import React from "react";
import { Card } from "antd";
import { useTranslation } from "react-i18next";

export default function Page() {
  const { t } = useTranslation();
  return <Card>{t("common.todo")}</Card>;
}
