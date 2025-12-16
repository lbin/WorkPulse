import { api } from "./client";

export interface AuthResponse {
  token: string;
  user: {
    id: string;
    org_id: string;
    email: string;
    display_name: string;
  };
}

export const login = async (email: string, password: string) => {
  const res = await api.post<AuthResponse>("/auth/login", { email, password });
  return res.data;
};

export const register = async (
  email: string,
  password: string,
  displayName?: string,
  orgName?: string
) => {
  const res = await api.post<AuthResponse>("/auth/register", {
    email,
    password,
    display_name: displayName,
    org_name: orgName,
  });
  return res.data;
};
