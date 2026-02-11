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

export interface CategoryListParams {
  page?: number;
  page_size?: number;
  search?: string;
  is_active?: boolean | null;
  sort_by?: 'name' | 'created_at' | 'updated_at';
  sort_order?: 'asc' | 'desc';
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}
