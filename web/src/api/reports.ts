import { api } from "./client";

export type ReportLink = {
  id?: string;
  target_type: string;
  target_id: string;
  relation?: string;
};

export type Report = {
  id: string;
  org_id: string;
  author_user_id: string;
  reviewer_user_id?: string;
  type: string;
  period_start: string;
  period_end: string;
  status: string;
  title?: string;
  summary?: string;
  content?: any;
  links?: ReportLink[];
  review_comment?: string;
  submitted_at?: string;
  reviewed_at?: string;
  created_at?: string;
  updated_at?: string;
};

export type ReportFilters = {
  type?: string;
  status?: string;
  author_user_id?: string;
  from?: string;
  to?: string;
  page?: number;
  page_size?: number;
};

export async function listReports(params: ReportFilters) {
  const res = await api.get("/reports", { params });
  return res.data.data as { items: Report[]; total: number };
}

export async function getReport(id: string) {
  const res = await api.get(`/reports/${id}`);
  return res.data.data as { report: Report; links: ReportLink[] };
}

export async function createReport(payload: any) {
  const res = await api.post("/reports", payload);
  return res.data.data as { report: Report; links: ReportLink[] };
}

export async function updateReport(id: string, payload: any) {
  const res = await api.put(`/reports/${id}`, payload);
  return res.data.data as { report: Report; links: ReportLink[] };
}

export async function changeReportStatus(id: string, action: "submit" | "approve" | "reject", comment?: string) {
  const res = await api.post(`/reports/${id}/${action}`, { comment });
  return res.data.data as Report;
}

export async function exportReports(params: ReportFilters) {
  const res = await api.get("/reports/export", { params, responseType: "blob" });
  return res.data as Blob;
}
