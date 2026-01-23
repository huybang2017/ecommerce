"use client";

import Link from "next/link";
import { useState, useEffect, useMemo } from "react";
import { useCartContext as useCart } from "@/contexts/CartContext";
import {
  useSetItemSelection,
  useSetSelection,
  useClearSelectedItems,
  useValidateCart,
} from "@/hooks/useCart";
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
  const [searchKeyword, setSearchKeyword] = useState("");

  const visibleItems = useMemo(() => {
    const q = searchKeyword.trim().toLowerCase();
    if (!q) return items;
    return items.filter((i) => {
      const name = (i.name || i.product_name || "").toString().toLowerCase();
      const sku = (i.sku || "").toString().toLowerCase();
      return name.includes(q) || sku.includes(q);
    });
  }, [items, searchKeyword]);
  const [updating, setUpdating] = useState<number | null>(null);

  // Selection state for items (to support selecting items like Shopee)
  const [selected, setSelected] = useState<number[]>([]);

  useEffect(() => {
    // Initialize selected state from backend `is_selected` field when cart changes
    const initiallySelected = items
      .filter((i) => (i as any).is_selected)
      .map((i) => i.product_item_id);
    if (initiallySelected.length > 0) setSelected(initiallySelected);
  }, [items]);

  const isAllSelected =
    visibleItems.length > 0 &&
    visibleItems.every((i) => selected.includes(i.product_item_id));
  const selectedTotal = items
    .filter((i) => selected.includes(i.product_item_id))
    .reduce((s, i) => s + i.price * i.quantity, 0);

  const toggleSelect = (id: number) => {
    const next = selected.includes(id)
      ? selected.filter((x) => x !== id)
      : [...selected, id];

    // optimistically update UI
    setSelected(next);

    // persist selection to backend
    setItemSelection.mutate({ productItemId: id, selected: next.includes(id) });
  };

  const toggleSelectAll = () => {
    const selectAll = !isAllSelected;
    const visibleIds = visibleItems.map((i) => i.product_item_id);
    if (selectAll) setSelected(visibleIds);
    else setSelected([]);

    // persist select all state for backend (best-effort)
    setSelection.mutate({ selected: selectAll });
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
    quantity: number,
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
    // clear selected items on backend
    if (selected.length > 0) {
      clearSelected.mutate();
      setSelected([]);
    } else {
      // fallback: clear entire cart
      await clear();
    }
  };

  const handleValidate = async () => {
    validateCart.mutate();
  };

  const setItemSelection = useSetItemSelection();
  const setSelection = useSetSelection();
  const clearSelected = useClearSelectedItems();
  const validateCart = useValidateCart();

  return (
    <div className="min-h-screen bg-white">
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8 pb-24">
        {/* <h1 className="mb-8 text-4xl font-semibold tracking-tight text-neutral-900">
          Shopping Cart
        </h1> */}

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
            <div className="grid gap-8">
              {/* Cart Items */}
              <div>
                {/* Header: checkbox, column labels and cart-local search */}
                <div className="mb-4 flex flex-col gap-3 rounded-t-lg border border-neutral-200 bg-white p-3">
                  <div className="flex items-center justify-between">
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

                  <div className="mx-auto w-full max-w-[720px]">
                    <div className="relative">
                      <input
                        type="text"
                        value={searchKeyword}
                        onChange={(e) => setSearchKeyword(e.target.value)}
                        placeholder="TÌM KIẾM TRONG GIỎ HÀNG"
                        className="w-full rounded-sm border border-[#ee4d2d] py-2 pl-4 pr-36 text-sm outline-none"
                      />
                      <button
                        onClick={(e) => e.preventDefault()}
                        className="absolute right-1 top-1 bottom-1 flex h-[36px] w-14 items-center justify-center rounded-sm bg-[#ee4d2d] text-white"
                      >
                        Tìm
                      </button>
                    </div>
                  </div>
                </div>

                <div className="space-y-3">
                  {visibleItems.map((item) => (
                    <CartItemRow
                      key={item.product_item_id}
                      item={item}
                      updating={updating}
                      checked={
                        selected.includes(item.product_item_id) ||
                        (item as any).is_selected
                      }
                      onToggle={toggleSelect}
                      onRemove={handleRemove}
                      onQuantityChange={handleQuantityChange}
                    />
                  ))}
                  {visibleItems.length === 0 && (
                    <div className="p-6 text-center text-neutral-500">
                      Không có sản phẩm phù hợp
                    </div>
                  )}
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
