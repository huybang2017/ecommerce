import { getProducts } from '@/lib/api';
import ProductCard from '@/components/ProductCard';
import { ProductCardSkeleton } from '@/components/Loading';
import Error from '@/components/Error';
import { Product } from '@/lib/types';

export default async function Home() {
  let products: Product[] = [];
  let error: string | null = null;

  try {
    const response = await getProducts({ page: 1, limit: 12, status: 'ACTIVE' });
    products = response.products || [];
  } catch (err: any) {
    error = err?.message || 'Failed to load products';
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            Featured Products
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Discover our latest collection
          </p>
        </div>

        {error ? (
          <Error message={error} />
        ) : (
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {products.length === 0 ? (
              <div className="col-span-full text-center text-gray-500 dark:text-gray-400">
                No products available
              </div>
            ) : (
              products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))
            )}
          </div>
        )}
      </main>
    </div>
  );
}
