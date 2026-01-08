import type { Category } from "./product";

export interface CategoryDetailResponse extends Category {
  parent?: Category;
  children?: Category[];
}

export interface BreadcrumbItem {
  name: string;
  path: string;
  slug: string;
}

export interface CategoryDetail extends Category {
  path: BreadcrumbItem[];
  children: Category[];
}

export interface ProductFilters {
  page?: number;
  limit?: number;
  sort_by?: string;
  order?: "asc" | "desc";
  min_price?: number;
  max_price?: number;
}
