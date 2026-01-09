import { useQuery } from "@tanstack/react-query";
import { apiClient } from "./axios-client";
import type { Category, ProductsResponse } from "@/hooks/useProducts";

interface CategoryDetailResponse extends Category {
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

interface ProductFilters {
  page?: number;
  limit?: number;
  sort_by?: string;
  order?: "asc" | "desc";
  min_price?: number;
  max_price?: number;
}

/**
 * Build breadcrumb path from category data
 */
function buildBreadcrumbPath(
  category: CategoryDetailResponse
): BreadcrumbItem[] {
  const breadcrumbPath: BreadcrumbItem[] = [];

  // Add parent to breadcrumb if exists (for child categories)
  if (category.parent) {
    breadcrumbPath.push({
      name: category.parent.name,
      path: `/${category.parent.slug}.${category.parent.id}`,
      slug: `${category.parent.slug}.${category.parent.id}`,
    });

    // Add current child category
    breadcrumbPath.push({
      name: category.name,
      path: `/${category.parent.slug}.${category.parent.id}.${category.id}`,
      slug: `${category.parent.slug}.${category.parent.id}.${category.id}`,
    });
  } else {
    // Parent category only - just add itself
    breadcrumbPath.push({
      name: category.name,
      path: `/${category.slug}.${category.id}`,
      slug: `${category.slug}.${category.id}`,
    });
  }

  return breadcrumbPath;
}

/**
 * Fetch category detail by ID
 */
async function fetchCategoryDetail(id: number): Promise<CategoryDetail> {
  const { data } = await apiClient.get<CategoryDetailResponse>(
    `/api/v1/categories/${id}`
  );

  return {
    id: data.id,
    name: data.name,
    slug: data.slug,
    description: data.description,
    parent_id: data.parent_id,
    path: buildBreadcrumbPath(data),
    children: data.children || [],
    is_active: data.is_active,
    created_at: data.created_at,
    updated_at: data.updated_at,
  };
}

/**
 * Fetch products by category ID with filters
 */
async function fetchCategoryProducts(
  categoryId: number,
  filters: ProductFilters = {}
) {
  const params = {
    page: filters.page || 1,
    limit: filters.limit || 20,
    sort_by: filters.sort_by || "created_at",
    order: filters.order || "desc",
    ...(filters.min_price && { min_price: filters.min_price }),
    ...(filters.max_price && { max_price: filters.max_price }),
  };

  const { data } = await apiClient.get<ProductsResponse>(
    `/api/v1/categories/${categoryId}/products`,
    { params }
  );

  return {
    data: data.products || [],
    total: data.total || 0,
    page: data.page || 1,
    last_page: data.total_pages || 1,
  };
}

/**
 * React Query hook for fetching category detail
 */
export function useCategoryDetail(categoryId: number) {
  return useQuery({
    queryKey: ["category", categoryId],
    queryFn: () => fetchCategoryDetail(categoryId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: !!categoryId,
  });
}

/**
 * React Query hook for fetching category products
 */
export function useCategoryProducts(
  categoryId: number,
  filters: ProductFilters = {}
) {
  return useQuery({
    queryKey: ["category-products", categoryId, filters],
    queryFn: () => fetchCategoryProducts(categoryId, filters),
    staleTime: 1 * 60 * 1000, // 1 minute
    enabled: !!categoryId,
  });
}
