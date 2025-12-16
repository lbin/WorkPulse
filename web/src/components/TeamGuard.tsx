import React from "react";
import { Navigate, Outlet } from "react-router-dom";

// TeamGuard redirects users without a selected team to the setup flow.
export default function TeamGuard() {
  const teamId = localStorage.getItem("activeTeamId");
  if (!teamId) {
    return <Navigate to="/team-setup" replace />;
  }
  return <Outlet />;
}
