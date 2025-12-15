import axios from "axios";

export const api = axios.create({
  baseURL: "/api/v1",
  timeout: 15000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  const orgUnit = localStorage.getItem("active_org_unit");
  if (orgUnit) {
    config.headers["X-Org-Unit-ID"] = orgUnit;
  }
  return config;
});
