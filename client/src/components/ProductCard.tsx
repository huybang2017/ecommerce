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
      className="group relative flex flex-col overflow-hidden rounded-xl border border-neutral-200 bg-white transition-all hover:border-neutral-300 hover:shadow-md"
    >
      <div className="relative aspect-square w-full overflow-hidden bg-neutral-50">
        {imageUrl.startsWith('http') ? (
          <img
            src={imageUrl}
            alt={product.name}
            className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-neutral-100 text-neutral-400">
            <span className="text-sm font-medium">No Image</span>
          </div>
        )}
        {product.status === 'INACTIVE' && (
          <div className="absolute top-3 right-3 rounded-full bg-neutral-900 px-3 py-1 text-xs font-medium text-white">
            Out of Stock
          </div>
        )}
      </div>
      <div className="flex flex-1 flex-col p-5">
        <h3 className="mb-2 line-clamp-2 text-base font-semibold text-neutral-900 leading-tight">
          {product.name}
        </h3>
        {product.description && (
          <p className="mb-4 line-clamp-2 text-sm text-neutral-500 leading-relaxed">
            {product.description}
          </p>
        )}
        <div className="mt-auto flex items-center justify-between">
          <span className="text-xl font-semibold text-neutral-900">
            {formatPrice(product.price)}
          </span>
          {product.stock > 0 && (
            <span className="text-xs font-medium text-neutral-400">
              {product.stock} in stock
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}

