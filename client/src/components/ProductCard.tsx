'use client';

import Link from 'next/link';
import Image from 'next/image';
import { Product } from '@/lib/types';
import { useCartContext as useCart } from '@/contexts/CartContext';
import { useState } from 'react';

interface ProductCardProps {
  product: Product;
}

export default function ProductCard({ product }: ProductCardProps) {
  const { addItem } = useCart();
  const [adding, setAdding] = useState(false);

  const handleAddToCart = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (product.status !== 'ACTIVE' || product.stock === 0) return;

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
        <div className="mt-auto space-y-3">
          <div className="flex items-center justify-between">
            <span className="text-xl font-semibold text-neutral-900">
              {formatPrice(product.price)}
            </span>
            {product.stock > 0 && (
              <span className="text-xs font-medium text-neutral-400">
                {product.stock} in stock
              </span>
            )}
          </div>
          <button
            onClick={handleAddToCart}
            disabled={product.status !== 'ACTIVE' || product.stock === 0 || adding}
            className="w-full rounded-lg bg-neutral-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-neutral-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {adding ? 'Adding...' : product.status === 'ACTIVE' && product.stock > 0 ? 'Add to Cart' : 'Out of Stock'}
          </button>
        </div>
      </div>
    </Link>
  );
}

