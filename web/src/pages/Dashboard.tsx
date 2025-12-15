import React, { useState } from "react";
import { Card, Row, Col, Button, Typography } from "antd";
import WorkItemTable from "../components/WorkItemTable";
import LinkOKRModal from "../components/LinkOKRModal";
import { createWorkItem, addOKRLink, WorkItem } from "../api/workItems";

const { Text } = Typography;

// NOTE: replace with a real user_id after you implement auth/login.
// For quick UI smoke test only.
const DUMMY_USER_ID = "00000000-0000-0000-0000-000000000000";

export default function Dashboard() {
  const [items, setItems] = useState<WorkItem[]>([]);
  const [linkOpen, setLinkOpen] = useState(false);
  const [current, setCurrent] = useState<WorkItem | null>(null);

  return (
    <Row gutter={16}>
      <Col span={14}>
        <Card
          title="今日产出 / 待办"
          extra={
            <Button
              onClick={async () => {
                const wi = await createWorkItem({
                  owner_user_id: DUMMY_USER_ID,
                  type: "TASK",
                  title: "新任务",
                  priority: 3
                });
                setItems([wi, ...items]);
              }}
            >
              新建工作项
            </Button>
          }
        >
          {items.length === 0 ? (
            <Text type="secondary">暂无数据。点击右上角“新建工作项”做一次联调。</Text>
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
        <Card title="OKR 贡献提示（占位）">
          这里接入 analytics 接口后展示：覆盖率、投入偏差、无归属工作 Top。
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
