import { getProducts, getCategories } from '@/lib/api';
import ProductCard from '@/components/ProductCard';
import { ProductCardSkeleton } from '@/components/Loading';
import Error from '@/components/Error';
import { Product, Category } from '@/lib/types';
import Link from 'next/link';

interface ProductsPageProps {
  searchParams: Promise<{
    page?: string;
    limit?: string;
    category_id?: string;
  }>;
}

export default async function ProductsPage({ searchParams }: ProductsPageProps) {
  const params = await searchParams;
  const page = parseInt(params.page || '1', 10);
  const limit = parseInt(params.limit || '20', 10);
  const categoryId = params.category_id
    ? parseInt(params.category_id, 10)
    : undefined;

  let products: Product[] = [];
  let total = 0;
  let categories: Category[] = [];
  let error: string | null = null;

  try {
    const [productsResponse, categoriesResponse] = await Promise.all([
      getProducts({
        page,
        limit,
        category_id: categoryId,
        status: 'ACTIVE',
      }),
      getCategories(),
    ]);

    products = productsResponse.products || [];
    total = productsResponse.total || 0;
    categories = categoriesResponse || [];
  } catch (err: any) {
    error = err?.message || 'Failed to load products';
  }

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="min-h-screen bg-white">
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="mb-10">
          <h1 className="text-4xl font-semibold tracking-tight text-neutral-900 sm:text-5xl">
            All Products
          </h1>
          <p className="mt-4 text-lg text-neutral-600">
            Browse our complete collection
          </p>
        </div>

        <div className="mb-10">
          <div className="flex flex-wrap gap-3">
            <Link
              href="/products"
              className={`rounded-full px-5 py-2.5 text-sm font-medium transition-all ${
                !categoryId
                  ? 'bg-neutral-900 text-white shadow-sm'
                  : 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200'
              }`}
            >
              All
            </Link>
            {categories.map((category) => (
              <Link
                key={category.id}
                href={`/products?category_id=${category.id}`}
                className={`rounded-full px-5 py-2.5 text-sm font-medium transition-all ${
                  categoryId === category.id
                    ? 'bg-neutral-900 text-white shadow-sm'
                    : 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200'
                }`}
              >
                {category.name}
              </Link>
            ))}
          </div>
        </div>

        {error ? (
          <Error message={error} />
        ) : (
          <>
            <div className="mb-10 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {products.length === 0 ? (
                <div className="col-span-full py-16 text-center">
                  <p className="text-neutral-500">No products found</p>
                </div>
              ) : (
                products.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))
              )}
            </div>

            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-3">
                {page > 1 && (
                  <Link
                    href={`/products?page=${page - 1}${categoryId ? `&category_id=${categoryId}` : ''}`}
                    className="rounded-lg border border-neutral-200 bg-white px-5 py-2.5 text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50 hover:border-neutral-300"
                  >
                    Previous
                  </Link>
                )}
                <span className="px-5 py-2.5 text-sm font-medium text-neutral-600">
                  Page {page} of {totalPages}
                </span>
                {page < totalPages && (
                  <Link
                    href={`/products?page=${page + 1}${categoryId ? `&category_id=${categoryId}` : ''}`}
                    className="rounded-lg border border-neutral-200 bg-white px-5 py-2.5 text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50 hover:border-neutral-300"
                  >
                    Next
                  </Link>
                )}
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}

