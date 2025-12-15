import { api } from "./client";

export type WorkItem = {
  id: string;
  type: string;
  title: string;
  status: string;
  priority: number;
  effort_minutes?: number;
};

export async function createWorkItem(payload: any) {
  const res = await api.post("/work-items", payload);
  return res.data.data as WorkItem;
}

export async function addOKRLink(id: string, payload: any) {
  const res = await api.post(`/work-items/${id}/okr-links`, payload);
  return res.data.data;
}
