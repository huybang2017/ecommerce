import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

// Create axios instance with cookie-based authentication
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  withCredentials: true, // CRITICAL: Send all cookies (access_token, refresh_token, session_id)
  headers: {
    "Content-Type": "application/json",
  },
});

// Flag to prevent infinite refresh loops
let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

// Request interceptor - just log requests (no Authorization header needed)
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    console.log(`🚀 [API] ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error("❌ [API] Request error:", error);
    return Promise.reject(error);
  }
);

// Response interceptor with automatic token refresh
apiClient.interceptors.response.use(
  (response) => {
    console.log(
      `✅ [API] ${response.config.method?.toUpperCase()} ${
        response.config.url
      } - ${response.status}`
    );
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // If 401 and not already retrying
    if (error.response?.status === 401 && !originalRequest._retry) {
      // Don't retry login/register/refresh endpoints
      if (
        originalRequest.url?.includes("/auth/login") ||
        originalRequest.url?.includes("/auth/register") ||
        originalRequest.url?.includes("/auth/refresh")
      ) {
        console.error("❌ [API] Auth endpoint failed:", originalRequest.url);
        return Promise.reject(error);
      }

      // If already refreshing, wait for that refresh to complete
      if (isRefreshing && refreshPromise) {
        try {
          await refreshPromise;
          originalRequest._retry = true;
          return apiClient(originalRequest);
        } catch (refreshError) {
          return Promise.reject(refreshError);
        }
      }

      // Mark as retrying to prevent infinite loops
      originalRequest._retry = true;
      isRefreshing = true;

      // Create refresh promise
      refreshPromise = (async (): Promise<void> => {
        try {
          console.log("🔄 [API] Attempting token refresh...");
          await apiClient.post("/api/v1/auth/refresh");
          console.log("✅ [API] Token refreshed successfully");

          isRefreshing = false;
          refreshPromise = null;
        } catch (refreshError) {
          console.error("❌ [API] Token refresh failed");
          isRefreshing = false;
          refreshPromise = null;

          // Redirect to login on refresh failure
          if (typeof window !== "undefined") {
            window.location.href = "/login";
          }
          throw refreshError;
        }
      })();

      try {
        await refreshPromise;
        // Retry original request (new access_token cookie set automatically)
        return apiClient(originalRequest);
      } catch (refreshError) {
        return Promise.reject(refreshError);
      }
    }

    console.error(
      `❌ [API] ${error.config?.method?.toUpperCase()} ${error.config?.url} - ${
        error.response?.status || "Network Error"
      }`
    );
    return Promise.reject(error);
  }
);

export default apiClient;
