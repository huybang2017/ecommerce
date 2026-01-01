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
  addItem: (productId: number, quantity: number) => Promise<void>;
  updateItem: (productId: number, quantity: number) => Promise<void>;
  removeItem: (productId: number) => Promise<void>;
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
  const itemCount = cart
    ? cart.items.reduce((sum: number, item) => sum + item.quantity, 0)
    : 0;
  const total = cart?.total_price || 0;

  // Refresh cart
  const refreshCart = () => {
    refetch();
  };

  // Add item to cart
  const addItem = async (productId: number, quantity: number) => {
    await addToCartMutation.mutateAsync({ product_id: productId, quantity });
  };

  // Update cart item
  const updateItem = async (productId: number, quantity: number) => {
    await updateCartMutation.mutateAsync({ productId, quantity });
  };

  // Remove item from cart
  const removeItem = async (productId: number) => {
    await removeCartMutation.mutateAsync(productId);
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
