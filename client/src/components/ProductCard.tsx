import Link from 'next/link';
import Image from 'next/image';
import { Product } from '@/lib/types';

interface ProductCardProps {
  product: Product;
}

export default function ProductCard({ product }: ProductCardProps) {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
    }).format(price);
  };

  const imageUrl =
    product.images && Array.isArray(product.images) && product.images.length > 0
      ? product.images[0]
      : '/placeholder-product.jpg';

  return (
    <Link
      href={`/products/${product.id}`}
      className="group relative flex flex-col overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm transition-all hover:shadow-lg dark:border-gray-800 dark:bg-gray-900"
    >
      <div className="relative aspect-square w-full overflow-hidden bg-gray-100 dark:bg-gray-800">
        {imageUrl.startsWith('http') ? (
          <img
            src={imageUrl}
            alt={product.name}
            className="h-full w-full object-cover transition-transform group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-gray-200 text-gray-400 dark:bg-gray-700 dark:text-gray-500">
            <span className="text-sm">No Image</span>
          </div>
        )}
        {product.status === 'INACTIVE' && (
          <div className="absolute top-2 right-2 rounded bg-red-500 px-2 py-1 text-xs font-semibold text-white">
            Out of Stock
          </div>
        )}
      </div>
      <div className="flex flex-1 flex-col p-4">
        <h3 className="mb-2 line-clamp-2 text-lg font-semibold text-gray-900 dark:text-gray-100">
          {product.name}
        </h3>
        {product.description && (
          <p className="mb-3 line-clamp-2 text-sm text-gray-600 dark:text-gray-400">
            {product.description}
          </p>
        )}
        <div className="mt-auto flex items-center justify-between">
          <span className="text-xl font-bold text-blue-600 dark:text-blue-400">
            {formatPrice(product.price)}
          </span>
          {product.stock > 0 && (
            <span className="text-xs text-gray-500 dark:text-gray-400">
              {product.stock} in stock
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}

