"use client";

import { useCartContext as useCart } from "@/contexts/CartContext";
import Link from "next/link";
import { useState } from "react";

export default function CartPage() {
  const {
    cart,
    loading,
    error,
    itemCount,
    total,
    updateItem,
    removeItem,
    clear,
  } = useCart();
  const [updating, setUpdating] = useState<number | null>(null);

  const handleQuantityChange = async (
    productId: number,
    newQuantity: number
  ) => {
    if (newQuantity < 1) {
      await removeItem(productId);
      return;
    }
    setUpdating(productId);
    try {
      await updateItem(productId, newQuantity);
    } finally {
      setUpdating(null);
    }
  };

  const handleRemove = async (productId: number) => {
    setUpdating(productId);
    try {
      await removeItem(productId);
    } finally {
      setUpdating(null);
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  if (loading && !cart) {
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
        <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
          <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-red-700">
            {error}
          </div>
        </main>
      </div>
    );
  }

  const items = cart?.items ? Object.values(cart.items) : [];

  return (
    <div className="min-h-screen bg-white">
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <h1 className="mb-8 text-4xl font-semibold tracking-tight text-neutral-900">
          Shopping Cart
        </h1>

        {items.length === 0 ? (
          <div className="rounded-xl border border-neutral-200 bg-neutral-50 p-12 text-center">
            <svg
              className="mx-auto mb-4 h-16 w-16 text-neutral-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
              />
            </svg>
            <p className="mb-4 text-lg font-medium text-neutral-600">
              Your cart is empty
            </p>
            <Link
              href="/products"
              className="inline-block rounded-lg bg-neutral-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-neutral-800"
            >
              Continue Shopping
            </Link>
          </div>
        ) : (
          <div className="grid gap-8 lg:grid-cols-3">
            {/* Cart Items */}
            <div className="lg:col-span-2">
              <div className="space-y-4">
                {items.map((item) => (
                  <div
                    key={item.product_id}
                    className="flex gap-6 rounded-xl border border-neutral-200 bg-white p-6"
                  >
                    {/* Product Image */}
                    <div className="h-24 w-24 flex-shrink-0 overflow-hidden rounded-lg border border-neutral-200 bg-neutral-100">
                      {item.image ? (
                        <img
                          src={item.image}
                          alt={item.name}
                          className="h-full w-full object-cover"
                        />
                      ) : (
                        <div className="flex h-full w-full items-center justify-center text-neutral-400">
                          <svg
                            className="h-8 w-8"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                            />
                          </svg>
                        </div>
                      )}
                    </div>

                    {/* Product Info */}
                    <div className="flex flex-1 flex-col">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <Link
                            href={`/products/${item.product_id}`}
                            className="text-lg font-medium text-neutral-900 hover:text-neutral-700"
                          >
                            {item.name}
                          </Link>
                          {item.sku && (
                            <p className="mt-1 text-sm text-neutral-500">
                              SKU: {item.sku}
                            </p>
                          )}
                        </div>
                        <button
                          onClick={() => handleRemove(item.product_id)}
                          disabled={updating === item.product_id}
                          className="ml-4 text-neutral-400 transition-colors hover:text-red-600 disabled:opacity-50"
                        >
                          <svg
                            className="h-5 w-5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M6 18L18 6M6 6l12 12"
                            />
                          </svg>
                        </button>
                      </div>

                      <div className="mt-4 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <label className="text-sm font-medium text-neutral-700">
                            Quantity:
                          </label>
                          <div className="flex items-center gap-2">
                            <button
                              onClick={() =>
                                handleQuantityChange(
                                  item.product_id,
                                  item.quantity - 1
                                )
                              }
                              disabled={
                                updating === item.product_id ||
                                item.quantity <= 1
                              }
                              className="flex h-8 w-8 items-center justify-center rounded border border-neutral-300 bg-white text-neutral-700 transition-colors hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                              <svg
                                className="h-4 w-4"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M20 12H4"
                                />
                              </svg>
                            </button>
                            <span className="w-12 text-center text-sm font-medium text-neutral-900">
                              {item.quantity}
                            </span>
                            <button
                              onClick={() =>
                                handleQuantityChange(
                                  item.product_id,
                                  item.quantity + 1
                                )
                              }
                              disabled={updating === item.product_id}
                              className="flex h-8 w-8 items-center justify-center rounded border border-neutral-300 bg-white text-neutral-700 transition-colors hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                              <svg
                                className="h-4 w-4"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M12 4v16m8-8H4"
                                />
                              </svg>
                            </button>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="text-lg font-semibold text-neutral-900">
                            {formatPrice(item.price * item.quantity)}
                          </p>
                          <p className="text-sm text-neutral-500">
                            {formatPrice(item.price)} each
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* Clear Cart Button */}
              <div className="mt-6">
                <button
                  onClick={clear}
                  disabled={loading}
                  className="rounded-lg border border-neutral-300 bg-white px-4 py-2 text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50 disabled:opacity-50"
                >
                  Clear Cart
                </button>
              </div>
            </div>

            {/* Order Summary */}
            <div className="lg:col-span-1">
              <div className="sticky top-8 rounded-xl border border-neutral-200 bg-neutral-50 p-6">
                <h2 className="mb-4 text-xl font-semibold text-neutral-900">
                  Order Summary
                </h2>
                <div className="space-y-3 border-b border-neutral-200 pb-4">
                  <div className="flex justify-between text-sm text-neutral-600">
                    <span>Items ({itemCount})</span>
                    <span>{formatPrice(total)}</span>
                  </div>
                </div>
                <div className="mt-4 flex justify-between text-lg font-semibold text-neutral-900">
                  <span>Total</span>
                  <span>{formatPrice(total)}</span>
                </div>
                <Link
                  href="/checkout"
                  className="mt-6 block w-full rounded-sm bg-[#ee4d2d] px-6 py-3 text-center text-sm font-medium text-white transition-colors hover:bg-[#d73211] disabled:opacity-50 disabled:cursor-not-allowed uppercase shadow-sm"
                >
                  Mua Hàng
                </Link>
                <Link
                  href="/products"
                  className="mt-3 block w-full rounded-sm border border-neutral-300 bg-white px-6 py-3 text-center text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50"
                >
                  Tiếp Tục Mua Sắm
                </Link>
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
