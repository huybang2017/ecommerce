"use client";

import React, { createContext, useContext } from "react";
import { Cart } from "@/lib/types";
import {
  useCart as useCartQuery,
  useAddToCart,
  useUpdateCartItem,
  useRemoveFromCart,
  useClearCart,
} from "@/hooks/useCart";
import { useAuth } from "./AuthContext";

interface CartContextType {
  cart: Cart | null;
  loading: boolean;
  error: string | null;
  itemCount: number;
  total: number;
  refreshCart: () => void;
  addItem: (productItemId: number, quantity: number) => Promise<void>;
  updateItem: (productItemId: number, quantity: number) => Promise<void>;
  removeItem: (productItemId: number) => Promise<void>;
  clear: () => Promise<void>;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();

  // Use React Query hooks (renamed to avoid conflict)
  const {
    data: cart,
    isLoading: loading,
    error: queryError,
    refetch,
  } = useCartQuery(isAuthenticated);
  const addToCartMutation = useAddToCart();
  const updateCartMutation = useUpdateCartItem();
  const removeCartMutation = useRemoveFromCart();
  const clearCartMutation = useClearCart();

  // Calculate item count and total
  const itemCount = cart?.items
    ? cart.items.reduce((sum: number, item) => sum + item.quantity, 0)
    : 0;
  const total = cart?.total_price || 0;

  // Refresh cart
  const refreshCart = () => {
    refetch();
  };

  // Add item to cart
  const addItem = async (productItemId: number, quantity: number) => {
    await addToCartMutation.mutateAsync({
      product_item_id: productItemId,
      quantity,
    });
  };

  // Update cart item
  const updateItem = async (productItemId: number, quantity: number) => {
    console.log("Updating item", productItemId, quantity);
    await updateCartMutation.mutateAsync({ productItemId, quantity });
  };

  // Remove item from cart
  const removeItem = async (productItemId: number) => {
    await removeCartMutation.mutateAsync(productItemId);
  };

  // Clear cart
  const clear = async () => {
    await clearCartMutation.mutateAsync();
  };

  const value: CartContextType = {
    cart: cart || null,
    loading,
    error: queryError ? String(queryError) : null,
    itemCount,
    total,
    refreshCart,
    addItem,
    updateItem,
    removeItem,
    clear,
  };

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCartContext() {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error("useCartContext must be used within a CartProvider");
  }
  return context;
}
