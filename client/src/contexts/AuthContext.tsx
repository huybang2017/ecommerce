"use client";

import React, { createContext, useContext, ReactNode } from "react";
import { usePathname } from "next/navigation";
import {
  useProfile,
  useLogin,
  useLogout,
  useRegister,
  User,
  LoginRequest,
  RegisterRequest,
} from "@/hooks/useAuth";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
  refreshUser: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const pathname = usePathname();

  // CRITICAL: Don't fetch profile on public pages to prevent refresh loop
  const isPublicPage = pathname === "/login" || pathname === "/register";
  const shouldFetchProfile = !isPublicPage;

  // Use React Query hooks - only fetch if authenticated
  const {
    data: user,
    isLoading: loading,
    refetch,
  } = useProfile(shouldFetchProfile);

  const loginMutation = useLogin();
  const registerMutation = useRegister();
  const logoutMutation = useLogout();

  const login = async (credentials: LoginRequest) => {
    await loginMutation.mutateAsync(credentials);
    // User will be automatically set in cache by mutation onSuccess
  };

  const register = async (data: RegisterRequest) => {
    await registerMutation.mutateAsync(data);
    // User will be automatically set in cache by mutation onSuccess
  };

  const logout = async () => {
    await logoutMutation.mutateAsync();
    // Cache will be cleared by mutation onSuccess
  };

  const refreshUser = () => {
    refetch();
  };

  const value: AuthContextType = {
    user: user || null,
    loading,
    login,
    register,
    logout,
    isAuthenticated: !!user,
    refreshUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
