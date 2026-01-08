"use client";

import { useCategoryChildren } from "@/hooks/useProducts";
import { FilterSidebar } from "./FilterSidebar";

interface CategoryFilterWrapperProps {
  categoryId: number;
  categoryName: string;
  isParent?: boolean;
  currentSlug: string;
  currentCategoryId: number;
}

/**
 * Client Component wrapper for FilterSidebar that fetches children categories
 * using React Query. This follows Next.js 13+ best practices by separating
 * server-side initial data from client-side interactive data fetching.
 */
export function CategoryFilterWrapper({
  categoryId,
  categoryName,
  isParent = true,
  currentSlug,
  currentCategoryId,
}: CategoryFilterWrapperProps) {
  const {
    data: children,
    isLoading,
    error,
  } = useCategoryChildren(categoryId, isParent);

  // Extract parent slug and ID from currentSlug for child category generation
  // Format: thoi-trang-nam.6 or thoi-trang-nam.6.23
  const slugParts = currentSlug.split(".");
  const baseSlug = slugParts[0]; // thoi-trang-nam
  const parentId = slugParts.length >= 2 ? slugParts[1] : categoryId.toString();

  // Loading state - show skeleton or placeholder
  if (isLoading) {
    return (
      <aside className="w-[230px] flex-shrink-0 hidden md:block">
        <div className="animate-pulse">
          <div className="h-6 bg-neutral-200 rounded mb-4"></div>
          <div className="space-y-2">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-4 bg-neutral-200 rounded"></div>
            ))}
          </div>
        </div>
      </aside>
    );
  }

  // Error state - fallback to empty categories
  if (error) {
    console.error("Failed to fetch category children:", error);
    return (
      <FilterSidebar
        categories={[]}
        currentCategory={{ id: categoryId, name: categoryName }}
        isParent={isParent}
        currentSlug={currentSlug}
        currentCategoryId={currentCategoryId}
      />
    );
  }

  // Transform API data with correct slug format
  // Parent children: baseSlug.parentId.childId
  const categories = (children || []).map((cat) => ({
    id: cat.id,
    name: cat.name,
    slug: `${baseSlug}.${parentId}.${cat.id}`, // Format: thoi-trang-nam.6.23
  }));

  return (
    <FilterSidebar
      categories={categories}
      currentCategory={{ id: categoryId, name: categoryName }}
      isParent={isParent}
      currentSlug={currentSlug}
      currentCategoryId={currentCategoryId}
    />
  );
}
