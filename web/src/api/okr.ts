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

export async function listObjectives(filter: OKRFilter) {
  const { data } = await api.get("/okrs/objectives", { params: filter });
  return data.data as { objectives: OKRObjective[]; key_results: OKRKeyResult[] };
}

export async function updateKRProgress(id: string, current: number, confidence: number) {
  return api.post(`/okrs/key-results/${id}/progress`, { current, confidence });
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
