import React from "react";
import Link from "next/link";
import type { CartItem } from "@/hooks/useCart";

interface Props {
  item: CartItem;
  updating: number | null;
  checked: boolean;
  onToggle: (id: number) => void;
  onRemove: (id: number) => void;
  onQuantityChange: (id: number, qty: number) => void;
}

export default function CartItemRow({
  item,
  updating,
  checked,
  onToggle,
  onRemove,
  onQuantityChange,
}: Props) {
  const displayImage = item.image || item.product_image;
  const displayName = item.name || item.product_name;

  return (
    <div className="flex items-start gap-4 rounded-lg border border-neutral-200 bg-white p-4">
      <div className="flex items-start">
        <input
          type="checkbox"
          checked={checked}
          onChange={() => onToggle(item.product_item_id)}
          className="h-4 w-4 rounded border-neutral-300 text-[#ee4d2d]"
        />
      </div>

      <div className="flex flex-1 gap-4">
        <div className="h-20 w-20 flex-shrink-0 overflow-hidden rounded-lg border border-neutral-200 bg-neutral-100">
          {displayImage ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={displayImage}
              alt={displayName}
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

        <div className="flex flex-1 flex-col gap-1">
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <Link
                href={`/products/${item.product_item_id}`}
                className="text-sm font-medium text-neutral-900 hover:text-neutral-700"
              >
                {displayName}
              </Link>
              {item.sku && (
                <p className="mt-1 text-xs text-neutral-500">SKU: {item.sku}</p>
              )}
            </div>
            <div className="text-right">
              <p className="text-sm font-semibold text-neutral-900">
                {new Intl.NumberFormat("vi-VN", {
                  style: "currency",
                  currency: "VND",
                }).format(item.price * item.quantity)}
              </p>
              <p className="text-xs text-neutral-500">
                {new Intl.NumberFormat("vi-VN", {
                  style: "currency",
                  currency: "VND",
                }).format(item.price)}{" "}
                each
              </p>
            </div>
          </div>

          <div className="mt-2 flex items-center justify-between">
            <div className="flex items-center gap-3 text-sm text-neutral-600">
              <span className="text-sm">Số lượng</span>
              <div className="flex items-center gap-2">
                <button
                  onClick={() =>
                    onQuantityChange(item.product_item_id, item.quantity - 1)
                  }
                  disabled={
                    updating === item.product_item_id || item.quantity <= 1
                  }
                  className="flex h-7 w-7 items-center justify-center rounded border border-neutral-300 bg-white text-neutral-700 transition-colors hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  -
                </button>
                <span className="w-9 text-center text-sm font-medium text-neutral-900">
                  {item.quantity}
                </span>
                <button
                  onClick={() =>
                    onQuantityChange(item.product_item_id, item.quantity + 1)
                  }
                  disabled={updating === item.product_item_id}
                  className="flex h-7 w-7 items-center justify-center rounded border border-neutral-300 bg-white text-neutral-700 transition-colors hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  +
                </button>
              </div>
            </div>

            <div className="flex items-center gap-4 text-sm">
              <button
                onClick={() => onRemove(item.product_item_id)}
                disabled={updating === item.product_item_id}
                className="text-neutral-500 hover:text-red-600 disabled:opacity-50"
              >
                Xóa
              </button>
              <button className="text-neutral-500 hover:text-neutral-700">
                Tìm sản phẩm tương tự
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
