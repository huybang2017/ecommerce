// TypeScript types for Product and Category

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  sku: string;
  category_id?: number;
  category?: Category;
  status: string;
  images: string[] | null;
  stock: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number;
  parent?: Category;
  children?: Category[];
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface ProductsResponse {
  products: Product[];
  total: number;
  page: number;
  limit: number;
}

export interface ApiError {
  error: string;
  message?: string;
}

// Search API types
export interface SearchParams {
  q?: string;
  category_id?: number;
  min_price?: number;
  max_price?: number;
  status?: string;
  sort_field?: 'price' | 'name' | 'created_at';
  sort_order?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

export interface SearchResponse {
  products: Product[];
  total: number;
  page: number;
  limit: number;
}

