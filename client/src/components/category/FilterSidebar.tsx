"use client";

import { useState } from "react";
import { Filter, Star, ChevronDown, ChevronUp } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";

interface FilterSidebarProps {
  categories: { id: number; name: string; slug: string }[];
  currentCategory: { id: number; name: string };
  isParent: boolean;
}

export function FilterSidebar({
  categories,
  currentCategory,
  isParent,
}: FilterSidebarProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [priceRange, setPriceRange] = useState({ min: "", max: "" });

  const handleApplyPrice = () => {
    const params = new URLSearchParams(searchParams.toString());
    if (priceRange.min) params.set("min_price", priceRange.min);
    else params.delete("min_price");

    if (priceRange.max) params.set("max_price", priceRange.max);
    else params.delete("max_price");

    router.push(`?${params.toString()}`, { scroll: false });
  };

  const handleClearAll = () => {
    router.push(window.location.pathname, { scroll: false });
    setPriceRange({ min: "", max: "" });
  };

  return (
    <aside className="w-[230px] flex-shrink-0 hidden md:block">
      {/* Category Navigation */}
      <div className="mb-8">
        <h3 className="flex items-center gap-2 font-bold text-neutral-800 mb-4 pb-2 border-b border-neutral-200">
          <Filter size={16} /> Bộ lọc tìm kiếm
        </h3>

        <div className="mb-4">
          <div className="font-semibold mb-2 text-sm text-neutral-700">
            Theo Danh Mục
          </div>
          <ul className="text-sm space-y-2 pl-2">
            <li className={`font-medium ${!isParent ? "text-[#ee4d2d]" : ""}`}>
              {isParent ? (
                <span className="font-bold flex items-center gap-1">
                  <ChevronDown size={14} /> {currentCategory.name}
                </span>
              ) : (
                <span className="flex items-center gap-1">
                  <ChevronDown size={14} /> {currentCategory.name}
                </span>
              )}
            </li>
            {categories.map((cat) => (
              <li
                key={cat.id}
                className="cursor-pointer hover:text-[#ee4d2d] transition-colors pl-4"
              >
                {cat.name}
              </li>
            ))}
          </ul>
        </div>

        {/* Dynamic Filters - Mocked */}
        <div className="space-y-6">
          {/* Location - Example Attribute */}
          <div>
            <h4 className="text-sm font-medium mb-2">Nơi bán</h4>
            <div className="space-y-2 text-sm text-neutral-600">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  className="rounded-sm border-neutral-300"
                />
                <span>Hà Nội</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  className="rounded-sm border-neutral-300"
                />
                <span>TP. Hồ Chí Minh</span>
              </label>
            </div>
          </div>

          {/* Price Range */}
          <div>
            <h4 className="text-sm font-medium mb-2">Khoảng Giá</h4>
            <div className="flex items-center gap-1 mb-2">
              <input
                type="number"
                placeholder="₫ TỪ"
                className="w-full border border-neutral-300 rounded-sm px-2 py-1 text-xs focus:border-neutral-500 outline-none"
                value={priceRange.min}
                onChange={(e) =>
                  setPriceRange({ ...priceRange, min: e.target.value })
                }
              />
              <span className="text-neutral-400">-</span>
              <input
                type="number"
                placeholder="₫ ĐẾN"
                className="w-full border border-neutral-300 rounded-sm px-2 py-1 text-xs focus:border-neutral-500 outline-none"
                value={priceRange.max}
                onChange={(e) =>
                  setPriceRange({ ...priceRange, max: e.target.value })
                }
              />
            </div>
            <button className="w-full bg-[#ee4d2d] text-white text-sm py-1 rounded-sm uppercase font-medium hover:bg-[#d73211] transition-colors">
              Áp dụng
            </button>
          </div>

          {/* Rating */}
          <div>
            <h4 className="text-sm font-medium mb-2">Đánh Giá</h4>
            <div className="space-y-1">
              {[5, 4, 3, 2, 1].map((star) => (
                <div
                  key={star}
                  className="flex items-center text-sm gap-2 cursor-pointer hover:bg-neutral-50 py-1 px-2 -ml-2 rounded-sm"
                >
                  <div className="flex text-yellow-400 gap-0.5">
                    {Array.from({ length: 5 }).map((_, i) => (
                      <Star
                        key={i}
                        size={14}
                        fill={i < star ? "currentColor" : "none"}
                        className={i >= star ? "text-neutral-300" : ""}
                      />
                    ))}
                  </div>
                  {star < 5 && (
                    <span className="text-xs text-neutral-500">trở lên</span>
                  )}
                </div>
              ))}
            </div>
          </div>

          {/* Clear Filter */}
          <button className="w-full bg-[#ee4d2d] text-white text-sm py-2 rounded-sm font-medium hover:bg-[#d73211]">
            Xóa tất cả
          </button>
        </div>
      </div>
    </aside>
  );
}
