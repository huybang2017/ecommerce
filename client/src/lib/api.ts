/**
 * LEGACY API FILE
 * This file provides backward compatibility for pages still using old API functions.
 * TODO: Refactor all pages to use React Query hooks directly from @/hooks
 */

import apiClient from "./axios-client";
import { Product, Category } from "./types";

// ==================== PRODUCTS ====================

export interface GetProductsParams {
  page?: number;
  limit?: number;
  category_id?: number;
  status?: string;
  sort_by?: string;
  order?: string;
}

export interface ProductsResponse {
  products: Product[];
  total: number;
  page: number;
  limit: number;
}

export async function getProducts(
  params: GetProductsParams = {}
): Promise<ProductsResponse> {
  const { data } = await apiClient.get<ProductsResponse>("/api/v1/products", {
    params,
  });
  return data;
}

export async function getProduct(id: number): Promise<Product> {
  const { data } = await apiClient.get<Product>(`/api/v1/products/${id}`);
  return data;
}

// ==================== CATEGORIES ====================

export async function getCategories(): Promise<Category[]> {
  const { data } = await apiClient.get<Category[]>("/api/v1/categories");
  return data;
}

export async function getCategory(id: number): Promise<Category> {
  const { data } = await apiClient.get<Category>(`/api/v1/categories/${id}`);
  return data;
}

// ==================== SEARCH ====================

export interface SearchParams {
  q?: string;
  category_id?: number;
  min_price?: number;
  max_price?: number;
  page?: number;
  limit?: number;
}

export async function searchProductsAdvanced(
  params: SearchParams
): Promise<ProductsResponse> {
  const { data } = await apiClient.get<ProductsResponse>("/api/v1/search", {
    params,
  });
  return data;
}

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
  items: Array<{
    product_id: number;
    quantity: number;
  }>;
  shipping_address_id?: number;
}

export async function listOrders(): Promise<Order[]> {
  const { data } = await apiClient.get<Order[]>("/api/v1/orders");
  return data;
}

export async function getOrder(id: number): Promise<Order> {
  const { data } = await apiClient.get<Order>(`/api/v1/orders/${id}`);
  return data;
}

export async function createOrder(request: CreateOrderRequest): Promise<Order> {
  const { data } = await apiClient.post<Order>("/api/v1/orders", request);
  return data;
}

export async function cancelOrder(id: number): Promise<Order> {
  const { data } = await apiClient.post<Order>(`/api/v1/orders/${id}/cancel`);
  return data;
}
