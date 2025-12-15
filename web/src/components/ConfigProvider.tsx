import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import { message, Spin } from "antd";
import { ClientConfig, fetchClientConfig } from "../api/config";

interface ConfigState {
  config?: ClientConfig;
  isFeatureEnabled: (flag: string, defaultValue?: boolean) => boolean;
}

const ConfigContext = createContext<ConfigState | undefined>(undefined);

export function ClientConfigProvider({ children }: { children: React.ReactNode }) {
  const [config, setConfig] = useState<ClientConfig>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const cfg = await fetchClientConfig();
        setConfig(cfg);
      } catch (err: any) {
        console.error(err);
        message.warning(err?.response?.data?.error || "Failed to load configuration");
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  const isFeatureEnabled = useMemo(
    () =>
      (flag: string, defaultValue: boolean = true) => {
        return config?.feature_flags?.[flag] ?? defaultValue;
      },
    [config]
  );

  if (loading && !config) {
    return (
      <div style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "60vh" }}>
        <Spin />
      </div>
    );
  }

  return <ConfigContext.Provider value={{ config, isFeatureEnabled }}>{children}</ConfigContext.Provider>;
}

export function useConfig() {
  const ctx = useContext(ConfigContext);
  if (!ctx) throw new Error("useConfig must be used within ConfigProvider");
  return ctx;
}

export function useFeatureFlag(flag: string, defaultValue?: boolean) {
  const { isFeatureEnabled } = useConfig();
  return isFeatureEnabled(flag, defaultValue);
}
