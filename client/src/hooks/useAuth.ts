import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import apiClient, { setAccessToken } from "@/lib/axios-client";
import { AxiosError } from "axios";

// Types
export interface User {
  id: number;
  username: string;
  email: string;
  phone_number?: string;
  full_name?: string;
  avatar_url?: string;
  role: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  full_name: string;
  phone_number?: string;
}

export interface AuthResponse {
  message: string;
  access_token: string; // NOW returned in response body
  user: User;
}

export interface ErrorResponse {
  error: string;
  message?: string;
}

// API functions
const authApi = {
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>(
      "/api/v1/auth/login",
      credentials
    );
    // Store access_token in memory (refresh_token is in HttpOnly cookie)
    setAccessToken(data.access_token);
    return data;
  },

  register: async (userData: RegisterRequest): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>(
      "/api/v1/auth/register",
      userData
    );
    // Store access_token in memory (refresh_token is in HttpOnly cookie)
    setAccessToken(data.access_token);
    return data;
  },

  logout: async (): Promise<{ message: string }> => {
    const { data } = await apiClient.post<{ message: string }>(
      "/api/v1/auth/logout"
    );
    // Clear access_token from memory (refresh_token cookie cleared by server)
    setAccessToken(null);
    return data;
  },

  refreshToken: async (): Promise<void> => {
    await apiClient.post("/api/v1/auth/refresh");
  },

  getProfile: async (): Promise<User> => {
    const { data } = await apiClient.get<User>("/api/v1/users/profile");
    return data;
  },
};

// React Query hooks
export const useLogin = () => {
  const queryClient = useQueryClient();

  return useMutation<AuthResponse, AxiosError<ErrorResponse>, LoginRequest>({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      // Cache user profile
      queryClient.setQueryData(["user"], data.user);
      console.log("✅ Login successful:", data.user.email);
      console.log("🔑 Access token stored in memory");
    },
    onError: (error) => {
      console.error(
        "❌ Login failed:",
        error.response?.data?.error || error.message
      );
    },
  });
};

export const useRegister = () => {
  const queryClient = useQueryClient();

  return useMutation<AuthResponse, AxiosError<ErrorResponse>, RegisterRequest>({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      // Cache user profile
      queryClient.setQueryData(["user"], data.user);
      console.log("✅ Register successful:", data.user.email);
      console.log("🔑 Access token stored in memory");
    },
    onError: (error) => {
      console.error(
        "❌ Register failed:",
        error.response?.data?.error || error.message
      );
    },
  });
};

export const useLogout = () => {
  const queryClient = useQueryClient();

  return useMutation<{ message: string }, AxiosError<ErrorResponse>, void>({
    mutationFn: authApi.logout,
    onSuccess: () => {
      // Clear all cached data
      queryClient.clear();
      console.log("✅ Logout successful");
      console.log("🔑 Access token cleared from memory");
    },
    onError: (error) => {
      console.error(
        "❌ Logout failed:",
        error.response?.data?.error || error.message
      );
      // Clear cache and token anyway on logout failure
      setAccessToken(null);
      queryClient.clear();
    },
  });
};

export const useProfile = (enabled: boolean = true) => {
  return useQuery<User, AxiosError<ErrorResponse>>({
    queryKey: ["user"],
    queryFn: authApi.getProfile,
    enabled,
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Export authApi for use in AuthContext if needed
export { authApi };
