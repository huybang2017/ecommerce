import { apiClient } from "@/lib/axios-client";
import type { Product } from "@/types/product";

// ==================== ORDERS ====================

export interface Order {
  id: number;
  user_id: number;
  total_amount: number;
  status: string;
  created_at: string;
  updated_at: string;
  items?: OrderItem[];
}

export interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
  quantity: number;
  price: number;
  product?: Product;
}

export interface CreateOrderRequest {
  // User/session context
  user_id?: number;
  session_id?: string;

  // Shipping information
  shipping_name: string;
  shipping_phone: string;
  shipping_address: string;
  shipping_city: string;
  shipping_province?: string;
  shipping_postal_code?: string;
  shipping_country?: string;
  shipping_address_id?: number;

  // Financials
  shipping_fee?: number;
  shipping_discount?: number;
  voucher_discount?: number;
  tax?: number;
  discount?: number;

  // Payment
  payment_method?: string;

  // Legacy/compat fields
  items?: Array<{
    product_id: number;
    quantity: number;
  }>;
}

/**
 * Get list of orders
 */
export async function listOrders(): Promise<Order[]> {
  const { data } = await apiClient.get<Order[]>("/api/v1/orders");
  return data;
}

/**
 * Get single order by ID
 */
export async function getOrder(id: number): Promise<Order> {
  const { data } = await apiClient.get<Order>(`/api/v1/orders/${id}`);
  return data;
}

/**
 * Create new order
 */
export async function createOrder(request: CreateOrderRequest): Promise<Order> {
  const { data } = await apiClient.post<Order>("/api/v1/orders", request);
  return data;
}

/**
 * Cancel order
 */
export async function cancelOrder(id: number): Promise<Order> {
  const { data } = await apiClient.post<Order>(`/api/v1/orders/${id}/cancel`);
  return data;
}
