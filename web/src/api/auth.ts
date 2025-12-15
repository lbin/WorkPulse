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

export interface AuthContext {
  user: UserProfile;
  org_units: OrgUnit[];
  roles: Role[];
  permissions: string[];
  active_org_unit?: string;
}

export async function fetchCurrentUser(orgUnitId?: string): Promise<AuthContext> {
  const res = await api.get<AuthContext>("/auth/me", {
    headers: orgUnitId ? { "X-Org-Unit-ID": orgUnitId } : undefined,
  });
  return res.data;
}
