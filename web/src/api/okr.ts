import { api } from "./client";

export type OKRKeyResult = {
  id: string;
  objective_id: string;
  title: string;
  metric_type: string;
  target_value?: number;
  current_value?: number;
  confidence?: number;
  status: string;
};

export type OKRObjective = {
  id: string;
  cycle_id: string;
  team_id?: string;
  owner_user_id?: string;
  title: string;
  description?: string;
  status: string;
  tags?: string[];
};

export type OKRLink = {
  id: string;
  objective_id?: string;
  key_result_id?: string;
  entity_type: string;
  entity_id: string;
  relation: string;
};

export type OKRCycle = {
  id: string;
  name: string;
  type: string;
  team_id?: string;
  start_date?: string;
  end_date?: string;
  status: string;
};

export type OKRCoverage = {
  entity_type: string;
  coverage: number;
  linked_key_results: number;
  total_key_results: number;
  deviation: number;
};

export type OKRMetrics = {
  completion: number;
  deviation: number;
  unlinked_key_results: OKRKeyResult[];
  coverage: OKRCoverage[];
};

export type OKRFilter = {
  cycle_id?: string;
  team_id?: string;
  owner_user_id?: string;
  status?: string;
};

export async function listCycles() {
  const { data } = await api.get("/okrs/cycles");
  return data.data as OKRCycle[];
}

export async function createCycle(payload: {
  name: string;
  type?: string;
  status?: string;
  start_date?: string;
  end_date?: string;
  team_id?: string;
}) {
  const { data } = await api.post("/okrs/cycles", payload);
  return data.data as OKRCycle;
}

export async function updateCycle(id: string, payload: Partial<OKRCycle>) {
  const { data } = await api.put(`/okrs/cycles/${id}`, payload);
  return data.data as OKRCycle;
}

export async function listObjectives(filter: OKRFilter) {
  const { data } = await api.get("/okrs/objectives", { params: filter });
  return data.data as { objectives: OKRObjective[]; key_results: OKRKeyResult[] };
}

export async function createObjective(payload: {
  cycle_id: string;
  title: string;
  description?: string;
  status?: string;
  team_id?: string;
  owner_user_id?: string;
  tags?: string[];
}) {
  const { data } = await api.post("/okrs/objectives", payload);
  return data.data as OKRObjective;
}

export async function updateKRProgress(id: string, current: number, confidence: number) {
  return api.post(`/okrs/key-results/${id}/progress`, { current, confidence });
}

export async function createKeyResult(payload: {
  objective_id: string;
  title: string;
  metric_type?: string;
  target_value?: number;
  unit?: string;
  current_value?: number;
  confidence?: number;
  status?: string;
}) {
  const { data } = await api.post("/okrs/key-results", payload);
  return data.data as OKRKeyResult;
}

export async function addLink(link: Omit<OKRLink, "id">) {
  return api.post("/okrs/links", link);
}

export async function removeLink(id: string) {
  return api.delete(`/okrs/links/${id}`);
}

export async function fetchMetrics(filter: OKRFilter) {
  const { data } = await api.get("/okrs/metrics", { params: filter });
  return data.data as OKRMetrics;
}
