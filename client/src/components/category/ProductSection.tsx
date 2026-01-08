"use client";

import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { SortBar } from "./SortBar";
import { ProductGrid } from "./ProductGrid";
import { getProducts } from "@/lib/mock-category-api";

interface ProductSectionProps {
  initialProducts: any;
  categoryId: number;
}

/**
 * Client component that handles product filtering without page reload.
 * Similar to Shopee - only products update when clicking child categories.
 */
export function ProductSection({
  initialProducts,
  categoryId,
}: ProductSectionProps) {
  const searchParams = useSearchParams();
  const [products, setProducts] = useState(initialProducts);
  const [isLoading, setIsLoading] = useState(false);

  // Get active category ID from search params (when child category clicked)
  const activeCategoryId = searchParams.get("category_id");
  const sortBy = searchParams.get("sort_by");
  const order = searchParams.get("order");

  useEffect(() => {
    // Fetch products when category or filters change
    const fetchProducts = async () => {
      setIsLoading(true);
      try {
        const targetCategoryId = activeCategoryId
          ? parseInt(activeCategoryId)
          : categoryId;

        const filters = {
          sort_by: sortBy || undefined,
          order: order || undefined,
        };

        const newProducts = await getProducts(targetCategoryId, filters);
        setProducts(newProducts);
      } catch (error) {
        console.error("Failed to fetch products:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchProducts();
  }, [activeCategoryId, sortBy, order, categoryId]);

  return (
    <div className="flex-1">
      {/* Sort Bar */}
      <SortBar
        sortBy={sortBy || ""}
        order={order || ""}
        pageNum={products.page}
        totalPages={products.last_page}
      />

      {/* Product Grid with loading state */}
      {isLoading ? (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-2 mt-5">
          {[...Array(10)].map((_, i) => (
            <div key={i} className="bg-white animate-pulse">
              <div className="aspect-square bg-neutral-200"></div>
              <div className="p-2 space-y-2">
                <div className="h-4 bg-neutral-200 rounded"></div>
                <div className="h-4 bg-neutral-200 rounded w-2/3"></div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <ProductGrid products={products.data} />
      )}
    </div>
  );
}
