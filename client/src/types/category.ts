import { Category } from "@/types/product";

// API Response (backend no longer returns parent/children objects)
export interface CategoryDetailResponse extends Category {
  // ❌ Removed: parent, children (backend refactored domain)
}

export interface BreadcrumbItem {
  name: string;
  path: string;
  slug: string;
}

export interface CategoryDetail extends Category {
  path: BreadcrumbItem[];
  children: Category[]; // Always empty array (for backward compat)
}

export interface ProductFilters {
  page?: number;
  limit?: number;
  sort_by?: string;
  order?: "asc" | "desc";
  min_price?: number;
  max_price?: number;
}
