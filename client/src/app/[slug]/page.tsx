import { notFound } from "next/navigation";
import { getCategoryDetail, getProducts } from "@/lib/mock-category-api";
import { Breadcrumb } from "@/components/category/Breadcrumb";
import { FilterSidebar } from "@/components/category/FilterSidebar";
import { SortBar } from "@/components/category/SortBar";
import { ProductGrid } from "@/components/category/ProductGrid";

interface PageProps {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

export default async function CategoryPage({
  params,
  searchParams,
}: PageProps) {
  const { slug } = await params;
  const sp = await searchParams; // search params object

  const category = await getCategoryDetail(slug);

  if (!category) {
    notFound();
  }

  // Fetch products with filters
  const products = await getProducts(category.id, sp);

  return (
    <div className="bg-[#f5f5f5] min-h-screen pb-10">
      <div className="max-w-[1200px] mx-auto px-4 pt-5">
        {/* Breadcrumb */}
        <Breadcrumb items={category.path} />

        <div className="flex gap-5 mt-5">
          {/* Sidebar */}
          <FilterSidebar
            categories={category.children || []}
            currentCategory={{ id: category.id, name: category.name }}
            isParent={true}
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
