import React from "react";
import Unauthorized from "../pages/Unauthorized";
import { usePermission } from "./AuthProvider";

export default function RequirePermission({ permission, children }: { permission: string; children: React.ReactNode }) {
  const allowed = usePermission(permission);
  if (!allowed) {
    return <Unauthorized />;
  }
  return <>{children}</>;
}
