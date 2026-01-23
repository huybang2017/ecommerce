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
  const shopId = (item as any).shop_id;

  return (
    <div
      className={`flex items-start gap-4 rounded-md border border-[#eee] bg-white p-4 shadow-sm ${
        checked ? "ring-2 ring-[#ee4d2d]/20" : ""
      }`}
    >
      <div className="flex items-start pt-1">
        <input
          type="checkbox"
          checked={checked}
          onChange={() => onToggle(item.product_item_id)}
          className="h-4 w-4 rounded border-neutral-300 text-[#ee4d2d]"
        />
      </div>

      <div className="flex flex-1 gap-4">
        <div className="h-20 w-20 flex-shrink-0 overflow-hidden rounded-sm border border-neutral-200 bg-neutral-100">
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

        <div className="flex flex-1 flex-col">
          <div className="flex items-start justify-between">
            <div className="flex-1 pr-4">
              <div className="mb-1 flex items-center gap-2">
                <span className="rounded bg-[#fff2ef] px-2 py-0.5 text-xs font-medium text-[#ee4d2d]">
                  Yêu thích
                </span>
                {shopId ? (
                  <span className="rounded bg-neutral-100 px-2 py-0.5 text-xs text-neutral-600">{`Nhựa_Việt_Nhật`}</span>
                ) : null}
              </div>
              <Link
                href={`/products/${item.product_item_id}`}
                className="text-sm font-medium text-neutral-900 hover:text-neutral-700"
              >
                {displayName}
              </Link>
              {(item as any).variant && (
                <p className="mt-1 text-xs text-neutral-500">
                  Phân Loại Hàng: {(item as any).variant}
                </p>
              )}
            </div>

            <div className="text-right flex-shrink-0">
              <p className="text-sm font-semibold text-[#ee4d2d]">
                {new Intl.NumberFormat("vi-VN", {
                  style: "currency",
                  currency: "VND",
                }).format(item.price * item.quantity)}
              </p>
              <p className="mt-1 text-xs text-neutral-500">
                {new Intl.NumberFormat("vi-VN", {
                  style: "currency",
                  currency: "VND",
                }).format(item.price)}{" "}
                / cái
              </p>
            </div>
          </div>

          <div className="mt-3 flex items-center justify-between">
            <div className="flex items-center gap-4 text-sm text-neutral-600">
              <div className="flex items-center gap-2">
                <button
                  onClick={() =>
                    onQuantityChange(
                      item.product_item_id,
                      Math.max(1, item.quantity - 1),
                    )
                  }
                  disabled={
                    updating === item.product_item_id || item.quantity <= 1
                  }
                  className="flex h-8 w-8 items-center justify-center rounded border border-neutral-300 bg-white text-neutral-700 hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  -
                </button>
                <div className="w-11">
                  <span className="block w-full text-center text-sm font-medium text-neutral-900">
                    {item.quantity}
                  </span>
                </div>
                <button
                  onClick={() =>
                    onQuantityChange(item.product_item_id, item.quantity + 1)
                  }
                  disabled={updating === item.product_item_id}
                  className="flex h-8 w-8 items-center justify-center rounded border border-neutral-300 bg-white text-neutral-700 hover:bg-neutral-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  +
                </button>
              </div>
            </div>

            <div className="flex items-center gap-4 text-sm">
              <button
                onClick={() => onRemove(item.product_item_id)}
                disabled={updating === item.product_item_id}
                className="text-[#ee4d2d] hover:underline disabled:opacity-50"
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
