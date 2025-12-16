import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { createTeam, fetchCurrentUser, AuthContext as ServerAuth, OrgUnit } from "../api/auth";
import { Spin, message } from "antd";

interface AuthState {
  profile?: ServerAuth;
  loading: boolean;
  activeOrgUnit?: string;
  refresh: (orgUnitId?: string) => Promise<void>;
  setActiveOrgUnit: (id?: string) => void;
  hasPermission: (perm: string) => boolean;
  isAuthenticated: boolean;
  logout: () => void;
}

const AuthContext = createContext<AuthState | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [profile, setProfile] = useState<ServerAuth>();
  const [loading, setLoading] = useState(true);
  const [activeOrgUnit, setActiveOrgUnit] = useState<string | undefined>(() => localStorage.getItem("active_org_unit") || undefined);
  const [autoTeamCreated, setAutoTeamCreated] = useState(false);

  const refresh = useCallback(async (orgUnitId?: string) => {
    setLoading(true);
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        setProfile(undefined);
        setLoading(false);
        return;
      }
      const ctx = await fetchCurrentUser(orgUnitId || activeOrgUnit || undefined);
      setProfile(ctx);
      if (ctx.active_org_unit) {
        setActiveOrgUnit(ctx.active_org_unit);
        localStorage.setItem("active_org_unit", ctx.active_org_unit);
      }
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || "Failed to load profile");
      if (err?.response?.status === 401) {
        localStorage.removeItem("token");
        setProfile(undefined);
      }
    } finally {
      setLoading(false);
    }
  }, [activeOrgUnit]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  useEffect(() => {
    const shouldAutocreate =
      !autoTeamCreated &&
      !!profile?.user &&
      (!profile.org_units || profile.org_units.length === 0);

    if (!shouldAutocreate) return;

    setAutoTeamCreated(true);
    const fallbackName = `${profile.user.display_name || "My"} Team`;
    createTeam(fallbackName)
      .then(() => refresh())
      .catch((err: any) => {
        console.error(err);
        message.error(err?.response?.data?.error || "Failed to create team automatically");
      });
  }, [autoTeamCreated, profile, refresh]);

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

  const logout = useCallback(() => {
    localStorage.removeItem("token");
    localStorage.removeItem("active_org_unit");
    setProfile(undefined);
    setActiveOrgUnit(undefined);
  }, []);

  const value: AuthState = {
    profile,
    loading,
    activeOrgUnit,
    refresh,
    setActiveOrgUnit: setActive,
    hasPermission: (perm: string) => permissions.has(perm),
    isAuthenticated: !!profile,
    logout,
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
