import { apiClient } from "@/lib/axios-client";
import type { Product } from "@/types/product";

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
  total_pages: number;
}

/**
 * Get list of products with filters
 */
export async function getProducts(
  params: GetProductsParams = {}
): Promise<ProductsResponse> {
  const { data } = await apiClient.get<ProductsResponse>("/api/v1/products", {
    params,
  });
  return data;
}

/**
 * Get single product by ID
 */
export async function getProduct(id: number): Promise<Product> {
  const { data } = await apiClient.get<Product>(`/api/v1/products/${id}`);
  return data;
}
