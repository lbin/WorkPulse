import { api } from "./client";

export interface Team {
  id: string;
  name: string;
  path: string;
}

export interface Membership {
  id: string;
  team_id: string;
  user_id: string;
  role_in_team: string;
}

export const listOrgTeams = async () => {
  const res = await api.get<{ data: Team[] }>("/teams");
  return res.data.data;
};

export const listMyTeams = async () => {
  const res = await api.get<{ data: Membership[] }>("/teams/my");
  return res.data.data;
};

export const createTeam = async (name: string) => {
  const res = await api.post<{ data: Team }>("/teams", { name });
  return res.data.data;
};

export const joinTeam = async (teamId: string) => {
  await api.post("/teams/join", { team_id: teamId });
};
