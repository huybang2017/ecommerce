import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "@/lib/axios-client";
import { AxiosError } from "axios";
import { useAuth } from "@/contexts/AuthContext";

// Types
export interface Address {
  id: number;
  user_id?: number;
  recipient_name: string;
  phone_number: string;
  address_line: string;
  city?: string;
  district?: string;
  ward?: string;
  is_default?: boolean;
  label?: string;
}

export interface CreateAddressRequest {
  recipient_name: string;
  phone_number: string;
  address_line: string;
  city?: string;
  district?: string;
  ward?: string;
  label?: string;
  // is_default usually set via dedicated endpoint, keep optional if supported
  is_default?: boolean;
}

export type UpdateAddressRequest = Partial<CreateAddressRequest>;

// API
const addressApi = {
  list: async (): Promise<Address[]> => {
    const resp = await apiClient.get<any>("/api/v1/addresses");
    const payload = resp.data;
    // Some services wrap the result in { data: [...] }
    if (payload && typeof payload === "object" && "data" in payload) {
      return payload.data as Address[];
    }
    return payload as Address[];
  },

  get: async (id: number): Promise<Address> => {
    const resp = await apiClient.get<any>(`/api/v1/addresses/${id}`);
    const payload = resp.data;
    if (payload && typeof payload === "object" && "data" in payload) {
      return payload.data as Address;
    }
    return payload as Address;
  },

  create: async (payload: CreateAddressRequest): Promise<Address> => {
    const resp = await apiClient.post<any>("/api/v1/addresses", payload);
    const body = resp.data;
    if (body && typeof body === "object" && "data" in body) return body.data as Address;
    return body as Address;
  },

  update: async (
    id: number,
    payload: UpdateAddressRequest,
  ): Promise<Address> => {
    const resp = await apiClient.put<any>(`/api/v1/addresses/${id}`, payload);
    const body = resp.data;
    if (body && typeof body === "object" && "data" in body) return body.data as Address;
    return body as Address;
  },

  remove: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/v1/addresses/${id}`);
  },

  setDefault: async (id: number): Promise<Address> => {
    const resp = await apiClient.put<any>(`/api/v1/addresses/${id}/default`);
    const body = resp.data;
    if (body && typeof body === "object" && "data" in body) return body.data as Address;
    return body as Address;
  },
};

// Hooks
export const useAddresses = (enabled?: boolean) => {
  const { isAuthenticated } = useAuth();
  const finalEnabled =
    enabled === undefined ? isAuthenticated : enabled && isAuthenticated;

  return useQuery<Address[], AxiosError>({
    queryKey: ["addresses"],
    queryFn: addressApi.list,
    enabled: finalEnabled,
    staleTime: 2 * 60 * 1000,
  });
};

export const useAddress = (id?: number, enabled?: boolean) => {
  const { isAuthenticated } = useAuth();
  const finalEnabled =
    (enabled === undefined ? isAuthenticated : enabled) &&
    Boolean(id) &&
    isAuthenticated;

  return useQuery<Address, AxiosError>({
    queryKey: ["addresses", id],
    queryFn: () =>
      id ? addressApi.get(id) : Promise.reject(new Error("Missing id")),
    enabled: finalEnabled,
  });
};

export const useCreateAddress = () => {
  const qc = useQueryClient();
  const { isAuthenticated } = useAuth();

  return useMutation<Address, AxiosError, CreateAddressRequest>({
    mutationFn: async (payload) => {
      if (!isAuthenticated) throw new Error("Not authenticated");
      return addressApi.create(payload);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["addresses"] }),
  });
};

export const useUpdateAddress = () => {
  const qc = useQueryClient();
  const { isAuthenticated } = useAuth();

  return useMutation<
    Address,
    AxiosError,
    { id: number; payload: UpdateAddressRequest }
  >({
    mutationFn: async ({ id, payload }) => {
      if (!isAuthenticated) throw new Error("Not authenticated");
      return addressApi.update(id, payload);
    },
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: ["addresses"] });
      qc.setQueryData(["addresses", data.id], data);
    },
  });
};

export const useDeleteAddress = () => {
  const qc = useQueryClient();
  const { isAuthenticated } = useAuth();

  return useMutation<void, AxiosError, number>({
    mutationFn: async (id) => {
      if (!isAuthenticated) throw new Error("Not authenticated");
      return addressApi.remove(id);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["addresses"] }),
  });
};

export const useSetDefaultAddress = () => {
  const qc = useQueryClient();
  const { isAuthenticated } = useAuth();

  return useMutation<Address, AxiosError, number>({
    mutationFn: async (id) => {
      if (!isAuthenticated) throw new Error("Not authenticated");
      return addressApi.setDefault(id);
    },
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: ["addresses"] });
      qc.setQueryData(["addresses", data.id], data);
    },
  });
};

export { addressApi };
