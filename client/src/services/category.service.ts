import { apiClient } from "@/lib/axios-client";
import type { ProductsResponse } from "./product.service";
import type { Category } from "@/types/product";
import type {
  CategoryDetailResponse,
  CategoryDetail,
  BreadcrumbItem,
  ProductFilters,
} from "@/types/category";

// ==================== CATEGORIES ====================

/**
 * Get list of all categories
 */
export async function getCategories(): Promise<Category[]> {
  const { data } = await apiClient.get<Category[]>("/api/v1/categories");
  return data;
}

/**
 * Get single category by ID (basic info)
 */
export async function getCategory(id: number): Promise<Category> {
  const { data } = await apiClient.get<Category>(`/api/v1/categories/${id}`);
  return data;
}

/**
 * Build breadcrumb path from category data
 * Note: Backend no longer returns parent/children objects, only parent_id
 */
async function buildBreadcrumbPath(
  category: CategoryDetailResponse
): Promise<BreadcrumbItem[]> {
  const breadcrumbPath: BreadcrumbItem[] = [];

  // If has parent_id, fetch parent to build full path
  if (category.parent_id) {
    try {
      const { data: parent } = await apiClient.get<CategoryDetailResponse>(
        `/api/v1/categories/${category.parent_id}`
      );

      // Add parent first
      breadcrumbPath.push({
        name: parent.name,
        path: `/${parent.slug}.${parent.id}`,
        slug: `${parent.slug}.${parent.id}`,
      });
    } catch (err) {
      console.error("Failed to fetch parent category", err);
    }
  }

  // Add current category
  breadcrumbPath.push({
    name: category.name,
    path: category.parent_id
      ? `/${category.slug}.${category.parent_id}.${category.id}`
      : `/${category.slug}.${category.id}`,
    slug: `${category.slug}.${category.id}`,
  });

  return breadcrumbPath;
}

/**
 * Fetch category detail by ID
 */
export async function getCategoryById(id: number): Promise<CategoryDetail> {
  const { data } = await apiClient.get<CategoryDetailResponse>(
    `/api/v1/categories/${id}`
  );

  return {
    id: data.id,
    name: data.name,
    slug: data.slug,
    description: data.description,
    parent_id: data.parent_id,
    path: await buildBreadcrumbPath(data),
    children: [],
    is_active: data.is_active,
    created_at: data.created_at,
    updated_at: data.updated_at,
  };
}

/**
 * Fetch products by category ID with filters
 */
export async function getCategoryProducts(
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
