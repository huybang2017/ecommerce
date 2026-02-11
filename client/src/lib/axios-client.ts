/**
 * ============================================================================
 * AXIOS CLIENT - Unified for Client & Server Components
 * ============================================================================
 *
 * Challenge:
 * - Client Components run in browser → Cookies sent automatically
 * - Server Components run on server → Must forward cookies manually
 *
 * Solution:
 * - clientAxios: For Client Components (withCredentials)
 * - serverAxios: For Server Components (manual cookie forwarding)
 * - api: Smart wrapper that auto-detects environment
 *
 * Usage:
 * - Client Component: import { api } from '@/lib/axios-client'
 * - Server Component: import { api } from '@/lib/axios-client'
 * - Both work seamlessly!
 * ============================================================================
 */

import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios";
import { cookies } from "next/headers";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

// ────────────────────────────────────────────────────────────────────────────
// CLIENT-SIDE AXIOS INSTANCE
// ────────────────────────────────────────────────────────────────────────────

/**
 * For Client Components and Client-side API calls
 * Cookies are sent automatically by browser
 */
export const clientAxios: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  withCredentials: true, // ✅ CRITICAL: Send httpOnly cookies
  headers: {
    "Content-Type": "application/json",
    "X-App-Role": "buyer", // ✅ CRITICAL: Role header
  },
});

// ── Request Interceptor ───────────────────────────────────────────────────
clientAxios.interceptors.request.use(
  (config) => {
    console.log(`🚀 [CLIENT AXIOS] ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error("❌ [CLIENT AXIOS] Request error:", error);
    return Promise.reject(error);
  }
);

// ── Response Interceptor (401 handling) ───────────────────────────────────
let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

clientAxios.interceptors.response.use(
  (response) => {
    console.log(`✅ [CLIENT AXIOS] ${response.status} ${response.config.url}`);
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // ── Handle 401 Unauthorized ────────────────────────────────────────────
    if (error.response?.status === 401 && !originalRequest._retry) {
      // Don't retry auth endpoints
      if (
        originalRequest.url?.includes("/auth/login") ||
        originalRequest.url?.includes("/auth/register") ||
        originalRequest.url?.includes("/auth/refresh")
      ) {
        console.error("❌ [CLIENT AXIOS] Auth endpoint failed");
        return Promise.reject(error);
      }

      // If already refreshing, wait for that refresh to complete
      if (isRefreshing && refreshPromise) {
        try {
          await refreshPromise;
          originalRequest._retry = true;
          return clientAxios(originalRequest);
        } catch (refreshError) {
          return Promise.reject(refreshError);
        }
      }

      // Mark as retrying
      originalRequest._retry = true;
      isRefreshing = true;

      // Create refresh promise
      refreshPromise = (async (): Promise<void> => {
        try {
          console.log("🔄 [CLIENT AXIOS] Refreshing token...");
          await axios.post(
            `${API_BASE_URL}/api/v1/auth/refresh`,
            {},
            {
              withCredentials: true,
              headers: { "X-App-Role": "buyer" },
            }
          );
          console.log("✅ [CLIENT AXIOS] Token refreshed");
          isRefreshing = false;
          refreshPromise = null;
        } catch (refreshError) {
          console.error("❌ [CLIENT AXIOS] Refresh failed");
          isRefreshing = false;
          refreshPromise = null;

          // Redirect to login
          if (typeof window !== "undefined") {
            window.location.href = "/login?reason=session_expired";
          }
          throw refreshError;
        }
      })();

      try {
        await refreshPromise;
        return clientAxios(originalRequest); // Retry with fresh token
      } catch (refreshError) {
        return Promise.reject(refreshError);
      }
    }

    console.error(
      `❌ [CLIENT AXIOS] ${error.response?.status || "Network Error"} ${error.config?.url}`
    );
    return Promise.reject(error);
  }
);

// ────────────────────────────────────────────────────────────────────────────
// SERVER-SIDE AXIOS INSTANCE (For Server Components & SSR)
// ────────────────────────────────────────────────────────────────────────────

/**
 * For Server Components and Server-Side Rendering
 * Cookies must be forwarded manually from request headers
 *
 * Usage:
 * ```typescript
 * // In Server Component
 * import { serverAxios } from '@/lib/axios-client';
 *
 * const { data } = await serverAxios.get('/api/v1/users/profile');
 * ```
 */
export const serverAxios: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
    "X-App-Role": "buyer",
  },
});

// ── Request Interceptor: Forward cookies from Next.js request ─────────────
serverAxios.interceptors.request.use(
  async (config) => {
    console.log(`🚀 [SERVER AXIOS] ${config.method?.toUpperCase()} ${config.url}`);

    // ✅ CRITICAL: Forward cookies from Next.js request to API
    // Server Components don't have access to browser cookies
    // We must read them from Next.js headers and forward manually

    try {
      const cookieStore = cookies(); // Next.js 13+ App Router
      const cookieHeader = cookieStore
        .getAll()
        .map((cookie) => `${cookie.name}=${cookie.value}`)
        .join("; ");

      if (cookieHeader) {
        config.headers["Cookie"] = cookieHeader;
        console.log(
          `[SERVER AXIOS] Forwarding cookies: ${cookieHeader.substring(0, 100)}...`
        );
      } else {
        console.warn("[SERVER AXIOS] No cookies to forward");
      }
    } catch (error) {
      console.error("[SERVER AXIOS] Error reading cookies:", error);
    }

    return config;
  },
  (error) => {
    console.error("❌ [SERVER AXIOS] Request error:", error);
    return Promise.reject(error);
  }
);

// ── Response Interceptor ──────────────────────────────────────────────────
serverAxios.interceptors.response.use(
  (response) => {
    console.log(`✅ [SERVER AXIOS] ${response.status} ${response.config.url}`);
    return response;
  },
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      console.error("[SERVER AXIOS] 401 Unauthorized - Session expired");
      // Server-side can't refresh, must let client handle it
    }

    console.error(
      `❌ [SERVER AXIOS] ${error.response?.status || "Network Error"} ${error.config?.url}`
    );
    return Promise.reject(error);
  }
);

// ────────────────────────────────────────────────────────────────────────────
// SMART AXIOS - Auto-detect environment
// ────────────────────────────────────────────────────────────────────────────

/**
 * Smart axios that chooses correct instance based on environment
 *
 * Usage (works in both Client & Server Components):
 * ```typescript
 * import { api } from '@/lib/axios-client';
 *
 * const { data } = await api.get('/api/v1/users/profile');
 * ```
 *
 * How it works:
 * - typeof window === 'undefined' → Server → Use serverAxios
 * - typeof window !== 'undefined' → Client → Use clientAxios
 */
export const api = typeof window === "undefined" ? serverAxios : clientAxios;

// ────────────────────────────────────────────────────────────────────────────
// EXPORTS
// ────────────────────────────────────────────────────────────────────────────

// Default export (smart API)
export default api;

// Named exports for explicit usage
export { clientAxios as apiClient }; // Backward compatibility

/**
 * ============================================================================
 * USAGE EXAMPLES
 * ============================================================================
 *
 * // ── CLIENT COMPONENT ──────────────────────────────────────────────────
 * 'use client';
 *
 * import { api } from '@/lib/axios-client';
 * import { useEffect, useState } from 'react';
 *
 * export default function ProfileClient() {
 *   const [user, setUser] = useState(null);
 *
 *   useEffect(() => {
 *     api.get('/api/v1/users/profile')
 *       .then(res => setUser(res.data))
 *       .catch(err => console.error(err));
 *   }, []);
 *
 *   return <div>{JSON.stringify(user)}</div>;
 * }
 *
 * // ── SERVER COMPONENT ──────────────────────────────────────────────────
 * import { api } from '@/lib/axios-client';
 *
 * export default async function ProfilePage() {
 *   try {
 *     const { data } = await api.get('/api/v1/users/profile');
 *     return <div>{JSON.stringify(data)}</div>;
 *   } catch (error) {
 *     return <div>Error: {error.message}</div>;
 *   }
 * }
 *
 * ============================================================================
 */
