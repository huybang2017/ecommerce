import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/axios-client";
import { AxiosError } from "axios";

// Types
export interface Product {
  id: number;
  name: string;
  description?: string;
  price: number;
  original_price?: number;
  category_id?: number;
  category?: Category;
  image_url?: string;
  images?: string[];
  stock_quantity: number;
  sku?: string;
  weight?: number;
  dimensions?: string;
  is_active: boolean;
  rating?: number;
  review_count?: number;
  sold_count?: number;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  description?: string;
  slug: string;
  parent_id?: number;
  image_url?: string;
  is_active: boolean;
  display_order?: number;
  product_count?: number;
  created_at: string;
  updated_at: string;
}

export interface ProductsResponse {
  products: Product[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface SearchParams {
  q?: string;
  category?: string;
  min_price?: number;
  max_price?: number;
  sort_by?: "price" | "rating" | "created_at" | "sold_count";
  order?: "asc" | "desc";
  page?: number;
  limit?: number;
}

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

  getCategoryProducts: async (
    id: number,
    params?: SearchParams
  ): Promise<ProductsResponse> => {
    const { data } = await apiClient.get<ProductsResponse>(
      `/api/v1/categories/${id}/products`,
      { params }
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

export const useCategoryProducts = (id: number, params?: SearchParams) => {
  return useQuery<ProductsResponse, AxiosError>({
    queryKey: ["category", id, "products", params],
    queryFn: () => productApi.getCategoryProducts(id, params),
    staleTime: 2 * 60 * 1000,
  });
};

// Export productApi for use outside hooks if needed
export { productApi };
