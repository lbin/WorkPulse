import React from "react";
import { Navigate } from "react-router-dom";
import AppLayout from "./components/AppLayout";
import { RequireAuth } from "./components/RequireAuth";

import Dashboard from "./pages/Dashboard";
import OKRCycles from "./pages/OKR/OKRCycles";
import OKRCycleDetail from "./pages/OKR/OKRCycleDetail";
import KRDetail from "./pages/OKR/KRDetail";
import Projects from "./pages/Projects/Projects";
import ProjectDetail from "./pages/Projects/ProjectDetail";
import Meetings from "./pages/Meetings/Meetings";
import MeetingDetail from "./pages/Meetings/MeetingDetail";
import Reports from "./pages/Reports/Reports";
import ReportEdit from "./pages/Reports/ReportEdit";
import Analytics from "./pages/Analytics/Analytics";
import RequirePermission from "./components/RequirePermission";
import Unauthorized from "./pages/Unauthorized";
import Docs from "./pages/Docs/Docs";
import { FeatureFlagGate } from "./components/FeatureFlagGate";
import Login from "./pages/Login";
import Onboarding from "./pages/Onboarding";

export const routes = [
  { path: "/login", element: <Login /> },
  { path: "/onboarding", element: <RequireAuth allowWithoutTeam><Onboarding /></RequireAuth> },
  { path: "/", element: <Navigate to="/dashboard" replace /> },
  {
    path: "/",
    element: (
      <RequireAuth>
        <AppLayout />
      </RequireAuth>
    ),
    children: [
      { path: "dashboard", element: <RequirePermission permission="dashboard.view"><Dashboard /></RequirePermission> },

      { path: "okr", element: <RequirePermission permission="okr.view"><OKRCycles /></RequirePermission> },
      { path: "okr/cycles/:id", element: <RequirePermission permission="okr.view"><OKRCycleDetail /></RequirePermission> },
      { path: "okr/krs/:id", element: <RequirePermission permission="okr.view"><KRDetail /></RequirePermission> },

      { path: "projects", element: <RequirePermission permission="projects.view"><Projects /></RequirePermission> },
      { path: "projects/:id", element: <RequirePermission permission="projects.view"><ProjectDetail /></RequirePermission> },

      { path: "meetings", element: <RequirePermission permission="meetings.view"><Meetings /></RequirePermission> },
      { path: "meetings/:id", element: <RequirePermission permission="meetings.view"><MeetingDetail /></RequirePermission> },

      { path: "reports", element: <RequirePermission permission="reports.view"><Reports /></RequirePermission> },
      { path: "reports/:id/edit", element: <RequirePermission permission="reports.manage"><ReportEdit /></RequirePermission> },

      {
        path: "analytics",
        element: (
          <RequirePermission permission="analytics.view">
            <FeatureFlagGate flag="analytics">
              <Analytics />
            </FeatureFlagGate>
          </RequirePermission>
        ),
      },
      {
        path: "docs",
        element: (
          <RequirePermission permission="dashboard.view">
            <FeatureFlagGate flag="apiDocs">
              <Docs />
            </FeatureFlagGate>
          </RequirePermission>
        ),
      },
    ],
  },
  { path: "/unauthorized", element: <Unauthorized /> },
];
