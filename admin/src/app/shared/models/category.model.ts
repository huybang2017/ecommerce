export interface Category {
  id: number;
  parent_id?: number | null;
  name: string;
  slug: string;
  description: string;
  image_url?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export type CategorySortBy = 'id' | 'name' | 'slug' | 'is_active' | 'created_at' | 'updated_at';

export interface CategoryListParams {
  page?: number;
  page_size?: number;
  search?: string;
  status?: 'all' | 'active' | 'inactive';
  sort_by?: CategorySortBy;
  sort_order?: 'asc' | 'desc';
}

export interface PaginatedResponse<T> {
  data: T[];
  limit: number;
  page: number;
  total: number;
  total_pages: number;
}
