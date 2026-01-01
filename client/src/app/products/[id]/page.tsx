'use client';

import { getProduct } from '@/lib/api';
import Error from '@/components/Error';
import { notFound, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useCartContext as useCart } from '@/contexts/CartContext';
import { useState, useEffect } from 'react';
import { Product } from '@/lib/types';

interface ProductDetailPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default function ProductDetailPage({
  params,
}: ProductDetailPageProps) {
  const router = useRouter();
  const { addItem } = useCart();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [productId, setProductId] = useState<number | null>(null);

  useEffect(() => {
    async function loadProduct() {
      const resolvedParams = await params;
      const id = parseInt(resolvedParams.id, 10);
      
      if (isNaN(id)) {
        setError('Invalid product ID');
        setLoading(false);
        return;
      }

      setProductId(id);

      try {
        const productData = await getProduct(id);
        setProduct(productData);
      } catch (err: any) {
        if (err?.message?.includes('404')) {
          notFound();
        }
        setError(err?.message || 'Failed to load product');
      } finally {
        setLoading(false);
      }
    }

    loadProduct();
  }, [params]);

  const handleAddToCart = async () => {
    if (!product || product.status !== 'ACTIVE' || product.stock === 0) return;

    setAdding(true);
    try {
      await addItem({
        product_id: product.id,
        name: product.name,
        price: product.price,
        quantity: 1,
        image: product.images && Array.isArray(product.images) && product.images.length > 0 
          ? product.images[0] 
          : undefined,
        sku: product.sku,
      });
    } catch (err) {
      console.error('Failed to add to cart:', err);
    } finally {
      setAdding(false);
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
    }).format(price);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
          <div className="flex items-center justify-center py-24">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-neutral-200 border-t-neutral-900"></div>
          </div>
        </main>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-white">
        <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
          <Error message={error} />
        </main>
      </div>
    );
  }

  if (!product) {
    notFound();
  }

  const images =
    product.images && Array.isArray(product.images) && product.images.length > 0
      ? product.images
      : [];

  return (
    <div className="min-h-screen bg-white">
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <Link
          href="/products"
          className="mb-8 inline-flex items-center text-sm font-medium text-neutral-600 hover:text-neutral-900 transition-colors"
        >
          <svg
            className="mr-2 h-4 w-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
          Back to Products
        </Link>

        <div className="grid gap-12 lg:grid-cols-2">
          {/* Product Images */}
          <div className="space-y-4">
            {images.length > 0 ? (
              <div className="relative aspect-square w-full overflow-hidden rounded-xl border border-neutral-200 bg-neutral-50">
                <img
                  src={images[0]}
                  alt={product.name}
                  className="h-full w-full object-cover"
                />
              </div>
            ) : (
              <div className="flex aspect-square w-full items-center justify-center rounded-xl border border-neutral-200 bg-neutral-50 text-neutral-400">
                <span className="font-medium">No Image Available</span>
              </div>
            )}
            {images.length > 1 && (
              <div className="grid grid-cols-4 gap-4">
                {images.slice(1, 5).map((image, index) => (
                  <div
                    key={index}
                    className="relative aspect-square overflow-hidden rounded-lg border border-neutral-200 bg-neutral-50"
                  >
                    <img
                      src={image}
                      alt={`${product.name} ${index + 2}`}
                      className="h-full w-full object-cover"
                    />
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Product Info */}
          <div className="space-y-8">
            <div>
              <h1 className="text-4xl font-semibold tracking-tight text-neutral-900">
                {product.name}
              </h1>
              {product.category && (
                <p className="mt-3 text-base text-neutral-600">
                  Category: {product.category.name}
                </p>
              )}
            </div>

            <div>
              <p className="text-4xl font-semibold text-neutral-900">
                {formatPrice(product.price)}
              </p>
              <p className="mt-2 text-sm font-medium text-neutral-500">
                SKU: {product.sku}
              </p>
            </div>

            {product.description && (
              <div>
                <h2 className="mb-3 text-lg font-semibold text-neutral-900">
                  Description
                </h2>
                <p className="text-neutral-600 leading-relaxed">
                  {product.description}
                </p>
              </div>
            )}

            <div className="space-y-3 rounded-xl border border-neutral-200 bg-neutral-50 p-6">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-neutral-700">
                  Status:
                </span>
                <span
                  className={`rounded-full px-3 py-1 text-xs font-semibold ${
                    product.status === 'ACTIVE'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-red-100 text-red-800'
                  }`}
                >
                  {product.status}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-neutral-700">
                  Stock:
                </span>
                <span className="text-sm font-medium text-neutral-600">
                  {product.stock} units
                </span>
              </div>
            </div>

            <div className="pt-4">
              <button
                onClick={handleAddToCart}
                disabled={product.status !== 'ACTIVE' || product.stock === 0 || adding}
                className="w-full rounded-lg bg-neutral-900 px-6 py-4 text-base font-medium text-white transition-colors hover:bg-neutral-800 disabled:cursor-not-allowed disabled:opacity-40"
              >
                {adding 
                  ? 'Adding...' 
                  : product.status === 'ACTIVE' && product.stock > 0
                  ? 'Add to Cart'
                  : 'Out of Stock'}
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
