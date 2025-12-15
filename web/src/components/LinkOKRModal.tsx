import React, { useState } from "react";
import { Modal, Form, Select, InputNumber, Switch } from "antd";

type Props = {
  open: boolean;
  onCancel: () => void;
  onOk: (v: any) => Promise<void> | void;
};

export default function LinkOKRModal({ open, onCancel, onOk }: Props) {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  return (
    <Modal
      title="关联 OKR"
      open={open}
      onCancel={onCancel}
      okButtonProps={{ loading }}
      onOk={async () => {
        const v = await form.validateFields();
        setLoading(true);
        try {
          await onOk(v);
          form.resetFields();
        } finally {
          setLoading(false);
        }
      }}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{ link_type: "actual", contribution_weight: 1.0, evidence_required: false }}
      >
        <Form.Item name="key_result_id" label="Key Result">
          <Select placeholder="选择 KR（接入后端列表）" options={[]} />
        </Form.Item>
        <Form.Item name="objective_id" label="Objective">
          <Select placeholder="选择 Objective（可选）" options={[]} />
        </Form.Item>
        <Form.Item name="link_type" label="关联类型" rules={[{ required: true }]}>
          <Select
            options={[
              { value: "planned", label: "计划" },
              { value: "actual", label: "实际" },
              { value: "both", label: "计划+实际" },
            ]}
          />
        </Form.Item>
        <Form.Item name="contribution_weight" label="贡献权重">
          <InputNumber min={0} max={1} step={0.1} style={{ width: "100%" }} />
        </Form.Item>
        <Form.Item name="evidence_required" label="需要证据" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Modal>
  );
}
