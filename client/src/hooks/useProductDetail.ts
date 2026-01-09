import { useState, useEffect, useCallback, useMemo } from "react";
import { ProductItem } from "@/types/product-item";
import { useProductItems } from "./useProductItems";

/**
 * Custom hook to manage product detail page state
 * Handles SKU selection, quantity, and cart operations
 */
export const useProductDetail = (productId: number) => {
  const [selectedSKU, setSelectedSKU] = useState<ProductItem | null>(null);
  const [quantity, setQuantity] = useState(1);

  // Fetch product items (SKUs)
  const {
    data: productItemsData,
    isLoading: isLoadingSKUs,
    error,
  } = useProductItems(productId);

  // Compute first available SKU
  const firstAvailableSKU = useMemo(() => {
    if (selectedSKU || !productItemsData?.items) return null;

    return (
      productItemsData.items.find(
        (item) => item.status === "ACTIVE" && item.qty_in_stock > 0
      ) || null
    );
  }, [productItemsData, selectedSKU]);

  // Auto-select first available SKU when data loads
  useEffect(() => {
    if (firstAvailableSKU && !selectedSKU) {
      setSelectedSKU(firstAvailableSKU);
    }
  }, [firstAvailableSKU, selectedSKU]);

  // Handle SKU change
  const handleSKUChange = useCallback((sku: ProductItem | null) => {
    setSelectedSKU(sku);
    setQuantity(1); // Reset quantity when SKU changes
  }, []);

  // Handle quantity change
  const handleQuantityChange = useCallback(
    (newQuantity: number) => {
      if (!selectedSKU) return;

      const maxQty = selectedSKU.qty_in_stock;
      const validQty = Math.max(1, Math.min(newQuantity, maxQty));
      setQuantity(validQty);
    },
    [selectedSKU]
  );

  // Validate if can add to cart
  const canAddToCart = useCallback(() => {
    if (!selectedSKU) return false;
    if (selectedSKU.qty_in_stock < quantity) return false;
    if (selectedSKU.status !== "ACTIVE") return false;
    return true;
  }, [selectedSKU, quantity]);

  return {
    // Data
    productItems: productItemsData?.items || [],
    selectedSKU,
    quantity,

    // Loading/Error states
    isLoadingSKUs,
    error,

    // Actions
    setSelectedSKU: handleSKUChange,
    setQuantity: handleQuantityChange,
    canAddToCart: canAddToCart(),
  };
};
