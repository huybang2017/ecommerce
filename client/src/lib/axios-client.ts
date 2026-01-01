import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

// In-memory token storage (NOT localStorage for security)
let accessToken: string | null = null;

// Getter and setter for access token
export const getAccessToken = () => accessToken;
export const setAccessToken = (token: string | null) => {
  accessToken = token;
};

// Create axios instance
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  withCredentials: true, // CRITICAL: Send refresh_token cookie automatically
  headers: {
    "Content-Type": "application/json",
  },
});

// Flag to prevent infinite refresh loops
let isRefreshing = false;
let refreshPromise: Promise<string> | null = null;

// Request interceptor - Add Bearer token from memory
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Add access_token from memory to Authorization header
    const token = getAccessToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
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
          const newToken = await refreshPromise;
          // Update request with new token
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${newToken}`;
          }
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
      refreshPromise = (async (): Promise<string> => {
        try {
          console.log("🔄 [API] Attempting token refresh...");
          const response = await apiClient.post("/api/v1/auth/refresh");

          // Extract new access_token from response body
          const newAccessToken = response.data.access_token;
          if (!newAccessToken) {
            throw new Error("No access_token in refresh response");
          }

          // Store new token in memory
          setAccessToken(newAccessToken);
          console.log("✅ [API] Token refreshed successfully");

          isRefreshing = false;
          refreshPromise = null;
          return newAccessToken;
        } catch (refreshError) {
          console.error("❌ [API] Token refresh failed");
          isRefreshing = false;
          refreshPromise = null;
          setAccessToken(null); // Clear invalid token

          // Redirect to login on refresh failure
          if (typeof window !== "undefined") {
            window.location.href = "/login";
          }
          throw refreshError;
        }
      })();

      try {
        const newToken = await refreshPromise;
        // Retry original request with new token
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
        }
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
