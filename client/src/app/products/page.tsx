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
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            All Products
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Browse our complete collection
          </p>
        </div>

        <div className="mb-6">
          <div className="flex flex-wrap gap-2">
            <Link
              href="/products"
              className={`rounded-lg px-4 py-2 text-sm font-medium transition-colors ${
                !categoryId
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
              }`}
            >
              All
            </Link>
            {categories.map((category) => (
              <Link
                key={category.id}
                href={`/products?category_id=${category.id}`}
                className={`rounded-lg px-4 py-2 text-sm font-medium transition-colors ${
                  categoryId === category.id
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-700 hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
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
            <div className="mb-6 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {products.length === 0 ? (
                <div className="col-span-full text-center text-gray-500 dark:text-gray-400">
                  No products found
                </div>
              ) : (
                products.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))
              )}
            </div>

            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2">
                {page > 1 && (
                  <Link
                    href={`/products?page=${page - 1}${categoryId ? `&category_id=${categoryId}` : ''}`}
                    className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
                  >
                    Previous
                  </Link>
                )}
                <span className="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
                  Page {page} of {totalPages}
                </span>
                {page < totalPages && (
                  <Link
                    href={`/products?page=${page + 1}${categoryId ? `&category_id=${categoryId}` : ''}`}
                    className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
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

