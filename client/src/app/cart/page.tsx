"use client";

import Link from "next/link";
import { useState, useEffect, useMemo } from "react";
import { useCartContext as useCart } from "@/contexts/CartContext";
import CartItemRow from "@/components/CartItemRow";

function formatPrice(value: number) {
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency: "VND",
  }).format(value);
}

export default function CartPage() {
  const { cart, loading, itemCount, total, updateItem, removeItem, clear } =
    useCart();
  const items = useMemo(() => cart?.items || [], [cart?.items]);
  const [updating, setUpdating] = useState<number | null>(null);

  // Selection state for items (to support selecting items like Shopee)
  const [selected, setSelected] = useState<number[]>([]);

  useEffect(() => {
    // default: select all items when cart changes
    setSelected(items.map((i) => i.product_item_id));
  }, [items]);

  const isAllSelected = items.length > 0 && selected.length === items.length;
  const selectedTotal = items
    .filter((i) => selected.includes(i.product_item_id))
    .reduce((s, i) => s + i.price * i.quantity, 0);

  const toggleSelect = (id: number) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
    );
  };

  const toggleSelectAll = () => {
    if (isAllSelected) setSelected([]);
    else setSelected(items.map((i) => i.product_item_id));
  };

  const handleRemove = async (productItemId: number) => {
    setUpdating(productItemId);
    try {
      await removeItem(productItemId);
    } finally {
      setUpdating(null);
    }
  };

  const handleQuantityChange = async (
    productItemId: number,
    quantity: number
  ) => {
    if (quantity < 1) return;
    setUpdating(productItemId);
    try {
      await updateItem(productItemId, quantity);
    } finally {
      setUpdating(null);
    }
  };

  const handleClear = async () => {
    await clear();
  };

  return (
    <div className="min-h-screen bg-white">
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8 pb-24">
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
          <>
            <div className="grid gap-8 lg:grid-cols-3">
              {/* Cart Items */}
              <div className="lg:col-span-2">
                {/* Header: checkbox and column labels */}
                <div className="mb-4 flex items-center justify-between rounded-t-lg border border-neutral-200 bg-white p-3">
                  <div className="flex items-center gap-4">
                    <input
                      type="checkbox"
                      checked={isAllSelected}
                      onChange={toggleSelectAll}
                      className="h-4 w-4"
                    />
                    <span className="text-sm font-medium">Sản Phẩm</span>
                  </div>

                  <div className="hidden sm:flex gap-6 text-sm text-neutral-600">
                    <span className="w-28 text-right">Đơn Giá</span>
                    <span className="w-24 text-center">Số Lượng</span>
                    <span className="w-28 text-right">Số Tiền</span>
                    <span className="w-24 text-right">Thao Tác</span>
                  </div>
                </div>

                <div className="space-y-3">
                  {items.map((item) => (
                    <CartItemRow
                      key={item.product_item_id}
                      item={item}
                      updating={updating}
                      checked={selected.includes(item.product_item_id)}
                      onToggle={toggleSelect}
                      onRemove={handleRemove}
                      onQuantityChange={handleQuantityChange}
                    />
                  ))}
                </div>

                {/* Vouchers / Shipping placeholders */}
                <div className="mt-4 rounded-lg border border-neutral-200 bg-white p-4 text-sm text-neutral-600">
                  <div>
                    Voucher giảm giá ·{" "}
                    <a className="text-[#1677ff]">Xem thêm</a>
                  </div>
                  <div className="mt-2">
                    Giảm phí vận chuyển đơn tối thiểu 0₫ ·{" "}
                    <a className="text-[#1677ff]">Tìm hiểu thêm</a>
                  </div>
                </div>

                {/* Actions */}
                <div className="mt-6 flex items-center gap-4">
                  <input
                    type="checkbox"
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                    className="h-4 w-4"
                  />
                  <span className="text-sm text-neutral-600">
                    Chọn Tất Cả ({items.length})
                  </span>
                  <button
                    onClick={handleClear}
                    disabled={loading}
                    className="ml-4 rounded-lg border border-neutral-300 bg-white px-4 py-2 text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50 disabled:opacity-50"
                  >
                    Xóa
                  </button>
                  <button className="text-sm text-[#ee4d2d]">
                    Lưu vào mục Đã thích
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
                    <span>Tổng</span>
                    <span>{formatPrice(total)}</span>
                  </div>

                  <Link
                    href="/checkout"
                    className="mt-6 block w-full rounded-sm bg-[#ee4d2d] px-6 py-3 text-center text-sm font-medium uppercase text-white transition-colors hover:bg-[#d73211] disabled:cursor-not-allowed disabled:opacity-50 shadow-sm"
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
            {/* Bottom sticky bar similar to Shopee */}
            <div className="fixed bottom-0 left-0 right-0 z-50 border-t bg-white">
              <div className="mx-auto max-w-7xl px-4 py-3 sm:px-6 lg:px-8">
                <div className="flex items-center gap-4">
                  <input
                    type="checkbox"
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                    className="h-4 w-4"
                  />
                  <span className="text-sm text-neutral-700">
                    Chọn Tất Cả ({items.length})
                  </span>
                  <button className="text-sm text-neutral-500">Xóa</button>
                  <div className="ml-auto flex items-center gap-4">
                    <div className="text-sm">
                      Tổng cộng ({selected.length} sản phẩm):{" "}
                      <span className="ml-2 text-lg font-semibold text-neutral-900">
                        {formatPrice(selectedTotal)}
                      </span>
                    </div>
                    <Link
                      href="/checkout"
                      className="rounded-sm bg-[#ee4d2d] px-5 py-2 text-sm font-medium text-white"
                    >
                      Mua Hàng
                    </Link>
                  </div>
                </div>
              </div>
            </div>{" "}
          </>
        )}
      </main>
    </div>
  );
}
