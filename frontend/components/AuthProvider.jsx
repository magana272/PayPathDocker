"use client";

import { createContext, useContext, useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { getToken, setToken, clearToken } from "@/lib/auth";

const AuthContext = createContext(null);

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}

export default function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const token = getToken();
    if (!token) {
      setLoading(false);
      return;
    }
    api.getMe()
      .then(setUser)
      .catch(() => clearToken())
      .finally(() => setLoading(false));
  }, []);

  const login = useCallback(async (email, password) => {
    const data = await api.login(email, password);
    setToken(data.token);
    const me = await api.getMe();
    setUser(me);
    router.push("/");
  }, [router]);

  const register = useCallback(async (email, password, name) => {
    const data = await api.register(email, password, name);
    setToken(data.token);
    const me = await api.getMe();
    setUser(me);
    router.push("/setup");
  }, [router]);

  const logout = useCallback(async () => {
    try { await api.logout(); } catch {}
    clearToken();
    setUser(null);
    router.push("/login");
  }, [router]);

  const deleteAccount = useCallback(async () => {
    await api.deleteAccount();
    clearToken();
    setUser(null);
    router.push("/login");
  }, [router]);

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout, deleteAccount }}>
      {children}
    </AuthContext.Provider>
  );
}
