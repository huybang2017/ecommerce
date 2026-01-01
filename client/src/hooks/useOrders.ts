import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "@/lib/axios-client";
import { AxiosError } from "axios";

// Types
export interface OrderItem {
  product_id: number;
  product_name: string;
  product_image?: string;
  price: number;
  quantity: number;
  subtotal: number;
}

export interface Order {
  id: number;
  user_id: number;
  order_number: string;
  status:
    | "pending"
    | "confirmed"
    | "processing"
    | "shipped"
    | "delivered"
    | "cancelled";
  items: OrderItem[];
  total_amount: number;
  shipping_address: string;
  payment_method: string;
  payment_status: "pending" | "paid" | "failed" | "refunded";
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateOrderRequest {
  items: Array<{
    product_id: number;
    quantity: number;
  }>;
  shipping_address: string;
  payment_method: string;
  notes?: string;
}

export interface OrdersResponse {
  orders: Order[];
  total: number;
  page: number;
  limit: number;
}

// API functions
const orderApi = {
  getOrders: async (
    page: number = 1,
    limit: number = 10
  ): Promise<OrdersResponse> => {
    const { data } = await apiClient.get<OrdersResponse>("/api/v1/orders", {
      params: { page, limit },
    });
    return data;
  },

  getOrder: async (id: number): Promise<Order> => {
    const { data } = await apiClient.get<Order>(`/api/v1/orders/${id}`);
    return data;
  },

  createOrder: async (request: CreateOrderRequest): Promise<Order> => {
    const { data } = await apiClient.post<Order>("/api/v1/orders", request);
    return data;
  },

  cancelOrder: async (id: number): Promise<Order> => {
    const { data } = await apiClient.post<Order>(`/api/v1/orders/${id}/cancel`);
    return data;
  },
};

// React Query Hooks

export const useOrders = (page: number = 1, limit: number = 10) => {
  return useQuery<OrdersResponse, AxiosError>({
    queryKey: ["orders", page, limit],
    queryFn: () => orderApi.getOrders(page, limit),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const useOrder = (id: number, enabled: boolean = true) => {
  return useQuery<Order, AxiosError>({
    queryKey: ["order", id],
    queryFn: () => orderApi.getOrder(id),
    enabled,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCreateOrder = () => {
  const queryClient = useQueryClient();

  return useMutation<Order, AxiosError, CreateOrderRequest>({
    mutationFn: orderApi.createOrder,
    onSuccess: (data) => {
      // Invalidate orders list to refetch
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      // Clear cart after successful order
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Order created:", data.order_number);
    },
    onError: (error) => {
      console.error("❌ Failed to create order:", error.response?.data);
    },
  });
};

export const useCancelOrder = () => {
  const queryClient = useQueryClient();

  return useMutation<Order, AxiosError, number>({
    mutationFn: orderApi.cancelOrder,
    onSuccess: (data) => {
      // Update order cache
      queryClient.setQueryData(["order", data.id], data);
      // Invalidate orders list
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      console.log("✅ Order cancelled:", data.order_number);
    },
    onError: (error) => {
      console.error("❌ Failed to cancel order:", error.response?.data);
    },
  });
};

// Export orderApi for use outside hooks if needed
export { orderApi };
