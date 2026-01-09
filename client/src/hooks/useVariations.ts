import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import type { Variation } from "@/types/variation";
import { getProductVariations } from "@/services/variation.service";

// ==================== REACT QUERY HOOKS ====================

/**
 * Hook to fetch product variations with options
 * Use this for variation selector UI (Shopee-style)
 */
export const useProductVariations = (
  productId: number,
  enabled: boolean = true
) => {
  return useQuery<Variation[], AxiosError>({
    queryKey: ["product", productId, "variations"],
    queryFn: () => getProductVariations(productId),
    enabled: enabled && productId > 0,
    staleTime: 5 * 60 * 1000, // 5 minutes (variations rarely change)
  });
};
