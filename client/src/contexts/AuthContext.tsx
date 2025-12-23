'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, LoginRequest, RegisterRequest, AuthResponse } from '@/lib/auth-types';
import { authApi } from '@/lib/auth-api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  loading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  // Load token from localStorage on mount
  useEffect(() => {
    const storedToken = localStorage.getItem('auth_token');
    if (storedToken) {
      setToken(storedToken);
      // Try to fetch user profile
      authApi
        .getProfile(storedToken)
        .then((userData) => {
          setUser(userData);
        })
        .catch(() => {
          // Token might be invalid, clear it
          localStorage.removeItem('auth_token');
          setToken(null);
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (credentials: LoginRequest) => {
    const response: AuthResponse = await authApi.login(credentials);
    if (response.token) {
      setToken(response.token);
      localStorage.setItem('auth_token', response.token);
      
      // Fetch user profile
      if (response.user) {
        setUser(response.user);
      } else {
        const userData = await authApi.getProfile(response.token);
        setUser(userData);
      }
    } else {
      throw new Error(response.message || 'Login failed');
    }
  };

  const register = async (data: RegisterRequest) => {
    const response: AuthResponse = await authApi.register(data);
    if (response.token) {
      setToken(response.token);
      localStorage.setItem('auth_token', response.token);
      
      // Fetch user profile
      if (response.user) {
        setUser(response.user);
      } else {
        const userData = await authApi.getProfile(response.token);
        setUser(userData);
      }
    } else {
      throw new Error(response.message || 'Registration failed');
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('auth_token');
  };

  const value: AuthContextType = {
    user,
    token,
    loading,
    login,
    register,
    logout,
    isAuthenticated: !!token && !!user,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

