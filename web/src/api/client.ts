import axios from "axios";

const apiBase = import.meta.env.VITE_API_BASE || "/api";
const apiVersion = import.meta.env.VITE_API_VERSION || "v1";

export const api = axios.create({
  baseURL: `${apiBase}/${apiVersion}`,
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
