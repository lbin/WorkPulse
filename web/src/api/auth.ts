import { api } from "./client";

export interface OrgUnit {
  id: string;
  name: string;
  path: string;
  unit_type: string;
}

export interface Role {
  id: string;
  name: string;
  org_unit_id?: string;
}

export interface UserProfile {
  id: string;
  email: string;
  display_name: string;
  avatar_url?: string;
}

export interface TeamSummary {
  id: string;
  name: string;
  path: string;
  parent_team_id?: string;
}

export interface AuthContext {
  user: UserProfile;
  org_units: OrgUnit[];
  roles: Role[];
  permissions: string[];
  active_org_unit?: string;
}

export interface LoginResponse {
  token: string;
  context: AuthContext;
}

export async function fetchCurrentUser(orgUnitId?: string): Promise<AuthContext> {
  const res = await api.get<AuthContext>("/auth/me", {
    headers: orgUnitId ? { "X-Org-Unit-ID": orgUnitId } : undefined,
  });
  return res.data;
}

// login exchanges email/password for a JWT.
export async function login(email: string, password: string): Promise<LoginResponse> {
  const res = await api.post<LoginResponse>("/auth/login", { email, password });
  return res.data;
}

// register provisions a new account using email/password.
export async function register(email: string, password: string, displayName?: string): Promise<LoginResponse> {
  const res = await api.post<LoginResponse>("/auth/register", { email, password, display_name: displayName });
  return res.data;
}

// createTeam builds a new team/org unit for the current org.
export async function createTeam(name: string, parent_team_id?: string): Promise<TeamSummary> {
  const res = await api.post<TeamSummary>("/auth/teams", { name, parent_team_id });
  return res.data;
}

// joinTeam subscribes the caller to an existing team.
export async function joinTeam(team_id: string): Promise<AuthContext> {
  const res = await api.post<AuthContext>("/auth/teams/join", { team_id });
  return res.data;
}

// listMyTeams returns memberships for the current user.
export async function listMyTeams(): Promise<TeamSummary[]> {
  const res = await api.get<TeamSummary[]>("/auth/teams");
  return res.data;
}
