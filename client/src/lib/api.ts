// API client for Product Service via API Gateway

import { Product, Category, ProductsResponse, ApiError } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
      });

      if (!response.ok) {
        const error: ApiError = await response.json().catch(() => ({
          error: 'Unknown error',
          message: `HTTP ${response.status}: ${response.statusText}`,
        }));
        throw new Error(error.message || error.error || 'Request failed');
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Network error or invalid response');
    }
  }

  // Product APIs
  async getProducts(params?: {
    page?: number;
    limit?: number;
    category_id?: number;
    status?: string;
  }): Promise<ProductsResponse> {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('limit', params.limit.toString());
    if (params?.category_id) searchParams.append('category_id', params.category_id.toString());
    if (params?.status) searchParams.append('status', params.status);

    const queryString = searchParams.toString();
    const endpoint = `/api/v1/products${queryString ? `?${queryString}` : ''}`;
    
    return this.request<ProductsResponse>(endpoint);
  }

  async getProduct(id: number): Promise<Product> {
    return this.request<Product>(`/api/v1/products/${id}`);
  }

  async searchProducts(query: string, category?: string): Promise<Product[]> {
    const searchParams = new URLSearchParams();
    searchParams.append('q', query);
    if (category) searchParams.append('category', category);

    const response = await this.request<{ products: Product[]; count: number }>(
      `/api/v1/products/search?${searchParams.toString()}`
    );
    return response.products;
  }

  // Category APIs
  async getCategories(): Promise<Category[]> {
    return this.request<Category[]>('/api/v1/categories');
  }

  async getCategory(id: number): Promise<Category> {
    return this.request<Category>(`/api/v1/categories/${id}`);
  }

  async getProductsByCategory(
    categoryId: number,
    page?: number,
    limit?: number
  ): Promise<ProductsResponse> {
    const searchParams = new URLSearchParams();
    if (page) searchParams.append('page', page.toString());
    if (limit) searchParams.append('limit', limit.toString());

    const queryString = searchParams.toString();
    const endpoint = `/api/v1/categories/${categoryId}/products${queryString ? `?${queryString}` : ''}`;
    
    return this.request<ProductsResponse>(endpoint);
  }
}

// Export singleton instance
export const apiClient = new ApiClient(API_BASE_URL);

// Export convenience functions
export const getProducts = (params?: {
  page?: number;
  limit?: number;
  category_id?: number;
  status?: string;
}) => apiClient.getProducts(params);

export const getProduct = (id: number) => apiClient.getProduct(id);

export const searchProducts = (query: string, category?: string) =>
  apiClient.searchProducts(query, category);

export const getCategories = () => apiClient.getCategories();

export const getCategory = (id: number) => apiClient.getCategory(id);

export const getProductsByCategory = (
  categoryId: number,
  page?: number,
  limit?: number
) => apiClient.getProductsByCategory(categoryId, page, limit);

