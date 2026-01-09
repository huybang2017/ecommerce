export interface Product {
  id: number;
  name: string;
  description?: string;
  base_price: number;
  category_id?: number;
  category?: Category;
  image_url?: string;
  images?: string[];
  is_active: boolean;
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
  category_id?: number;
  min_price?: number;
  max_price?: number;
  sort_by?: "price" | "created_at" | "sold_count" | "rating";
  order?: "asc" | "desc";
  page?: number;
  limit?: number;
}
