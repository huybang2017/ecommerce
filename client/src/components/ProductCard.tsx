"use client";

import Link from "next/link";
import Image from "next/image";
import { Product } from "@/lib/types";
import { useCartContext as useCart } from "@/contexts/CartContext";
import { useState } from "react";

interface ProductCardProps {
  product: Product;
}

export default function ProductCard({ product }: ProductCardProps) {
  const { addItem } = useCart();
  const [adding, setAdding] = useState(false);

  const handleAddToCart = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    setAdding(true);
    try {
      await addItem(product.id, 1);
    } catch (err) {
      console.error("Failed to add to cart:", err);
    } finally {
      setAdding(false);
    }
  };
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  const imageUrl =
    product.images && Array.isArray(product.images) && product.images.length > 0
      ? product.images[0]
      : "/placeholder-product.jpg";

  return (
    <Link
      href={`/products/${product.id}`}
      className="group relative flex flex-col overflow-hidden rounded-sm border border-transparent bg-white shadow-sm transition-all hover:-translate-y-0.5 hover:border-orange-500 hover:shadow-md hover:z-10"
    >
      <div className="relative aspect-square w-full overflow-hidden bg-neutral-50">
        {imageUrl.startsWith("http") ? (
          <img
            src={imageUrl}
            alt={product.name}
            className="h-full w-full object-cover transition-opacity duration-300 group-hover:opacity-90"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-neutral-100 text-neutral-400">
            <span className="text-xs">No Image</span>
          </div>
        )}

        {/* Discount Badge */}
        <div className="absolute right-0 top-0 flex flex-col items-center bg-yellow-400 px-1 py-0.5 text-[10px] font-semibold leading-tight text-red-600 after:absolute after:bottom-[-4px] after:left-0 after:border-l-[18px] after:border-r-[18px] after:border-t-[4px] after:border-l-transparent after:border-r-transparent after:border-t-yellow-400 after:content-['']">
          <span>50%</span>
          <span className="text-white uppercase">GIẢM</span>
        </div>

        {/* Generic Mall/Favorite Badge if needed */}
        <div className="absolute left-[-4px] top-2 rounded-r-sm bg-[#ee4d2d] px-1.5 py-0.5 text-[10px] font-medium text-white shadow-sm before:absolute before:bottom-[-3px] before:left-0 before:border-l-[3px] before:border-t-[3px] before:border-l-transparent before:border-t-[#9e311b] before:content-['']">
          Yêu thích
        </div>
      </div>

      <div className="flex flex-1 flex-col p-2">
        <h3 className="mb-1 line-clamp-2 text-xs text-neutral-800 leading-normal min-h-[2.5em]">
          {product.name}
        </h3>

        <div className="mt-auto pt-2">
          {/* Price Section */}
          <div className="flex items-baseline gap-1 flex-wrap">
            <span className="text-[10px] text-neutral-400 line-through">
              {formatPrice(product.base_price * 1.5)}
            </span>
            <span className="text-sm font-medium text-[#ee4d2d]">
              <span className="text-[10px] align-top">₫</span>
              {new Intl.NumberFormat("vi-VN").format(product.base_price)}
            </span>
          </div>

          {/* Footer: Rating & Sold */}
          <div className="mt-2 flex items-center justify-between text-[10px] text-neutral-500">
            <div className="flex items-center gap-0.5 text-yellow-400">
              <span>★★★★★</span>
            </div>
            <span>Đã bán {product.sold_count || 0}</span>
          </div>
        </div>
      </div>
    </Link>
  );
}
