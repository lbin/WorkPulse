import React from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";

// AuthGuard ensures user has a token before accessing protected routes.
export default function AuthGuard() {
  const token = localStorage.getItem("token");
  const location = useLocation();
  if (!token) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }
  return <Outlet />;
}
