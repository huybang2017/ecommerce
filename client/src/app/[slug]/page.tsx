import { notFound } from "next/navigation";
import { Breadcrumb } from "@/components/category/Breadcrumb";
import { CategoryFilterWrapper } from "@/components/category/CategoryFilterWrapper";
import { SortBar } from "@/components/category/SortBar";
import { ProductGrid } from "@/components/category/ProductGrid";
import {
  getCategoryById,
  getCategoryProducts,
} from "@/services/category.service";

interface PageProps {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

export default async function CategoryPage({
  params,
  searchParams,
}: PageProps) {
  const { slug } = await params;
  const sp = await searchParams;

  // Parse format: slug.parentId or slug.parentId.childId
  // Example: thoi-trang-nam.6 (parent) or thoi-trang-nam.6.23 (child)
  const slugParts = slug.split(".");
  const categoryId =
    slugParts.length === 2
      ? parseInt(slugParts[1]) // Parent: thoi-trang-nam.6 -> ID=6
      : parseInt(slugParts[2]); // Child: thoi-trang-nam.6.23 -> ID=23

  if (!categoryId || isNaN(categoryId)) {
    notFound();
  }

  // Fetch category detail by ID from API
  const category = await getCategoryById(categoryId);

  if (!category) {
    notFound();
  }

  // Determine if parent or child based on slug format
  const isParentCategory = slugParts.length === 2; // Only has slug.ID

  // For sidebar: parent shows its children, child shows siblings
  const categoryIdForChildren = isParentCategory
    ? category.id
    : slugParts.length === 3
    ? parseInt(slugParts[1])
    : category.parent_id!;

  // Fetch products for this specific category
  const products = await getCategoryProducts(category.id, sp);

  return (
    <div className="bg-[#f5f5f5] min-h-screen pb-10">
      <div className="max-w-[1200px] mx-auto px-4 pt-5">
        {/* Breadcrumb */}
        <Breadcrumb items={category.path} />

        <div className="flex gap-5 mt-5">
          {/* Sidebar - Always fetch parent's children for navigation */}
          <CategoryFilterWrapper
            categoryId={categoryIdForChildren}
            categoryName={
              isParentCategory
                ? category.name
                : category.parent?.name || category.name
            }
            isParent={isParentCategory}
            currentSlug={slug}
            currentCategoryId={category.id}
          />

          {/* Main Content */}
          <div className="flex-1">
            {/* Sort Bar */}
            <SortBar
              sortBy={sp.sort_by as string}
              order={sp.order as string}
              pageNum={products.page}
              totalPages={products.last_page}
            />

            {/* Product Grid */}
            <ProductGrid products={products.data} />
          </div>
        </div>
      </div>
    </div>
  );
}
