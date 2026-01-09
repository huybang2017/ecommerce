import { apiClient } from "@/lib/axios-client";
import type { ProductItemsResponse, ProductItem } from "@/types/product-item";

// ==================== PRODUCT ITEMS (SKU) ====================

/**
 * Get all SKUs (product items) for a product
 */
export async function getProductItems(
  productId: number
): Promise<ProductItemsResponse> {
  const { data } = await apiClient.get<ProductItemsResponse>(
    `/api/v1/products/${productId}/items`
  );
  return data;
}

/**
 * Get specific SKU by product ID and item ID
 */
export async function getProductItem(
  productId: number,
  itemId: number
): Promise<ProductItem> {
  const { data } = await apiClient.get<ProductItem>(
    `/api/v1/products/${productId}/items/${itemId}`
  );
  return data;
}

/**
 * Get SKU by SKU code or item ID directly
 */
export async function getProductItemBySKU(
  skuCodeOrId: string | number
): Promise<ProductItem> {
  const { data } = await apiClient.get<ProductItem>(
    `/api/v1/product-items/${skuCodeOrId}`
  );
  return data;
}

/**
 * Batch fetch multiple product items
 */
export async function getProductItemsBatch(
  itemIds: number[]
): Promise<ProductItem[]> {
  const { data } = await apiClient.get<ProductItem[]>(
    `/api/v1/product-items/batch`,
    {
      params: { ids: itemIds.join(",") },
    }
  );
  return data;
}
