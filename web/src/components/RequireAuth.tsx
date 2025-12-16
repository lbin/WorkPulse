import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { Spin } from "antd";
import { useAuth } from "./AuthProvider";

interface RequireAuthProps {
  children: React.ReactNode;
  allowWithoutTeam?: boolean;
}

// RequireAuth gates routes that need authentication and a selected team/org unit.
export function RequireAuth({ children, allowWithoutTeam }: RequireAuthProps) {
  const { loading, profile, isAuthenticated } = useAuth();
  const location = useLocation();

  if (loading) {
    return (
      <div style={{ minHeight: "60vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
        <Spin />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  if (!allowWithoutTeam && (!profile?.org_units || profile.org_units.length === 0)) {
    return <Navigate to="/onboarding" replace />;
  }

  return <>{children}</>;
}
