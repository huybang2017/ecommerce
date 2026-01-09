import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "@/lib/axios-client";
import { AxiosError } from "axios";

// Types
export interface CartItem {
  product_item_id: number;
  product_name: string;
  product_image?: string;
  name?: string;
  image?: string;
  sku?: string;
  price: number;
  quantity: number;
  subtotal: number;
}

export interface Cart {
  user_id: number;
  items: CartItem[];
  total_items: number;
  total_price: number;
  updated_at: string;
}

export interface AddToCartRequest {
  product_item_id: number;
  quantity: number;
}

export interface UpdateCartRequest {
  quantity: number;
}

// API functions
const cartApi = {
  getCart: async (): Promise<Cart> => {
    const { data } = await apiClient.get<Cart>("/api/v1/cart");
    return data;
  },

  addToCart: async (request: AddToCartRequest): Promise<Cart> => {
    const { data } = await apiClient.post<Cart>("/api/v1/cart/items", request);
    return data;
  },

  updateCartItem: async (
    productId: number,
    request: UpdateCartRequest
  ): Promise<Cart> => {
    const { data } = await apiClient.put<Cart>(
      `/api/v1/cart/items/${productId}`,
      request
    );
    return data;
  },

  removeFromCart: async (productId: number): Promise<Cart> => {
    const { data } = await apiClient.delete<Cart>(
      `/api/v1/cart/items/${productId}`
    );
    return data;
  },

  clearCart: async (): Promise<void> => {
    await apiClient.delete("/api/v1/cart");
  },
};

// React Query Hooks

export const useCart = (enabled: boolean = true) => {
  return useQuery<Cart, AxiosError>({
    queryKey: ["cart"],
    queryFn: () => cartApi.getCart(),
    enabled,
    staleTime: 1 * 60 * 1000, // 1 minute
    retry: false, // Don't retry if user not logged in
  });
};

export const useAddToCart = () => {
  const queryClient = useQueryClient();

  return useMutation<Cart, AxiosError, AddToCartRequest>({
    mutationFn: cartApi.addToCart,
    onSuccess: (data) => {
      // Update cart cache
      queryClient.setQueryData(["cart"], data);
      console.log("✅ Added to cart");
    },
    onError: (error) => {
      console.error("❌ Failed to add to cart:", error.response?.data);
    },
  });
};

export const useUpdateCartItem = () => {
  const queryClient = useQueryClient();

  return useMutation<Cart, AxiosError, { productId: number; quantity: number }>(
    {
      mutationFn: ({ productId, quantity }) =>
        cartApi.updateCartItem(productId, { quantity }),
      onSuccess: (data) => {
        queryClient.setQueryData(["cart"], data);
        console.log("✅ Cart updated");
      },
      onError: (error) => {
        console.error("❌ Failed to update cart:", error.response?.data);
      },
    }
  );
};

export const useRemoveFromCart = () => {
  const queryClient = useQueryClient();

  return useMutation<Cart, AxiosError, number>({
    mutationFn: cartApi.removeFromCart,
    onSuccess: (data) => {
      queryClient.setQueryData(["cart"], data);
      console.log("✅ Removed from cart");
    },
    onError: (error) => {
      console.error("❌ Failed to remove from cart:", error.response?.data);
    },
  });
};

export const useClearCart = () => {
  const queryClient = useQueryClient();

  return useMutation<void, AxiosError, void>({
    mutationFn: cartApi.clearCart,
    onSuccess: () => {
      queryClient.setQueryData(["cart"], null);
      console.log("✅ Cart cleared");
    },
    onError: (error) => {
      console.error("❌ Failed to clear cart:", error.response?.data);
    },
  });
};

// Export cartApi for use outside hooks if needed
export { cartApi };
