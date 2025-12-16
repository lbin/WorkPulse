import React from "react";
import { Navigate } from "react-router-dom";
import AppLayout from "./components/AppLayout";
import AuthGuard from "./components/AuthGuard";
import TeamGuard from "./components/TeamGuard";
import AuthPage from "./pages/Auth/AuthPage";
import TeamSetup from "./pages/Teams/TeamSetup";

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

export const routes = [
  { path: "/", element: <Navigate to="/dashboard" replace /> },
  { path: "/login", element: <AuthPage /> },
  {
    path: "/",
    element: <AuthGuard />,
    children: [
      { path: "team-setup", element: <TeamSetup /> },
      {
        path: "/",
        element: <TeamGuard />,
        children: [
          {
            path: "/",
            element: <AppLayout />,
            children: [
              { path: "dashboard", element: <Dashboard /> },

              { path: "okr", element: <OKRCycles /> },
              { path: "okr/cycles/:id", element: <OKRCycleDetail /> },
              { path: "okr/krs/:id", element: <KRDetail /> },

              { path: "projects", element: <Projects /> },
              { path: "projects/:id", element: <ProjectDetail /> },

              { path: "meetings", element: <Meetings /> },
              { path: "meetings/:id", element: <MeetingDetail /> },

              { path: "reports", element: <Reports /> },
              { path: "reports/:id/edit", element: <ReportEdit /> },

              { path: "analytics", element: <Analytics /> },
            ],
          },
        ],
      },
    ],
  },
];
