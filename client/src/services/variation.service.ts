import { apiClient } from "@/lib/axios-client";
import type { Variation } from "@/types/variation";

// ==================== VARIATIONS ====================

/**
 * Get variations with options for a product
 * Used for variation selector UI (Color, Size, etc.)
 */
export async function getProductVariations(
  productId: number
): Promise<Variation[]> {
  const { data } = await apiClient.get<Variation[]>(
    `/api/v1/products/${productId}/variations`
  );
  return data;
}
