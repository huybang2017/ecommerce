import { useQuery } from "@tanstack/react-query";
import {
  getCategoryById,
  getCategoryProducts,
} from "@/services/category.service";
import type { ProductFilters } from "@/types/category";

/**
 * React Query hook for fetching category detail
 */
export function useCategoryDetail(categoryId: number) {
  return useQuery({
    queryKey: ["category", categoryId],
    queryFn: () => getCategoryById(categoryId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: !!categoryId,
  });
}

/**
 * React Query hook for fetching category products
 */
export function useCategoryProducts(
  categoryId: number,
  filters: ProductFilters = {}
) {
  return useQuery({
    queryKey: ["category-products", categoryId, filters],
    queryFn: () => getCategoryProducts(categoryId, filters),
    staleTime: 1 * 60 * 1000, // 1 minute
    enabled: !!categoryId,
  });
}
