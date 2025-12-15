import React from "react";
import { Navigate } from "react-router-dom";
import { useFeatureFlag } from "./ConfigProvider";

interface Props {
  flag: string;
  defaultEnabled?: boolean;
  children: React.ReactNode;
}

export function FeatureFlagGate({ flag, defaultEnabled = true, children }: Props) {
  const enabled = useFeatureFlag(flag, defaultEnabled);
  if (!enabled) {
    return <Navigate to="/dashboard" replace />;
  }
  return <>{children}</>;
}
