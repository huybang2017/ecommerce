import { apiClient } from "@/lib/axios-client";
import type { ProductsResponse } from "./product.service";

// ==================== SEARCH ====================

export interface SearchParams {
  q?: string;
  category_id?: number;
  min_price?: number;
  max_price?: number;
  page?: number;
  limit?: number;
}

/**
 * Advanced search for products
 */
export async function searchProductsAdvanced(
  params: SearchParams
): Promise<ProductsResponse> {
  const { data } = await apiClient.get<ProductsResponse>("/api/v1/search", {
    params,
  });
  return data;
}
