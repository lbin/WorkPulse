import { api } from "./client";

export interface DocsConfig {
  openapi: string;
  graphql_sdl: string;
}

export interface ClientConfig {
  api_version: string;
  feature_flags: Record<string, boolean>;
  docs: DocsConfig;
}

export async function fetchClientConfig(): Promise<ClientConfig> {
  const res = await api.get<ClientConfig>("/config");
  return res.data;
}
