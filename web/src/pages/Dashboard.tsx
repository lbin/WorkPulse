import React, { useState } from "react";
import { Card, Row, Col, Button, Typography } from "antd";
import { useTranslation } from "react-i18next";
import WorkItemTable from "../components/WorkItemTable";
import LinkOKRModal from "../components/LinkOKRModal";
import { createWorkItem, addOKRLink, WorkItem } from "../api/workItems";

const { Text } = Typography;

// NOTE: replace with a real user_id after you implement auth/login.
// For quick UI smoke test only.
const DUMMY_USER_ID = "00000000-0000-0000-0000-000000000000";

export default function Dashboard() {
  const { t } = useTranslation();
  const [items, setItems] = useState<WorkItem[]>([]);
  const [linkOpen, setLinkOpen] = useState(false);
  const [current, setCurrent] = useState<WorkItem | null>(null);

  return (
    <Row gutter={16}>
      <Col span={14}>
        <Card
          title={t("dashboard.title")}
          extra={
            <Button
              onClick={async () => {
                const wi = await createWorkItem({
                  owner_user_id: DUMMY_USER_ID,
                  type: "TASK",
                  title: t("dashboard.newWorkItem"),
                  priority: 3
                });
                setItems([wi, ...items]);
              }}
            >
              {t("dashboard.newWorkItem")}
            </Button>
          }
        >
          {items.length === 0 ? (
            <Text type="secondary">{t("dashboard.empty")}</Text>
          ) : (
            <WorkItemTable
              data={items}
              onLinkOKR={(wi) => {
                setCurrent(wi);
                setLinkOpen(true);
              }}
            />
          )}
        </Card>
      </Col>

      <Col span={10}>
        <Card title={t("dashboard.okrHintTitle")}>
          {t("dashboard.okrHintBody")}
        </Card>
      </Col>

      <LinkOKRModal
        open={linkOpen}
        onCancel={() => setLinkOpen(false)}
        onOk={async (v) => {
          if (!current) return;
          await addOKRLink(current.id, v);
          setLinkOpen(false);
        }}
      />
    </Row>
  );
}
