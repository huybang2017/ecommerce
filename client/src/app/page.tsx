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
    <div className="min-h-screen bg-white">
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="mb-12 text-center">
          <h1 className="text-4xl font-semibold tracking-tight text-neutral-900 sm:text-5xl">
            Featured Products
          </h1>
          <p className="mt-4 text-lg text-neutral-600">
            Discover our latest collection
          </p>
        </div>

        {error ? (
          <Error message={error} />
        ) : (
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {products.length === 0 ? (
              <div className="col-span-full py-16 text-center">
                <p className="text-neutral-500">No products available</p>
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
