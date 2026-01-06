"use client";

import { ChevronLeft, ChevronRight } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";

interface SortBarProps {
  sortBy?: string;
  order?: string;
  pageNum: number;
  totalPages: number;
}

export function SortBar({
  sortBy = "popular",
  order = "desc",
  pageNum,
  totalPages,
}: SortBarProps) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const onSortChange = (newSort: string, newOrder: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("sort_by", newSort);
    params.set("order", newOrder);
    router.push(`?${params.toString()}`, { scroll: false });
  };

  return (
    <div className="bg-[#ededed] py-3 px-5 rounded-sm flex justify-between items-center mb-4">
      <div className="flex items-center gap-3 text-sm">
        <span className="text-neutral-600">Sắp xếp theo</span>
        <button
          onClick={() => onSortChange("popular", "desc")}
          className={`px-4 py-2 rounded-sm capitalize ${
            sortBy === "popular"
              ? "bg-[#ee4d2d] text-white"
              : "bg-white hover:bg-neutral-50"
          }`}
        >
          Phổ biến
        </button>
        <button
          onClick={() => onSortChange("newest", "desc")}
          className={`px-4 py-2 rounded-sm capitalize ${
            sortBy === "newest"
              ? "bg-[#ee4d2d] text-white"
              : "bg-white hover:bg-neutral-50"
          }`}
        >
          Mới nhất
        </button>
        <button
          onClick={() => onSortChange("sales", "desc")}
          className={`px-4 py-2 rounded-sm capitalize ${
            sortBy === "sales"
              ? "bg-[#ee4d2d] text-white"
              : "bg-white hover:bg-neutral-50"
          }`}
        >
          Bán chạy
        </button>
        <div className="relative group">
          <button
            className={`w-48 px-4 py-2 bg-white rounded-sm flex justify-between items-center ${
              sortBy === "price" ? "text-[#ee4d2d]" : ""
            }`}
          >
            <span>
              {sortBy === "price"
                ? order === "asc"
                  ? "Giá: Thấp đến Cao"
                  : "Giá: Cao đến Thấp"
                : "Giá"}
            </span>
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M19 9l-7 7-7-7"
              ></path>
            </svg>
          </button>
          <div className="absolute top-full left-0 w-full bg-white shadow-md rounded-sm z-10 hidden group-hover:block py-1">
            <button
              onClick={() => onSortChange("price", "asc")}
              className="block w-full text-left px-4 py-2 hover:bg-neutral-50 hover:text-[#ee4d2d]"
            >
              Giá: Thấp đến Cao
            </button>
            <button
              onClick={() => onSortChange("price", "desc")}
              className="block w-full text-left px-4 py-2 hover:bg-neutral-50 hover:text-[#ee4d2d]"
            >
              Giá: Cao đến Thấp
            </button>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {/* Pagination logic would go here */}
        <div className="text-sm">
          <span className="text-[#ee4d2d]">{pageNum}</span>
          <span className="text-neutral-500">/{totalPages}</span>
        </div>
      </div>
    </div>
  );
}
