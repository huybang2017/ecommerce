"use client";

import { ProductItem } from "@/types/product-item";

interface SKUSelectorProps {
  productItems: ProductItem[];
  selectedSKU: ProductItem | null;
  onSKUChange: (sku: ProductItem) => void;
  isLoading?: boolean;
}

export default function SKUSelector({
  productItems,
  selectedSKU,
  onSKUChange,
  isLoading,
}: SKUSelectorProps) {
  if (isLoading) {
    return (
      <div className="space-y-4 animate-pulse">
        <div className="h-8 bg-neutral-200 rounded w-32"></div>
        <div className="flex gap-2">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-10 w-20 bg-neutral-200 rounded"></div>
          ))}
        </div>
      </div>
    );
  }

  if (!productItems || productItems.length === 0) {
    return (
      <div className="text-sm text-neutral-500">
        Sản phẩm này chưa có biến thể
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="text-neutral-500 text-sm">Chọn phiên bản:</div>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
        {productItems.map((item) => {
          const isSelected = selectedSKU?.id === item.id;
          const isAvailable = item.status === "ACTIVE" && item.qty_in_stock > 0;

          return (
            <button
              key={item.id}
              onClick={() => isAvailable && onSKUChange(item)}
              disabled={!isAvailable}
              className={`
                relative border rounded-sm p-3 text-left transition-all
                ${
                  isSelected
                    ? "border-[#ee4d2d] bg-[#fff5f3] text-[#ee4d2d] font-medium"
                    : "border-neutral-200 hover:border-[#ee4d2d]"
                }
                ${
                  !isAvailable
                    ? "opacity-50 cursor-not-allowed bg-neutral-50"
                    : "cursor-pointer"
                }
              `}
            >
              {/* SKU Code */}
              <div className="text-xs font-medium mb-1">{item.sku_code}</div>

              {/* Price */}
              <div className="text-sm font-semibold">
                ₫{new Intl.NumberFormat("vi-VN").format(item.price)}
              </div>

              {/* Stock Status */}
              <div className="text-xs text-neutral-500 mt-1">
                {isAvailable ? (
                  <span className="text-green-600">
                    Còn {item.qty_in_stock}
                  </span>
                ) : (
                  <span className="text-red-500">Hết hàng</span>
                )}
              </div>

              {/* Selected Indicator */}
              {isSelected && (
                <div className="absolute top-0 right-0 w-0 h-0 border-t-[20px] border-r-[20px] border-t-[#ee4d2d] border-r-transparent">
                  <svg
                    className="absolute -top-[18px] -right-[2px] w-3 h-3 text-white"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" />
                  </svg>
                </div>
              )}
            </button>
          );
        })}
      </div>

      {/* Info Text */}
      {selectedSKU && (
        <div className="text-sm text-neutral-600 bg-blue-50 border border-blue-200 rounded p-3">
          <span className="font-medium">Đã chọn:</span> {selectedSKU.sku_code} -
          ₫{new Intl.NumberFormat("vi-VN").format(selectedSKU.price)}
        </div>
      )}
    </div>
  );
}
