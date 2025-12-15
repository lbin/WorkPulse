import React, { useState } from "react";
import { Modal, Form, Select, InputNumber, Switch } from "antd";
import { useTranslation } from "react-i18next";

type Props = {
  open: boolean;
  onCancel: () => void;
  onOk: (v: any) => Promise<void> | void;
};

export default function LinkOKRModal({ open, onCancel, onOk }: Props) {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  const { t } = useTranslation();

  return (
    <Modal
      title={t("okrLink.title")}
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
        <Form.Item name="key_result_id" label={t("okrLink.kr")}>
          <Select placeholder={t("okrLink.selectKR")} options={[]} />
        </Form.Item>

        <Form.Item name="objective_id" label={t("okrLink.objective")}>
          <Select placeholder={t("okrLink.selectObj")} options={[]} />
        </Form.Item>

        <Form.Item name="link_type" label={t("okrLink.linkType")} rules={[{ required: true }]}>
          <Select
            options={[
              { value: "planned", label: t("okrLink.planned") },
              { value: "actual", label: t("okrLink.actual") },
              { value: "both", label: t("okrLink.both") },
            ]}
          />
        </Form.Item>

        <Form.Item name="contribution_weight" label={t("okrLink.weight")}>
          <InputNumber min={0} max={1} step={0.1} style={{ width: "100%" }} />
        </Form.Item>

        <Form.Item name="evidence_required" label={t("okrLink.evidenceRequired")} valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Modal>
  );
}
