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
  shop_id?: number;
  is_selected?: boolean;
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

export interface ValidateCartSuccess {
  message: string;
}

export interface ValidateCartError {
  error: string;
}

export interface SetSelectionRequest {
  selected: boolean;
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
    productItemId: number,
    request: UpdateCartRequest,
  ): Promise<Cart> => {
    console.log("🔄 updateCartItem called:", {
      productItemId,
      request,
      url: `/api/v1/cart/items/${productItemId}`,
    });

    const { data } = await apiClient.put<Cart>(
      `/api/v1/cart/items/${productItemId}`,
      request,
    );

    console.log("✅ updateCartItem response:", data);
    return data;
  },

  removeFromCart: async (productItemId: number): Promise<Cart> => {
    const { data } = await apiClient.delete<Cart>(
      `/api/v1/cart/items/${productItemId}`,
    );
    return data;
  },

  clearCart: async (): Promise<void> => {
    await apiClient.delete("/api/v1/cart");
  },

  setItemSelection: async (
    productItemId: number,
    request: SetSelectionRequest,
  ): Promise<Cart> => {
    const { data } = await apiClient.patch<Cart>(
      `/api/v1/cart/items/${productItemId}/selection`,
      request,
    );
    return data;
  },

  setSelection: async (request: SetSelectionRequest): Promise<Cart> => {
    const { data } = await apiClient.patch<Cart>(
      `/api/v1/cart/selection`,
      request,
    );
    return data;
  },

  setShopSelection: async (
    shopId: number,
    request: SetSelectionRequest,
  ): Promise<Cart> => {
    const { data } = await apiClient.patch<Cart>(
      `/api/v1/cart/shops/${shopId}/selection`,
      request,
    );
    return data;
  },

  clearSelected: async (): Promise<Cart> => {
    const { data } = await apiClient.delete<Cart>(`/api/v1/cart/selected`);
    return data;
  },

  validateCart: async (): Promise<ValidateCartSuccess> => {
    const { data } = await apiClient.post<ValidateCartSuccess>(
      `/api/v1/cart/validate`,
    );
    return data;
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
      // Backend returns message only; refetch cart to refresh state
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Added to cart");
    },
    onError: (error) => {
      console.error("❌ Failed to add to cart:", error.response?.data);
    },
  });
};

export const useUpdateCartItem = () => {
  const queryClient = useQueryClient();

  return useMutation<
    Cart,
    AxiosError,
    { productItemId: number; quantity: number }
  >({
    mutationFn: ({ productItemId, quantity }) =>
      cartApi.updateCartItem(productItemId, { quantity }),
    onSuccess: (data) => {
      // Backend returns message only; refetch cart to get latest state
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Cart updated");
    },
    onError: (error) => {
      console.error("❌ Failed to update cart:", error.response?.data);
    },
  });
};

export const useRemoveFromCart = () => {
  const queryClient = useQueryClient();

  return useMutation<Cart, AxiosError, number>({
    mutationFn: cartApi.removeFromCart,
    onSuccess: (data) => {
      // Backend returns message only; refetch cart to refresh state
      queryClient.invalidateQueries({ queryKey: ["cart"] });
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

export const useSetItemSelection = () => {
  const queryClient = useQueryClient();

  return useMutation<
    Cart,
    AxiosError,
    { productItemId: number; selected: boolean }
  >({
    mutationFn: ({ productItemId, selected }) =>
      cartApi.setItemSelection(productItemId, { selected }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Item selection updated");
    },
    onError: (error) => {
      console.error(
        "❌ Failed to update item selection:",
        error.response?.data,
      );
    },
  });
};

export const useSetSelection = () => {
  const queryClient = useQueryClient();

  return useMutation<Cart, AxiosError, { selected: boolean }>({
    mutationFn: ({ selected }) => cartApi.setSelection({ selected }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Selection updated for all items");
    },
    onError: (error) => {
      console.error("❌ Failed to set selection:", error.response?.data);
    },
  });
};

export const useSetShopSelection = () => {
  const queryClient = useQueryClient();

  return useMutation<Cart, AxiosError, { shopId: number; selected: boolean }>({
    mutationFn: ({ shopId, selected }) =>
      cartApi.setShopSelection(shopId, { selected }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Shop selection updated");
    },
    onError: (error) => {
      console.error("❌ Failed to set shop selection:", error.response?.data);
    },
  });
};

export const useClearSelectedItems = () => {
  const queryClient = useQueryClient();

  return useMutation<Cart, AxiosError, void>({
    mutationFn: cartApi.clearSelected,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Cleared selected items");
    },
    onError: (error) => {
      console.error("❌ Failed to clear selected items:", error.response?.data);
    },
  });
};

export const useValidateCart = () => {
  const queryClient = useQueryClient();

  return useMutation<ValidateCartSuccess, AxiosError, void>({
    mutationFn: () => cartApi.validateCart(),
    onSuccess: (data) => {
      // validation may not change cart, but refetch to be safe
      queryClient.invalidateQueries({ queryKey: ["cart"] });
      console.log("✅ Cart validated", data.message);
    },
    onError: (error) => {
      console.error("❌ Failed to validate cart:", error.response?.data);
    },
  });
};

// Export cartApi for use outside hooks if needed
export { cartApi };
