import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/axios-client";
import { AxiosError } from "axios";
import type {
  Product,
  Category,
  ProductsResponse,
  SearchParams,
} from "@/types/product";

// Re-export types for backward compatibility
export type {
  Product,
  Category,
  ProductsResponse,
  SearchParams,
} from "@/types/product";

// API functions
const productApi = {
  getProducts: async (params?: SearchParams): Promise<ProductsResponse> => {
    const { data } = await apiClient.get<ProductsResponse>("/api/v1/products", {
      params,
    });
    return data;
  },

  getProduct: async (id: number): Promise<Product> => {
    const { data } = await apiClient.get<Product>(`/api/v1/products/${id}`);
    return data;
  },

  searchProducts: async (params: SearchParams): Promise<ProductsResponse> => {
    const { data } = await apiClient.get<ProductsResponse>(
      "/api/v1/products/search",
      { params }
    );
    return data;
  },

  getCategories: async (): Promise<Category[]> => {
    const { data } = await apiClient.get<Category[]>("/api/v1/categories");
    return data;
  },

  getCategory: async (id: number): Promise<Category> => {
    const { data } = await apiClient.get<Category>(`/api/v1/categories/${id}`);
    return data;
  },

  getCategoryChildren: async (parentId: number): Promise<Category[]> => {
    const { data } = await apiClient.get<Category[]>(
      `/api/v1/categories/${parentId}/children`
    );
    return data;
  },
};

// React Query Hooks

export const useProducts = (params?: SearchParams) => {
  return useQuery<ProductsResponse, AxiosError>({
    queryKey: ["products", params],
    queryFn: () => productApi.getProducts(params),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const useProduct = (id: number, enabled: boolean = true) => {
  return useQuery<Product, AxiosError>({
    queryKey: ["product", id],
    queryFn: () => productApi.getProduct(id),
    enabled,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useSearchProducts = (params: SearchParams) => {
  return useQuery<ProductsResponse, AxiosError>({
    queryKey: ["products", "search", params],
    queryFn: () => productApi.searchProducts(params),
    enabled: !!params.q, // Only run if search query exists
    staleTime: 1 * 60 * 1000, // 1 minute
  });
};

export const useCategories = () => {
  return useQuery<Category[], AxiosError>({
    queryKey: ["categories"],
    queryFn: () => productApi.getCategories(),
    staleTime: 10 * 60 * 1000, // 10 minutes (categories don't change often)
  });
};

export const useCategory = (id: number, enabled: boolean = true) => {
  return useQuery<Category, AxiosError>({
    queryKey: ["category", id],
    queryFn: () => productApi.getCategory(id),
    enabled,
    staleTime: 10 * 60 * 1000,
  });
};

export const useCategoryChildren = (
  parentId: number,
  enabled: boolean = true
) => {
  return useQuery<Category[], AxiosError>({
    queryKey: ["category", parentId, "children"],
    queryFn: () => productApi.getCategoryChildren(parentId),
    enabled,
    staleTime: 10 * 60 * 1000, // 10 minutes (categories don't change often)
  });
};

// Export productApi for use outside hooks if needed
export { productApi };
