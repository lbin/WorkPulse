import React from "react";
import { Table, Button } from "antd";
import type { ColumnsType } from "antd/es/table";
import type { WorkItem } from "../api/workItems";

export default function WorkItemTable(props: {
  data: WorkItem[];
  onLinkOKR: (wi: WorkItem) => void;
}) {
  const cols: ColumnsType<WorkItem> = [
    { title: "标题", dataIndex: "title" },
    { title: "类型", dataIndex: "type", width: 140 },
    { title: "状态", dataIndex: "status", width: 120 },
    { title: "优先级", dataIndex: "priority", width: 100 },
    { title: "操作", width: 160, render: (_, r) => <Button onClick={() => props.onLinkOKR(r)}>关联 OKR</Button> },
  ];
  return <Table rowKey="id" columns={cols} dataSource={props.data} pagination={false} />;
}
