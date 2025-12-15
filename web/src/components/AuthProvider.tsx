import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { fetchCurrentUser, AuthContext as ServerAuth, OrgUnit } from "../api/auth";
import { Spin, message } from "antd";

interface AuthState {
  profile?: ServerAuth;
  loading: boolean;
  activeOrgUnit?: string;
  refresh: (orgUnitId?: string) => Promise<void>;
  setActiveOrgUnit: (id?: string) => void;
  hasPermission: (perm: string) => boolean;
}

const AuthContext = createContext<AuthState | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [profile, setProfile] = useState<ServerAuth>();
  const [loading, setLoading] = useState(true);
  const [activeOrgUnit, setActiveOrgUnit] = useState<string | undefined>(() => localStorage.getItem("active_org_unit") || undefined);

  const refresh = useCallback(async (orgUnitId?: string) => {
    setLoading(true);
    try {
      const ctx = await fetchCurrentUser(orgUnitId || activeOrgUnit || undefined);
      setProfile(ctx);
      if (ctx.active_org_unit) {
        setActiveOrgUnit(ctx.active_org_unit);
        localStorage.setItem("active_org_unit", ctx.active_org_unit);
      }
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || "Failed to load profile");
    } finally {
      setLoading(false);
    }
  }, [activeOrgUnit]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const setActive = useCallback(
    (id?: string) => {
      setActiveOrgUnit(id);
      if (id) {
        localStorage.setItem("active_org_unit", id);
      } else {
        localStorage.removeItem("active_org_unit");
      }
      refresh(id);
    },
    [refresh]
  );

  const permissions = useMemo(() => new Set(profile?.permissions || []), [profile]);

  const value: AuthState = {
    profile,
    loading,
    activeOrgUnit,
    refresh,
    setActiveOrgUnit: setActive,
    hasPermission: (perm: string) => permissions.has(perm),
  };

  if (loading && !profile) {
    return (
      <div style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "60vh" }}>
        <Spin />
      </div>
    );
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}

export function usePermission(code: string) {
  const { hasPermission } = useAuth();
  return hasPermission(code);
}

export function useOrgUnits(): OrgUnit[] {
  const { profile } = useAuth();
  return profile?.org_units || [];
}
