import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import type { ProductItemsResponse, ProductItem } from "@/types/product-item";
import {
  getProductItems,
  getProductItem,
  getProductItemBySKU,
  getProductItemsBatch,
} from "@/services/product-item.service";

// ==================== REACT QUERY HOOKS ====================

/**
 * Hook to fetch all product items (SKUs) for a product
 * Use this to get available variations and their prices/stock
 */
export const useProductItems = (productId: number, enabled: boolean = true) => {
  return useQuery<ProductItemsResponse, AxiosError>({
    queryKey: ["product", productId, "items"],
    queryFn: () => getProductItems(productId),
    enabled: enabled && productId > 0,
    staleTime: 2 * 60 * 1000, // 2 minutes (stock changes frequently)
  });
};

/**
 * Hook to fetch a specific product item
 */
export const useProductItem = (
  productId: number,
  itemId: number,
  enabled: boolean = true
) => {
  return useQuery<ProductItem, AxiosError>({
    queryKey: ["product", productId, "item", itemId],
    queryFn: () => getProductItem(productId, itemId),
    enabled: enabled && productId > 0 && itemId > 0,
    staleTime: 2 * 60 * 1000,
  });
};

/**
 * Hook to fetch product item by SKU code
 */
export const useProductItemBySKU = (
  skuCode: string,
  enabled: boolean = true
) => {
  return useQuery<ProductItem, AxiosError>({
    queryKey: ["product-item", "sku", skuCode],
    queryFn: () => getProductItemBySKU(skuCode),
    enabled: enabled && !!skuCode,
    staleTime: 2 * 60 * 1000,
  });
};

/**
 * Hook to batch fetch multiple product items
 */
export const useProductItemsBatch = (
  itemIds: number[],
  enabled: boolean = true
) => {
  return useQuery<ProductItem[], AxiosError>({
    queryKey: ["product-items", "batch", itemIds],
    queryFn: () => getProductItemsBatch(itemIds),
    enabled: enabled && itemIds.length > 0,
    staleTime: 2 * 60 * 1000,
  });
};
