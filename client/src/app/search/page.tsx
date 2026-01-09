"use client";

import { useState, useEffect, useCallback, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { searchProductsAdvanced } from "@/services/search.service";
import { getCategories } from "@/services/category.service";
import ProductCard from "@/components/ProductCard";
import { ProductCardSkeleton } from "@/components/Loading";
import Error from "@/components/Error";
import { Product, Category } from "@/lib/types";
import type { ProductsResponse } from "@/types/product";

function SearchContent() {
  const searchParams = useSearchParams();
  const router = useRouter();

  // State
  const [query, setQuery] = useState(searchParams.get("q") || "");
  const [categoryId, setCategoryId] = useState<number | undefined>(
    searchParams.get("category_id")
      ? parseInt(searchParams.get("category_id")!, 10)
      : undefined
  );
  const [minPrice, setMinPrice] = useState<string>(
    searchParams.get("min_price") || ""
  );
  const [maxPrice, setMaxPrice] = useState<string>(
    searchParams.get("max_price") || ""
  );
  const [status, setStatus] = useState<string>(
    searchParams.get("status") || ""
  );
  const [sortField, setSortField] = useState<"price" | "name" | "created_at">(
    (searchParams.get("sort_field") as "price" | "name" | "created_at") ||
      "created_at"
  );
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">(
    (searchParams.get("sort_order") as "asc" | "desc") || "desc"
  );
  const [page, setPage] = useState(
    parseInt(searchParams.get("page") || "1", 10)
  );
  const limit = 20;

  const [products, setProducts] = useState<Product[]>([]);
  const [total, setTotal] = useState(0);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load categories
  useEffect(() => {
    getCategories()
      .then(setCategories)
      .catch(() => setCategories([]));
  }, []);

  // Search function
  const performSearch = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const params: any = {
        page,
        limit,
        sort_field: sortField,
        sort_order: sortOrder,
      };

      if (query.trim()) params.q = query.trim();
      if (categoryId) params.category_id = categoryId;
      if (minPrice) params.min_price = parseFloat(minPrice);
      if (maxPrice) params.max_price = parseFloat(maxPrice);
      if (status) params.status = status;

      const response: ProductsResponse = await searchProductsAdvanced(params);
      setProducts(response.products || []);
      setTotal(response.total || 0);
    } catch (err: any) {
      setError(err?.message || "Failed to search products");
      setProducts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [
    query,
    categoryId,
    minPrice,
    maxPrice,
    status,
    sortField,
    sortOrder,
    page,
    limit,
  ]);

  // Update URL when filters change
  useEffect(() => {
    const params = new URLSearchParams();
    if (query) params.set("q", query);
    if (categoryId) params.set("category_id", categoryId.toString());
    if (minPrice) params.set("min_price", minPrice);
    if (maxPrice) params.set("max_price", maxPrice);
    if (status) params.set("status", status);
    if (sortField) params.set("sort_field", sortField);
    if (sortOrder) params.set("sort_order", sortOrder);
    if (page > 1) params.set("page", page.toString());

    router.replace(`/search?${params.toString()}`, { scroll: false });
  }, [
    query,
    categoryId,
    minPrice,
    maxPrice,
    status,
    sortField,
    sortOrder,
    page,
    router,
  ]);

  // Perform search when filters change
  useEffect(() => {
    performSearch();
  }, [performSearch]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    performSearch();
  };

  const handleReset = () => {
    setQuery("");
    setCategoryId(undefined);
    setMinPrice("");
    setMaxPrice("");
    setStatus("");
    setSortField("created_at");
    setSortOrder("desc");
    setPage(1);
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="min-h-screen bg-[#f5f5f5] pb-10">
      <main className="mx-auto max-w-[1200px] px-0 pt-6">
        <div className="grid grid-cols-12 gap-5">
          {/* Sidebar Filters - 2 cols (approx 200px) */}
          <div className="col-span-2 hidden lg:block">
            {/* Categories */}
            <div className="mb-6">
              <h3 className="font-bold text-sm mb-3 flex items-center gap-2">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="w-4 h-4"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                  />
                </svg>
                Tất Cả Danh Mục
              </h3>
              <ul className="text-sm space-y-2 pl-2">
                <li
                  className={`cursor-pointer ${
                    !categoryId ? "text-[#ee4d2d] font-bold" : ""
                  }`}
                  onClick={() => setCategoryId(undefined)}
                >
                  Tất cả
                </li>
                {categories.map((cat) => (
                  <li
                    key={cat.id}
                    className={`cursor-pointer ${
                      categoryId === cat.id ? "text-[#ee4d2d] font-bold" : ""
                    }`}
                    onClick={() => setCategoryId(cat.id)}
                  >
                    {cat.name}
                  </li>
                ))}
              </ul>
            </div>

            {/* Search Filter Header */}
            <h3 className="font-bold text-sm mb-3 flex items-center gap-2 uppercase">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-4 h-4"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 3c2.755 0 5.455.232 8.083.678.533.09.917.556.917 1.096v1.044a2.25 2.25 0 01-.659 1.591l-5.432 5.432a2.25 2.25 0 00-.659 1.591v2.927a2.25 2.25 0 01-1.244 2.013L9.75 21v-6.568a2.25 2.25 0 00-.659-1.591L3.659 7.409A2.25 2.25 0 013 5.818V4.774c0-.54.384-1.006.917-1.096A48.32 48.32 0 0112 3z"
                />
              </svg>
              Bộ lọc tìm kiếm
            </h3>

            {/* Price Range */}
            <div className="mb-6">
              <div className="text-sm mb-2">Khoảng Giá</div>
              <div className="flex items-center gap-2 mb-2">
                <input
                  type="number"
                  placeholder="₫ TỪ"
                  value={minPrice}
                  onChange={(e) => setMinPrice(e.target.value)}
                  className="w-full text-xs px-1 py-1 border outline-none"
                />
                <span className="text-neutral-400">-</span>
                <input
                  type="number"
                  placeholder="₫ ĐẾN"
                  value={maxPrice}
                  onChange={(e) => setMaxPrice(e.target.value)}
                  className="w-full text-xs px-1 py-1 border outline-none"
                />
              </div>
              <button
                onClick={() => performSearch()}
                className="w-full bg-[#ee4d2d] text-white text-sm py-1 uppercase rounded-sm hover:opacity-90"
              >
                Áp dụng
              </button>
            </div>

            {/* Rating (Mock) */}
            <div className="mb-6">
              <div className="text-sm mb-2">Đánh Giá</div>
              <div className="space-y-1 text-sm pl-2">
                {[5, 4, 3, 2, 1].map((star) => (
                  <div
                    key={star}
                    className="flex items-center gap-1 cursor-pointer hover:opacity-80"
                  >
                    <div className="flex text-yellow-500 text-xs">
                      {[...Array(5)].map((_, i) => (
                        <span
                          key={i}
                          className={i < star ? "" : "text-neutral-300"}
                        >
                          ★
                        </span>
                      ))}
                    </div>
                    {star < 5 && <span>trở lên</span>}
                  </div>
                ))}
              </div>
            </div>

            <button
              onClick={handleReset}
              className="w-full bg-[#ee4d2d] text-white text-sm py-1 uppercase rounded-sm hover:opacity-90"
            >
              Xóa tất cả
            </button>
          </div>

          {/* Main Content - 10 cols */}
          <div className="col-span-12 lg:col-span-10">
            {/* Sort Bars */}
            <div className="bg-[#ededed] p-3 flex items-center justify-between text-sm mb-4 rounded-sm">
              <div className="flex items-center gap-2">
                <span className="mr-2">Sắp xếp theo</span>
                <button className="bg-[#ee4d2d] text-white px-4 py-2 rounded-sm">
                  Liên Quan
                </button>
                <button className="bg-white px-4 py-2 rounded-sm hover:bg-neutral-50 ml-1">
                  Mới Nhất
                </button>
                <button className="bg-white px-4 py-2 rounded-sm hover:bg-neutral-50 ml-1">
                  Bán Chạy
                </button>
                <div className="relative inline-block ml-1">
                  <select
                    value={`${sortField}-${sortOrder}`}
                    onChange={(e) => {
                      const [field, order] = e.target.value.split("-");
                      setSortField(field as any);
                      setSortOrder(order as any);
                    }}
                    className="bg-white px-4 py-2 rounded-sm outline-none cursor-pointer appearance-none min-w-[200px]"
                  >
                    <option value="price-asc">Giá: Thấp đến Cao</option>
                    <option value="price-desc">Giá: Cao đến Thấp</option>
                  </select>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <div>
                  <span className="text-[#ee4d2d] font-bold">{page}</span>
                  <span className="text-neutral-800">/{totalPages || 1}</span>
                </div>
                <div className="flex rounded-sm overflow-hidden border border-neutral-200">
                  <button
                    onClick={() => setPage(Math.max(1, page - 1))}
                    disabled={page === 1}
                    className="px-3 py-2 bg-white hover:bg-neutral-50 disabled:bg-neutral-100 disabled:text-neutral-300 border-r"
                  >
                    &lt;
                  </button>
                  <button
                    onClick={() => setPage(Math.min(totalPages, page + 1))}
                    disabled={page === totalPages}
                    className="px-3 py-2 bg-white hover:bg-neutral-50 disabled:bg-neutral-100 disabled:text-neutral-300"
                  >
                    &gt;
                  </button>
                </div>
              </div>
            </div>

            {/* Results */}
            {loading ? (
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2.5">
                {[...Array(10)].map((_, i) => (
                  <ProductCardSkeleton key={i} />
                ))}
              </div>
            ) : (
              <>
                {products.length > 0 ? (
                  <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2.5">
                    {products.map((p) => (
                      <ProductCard key={p.id} product={p} />
                    ))}
                  </div>
                ) : (
                  <div className="flex flex-col items-center justify-center p-20 bg-white shadow-sm rounded-sm">
                    <img
                      src="https://deo.shopeemobile.com/shopee/shopee-pcmall-live-sg/search/a60759ad1dabe909c46a.png"
                      width={130}
                      alt="no result"
                    />
                    <div className="mt-4 text-neutral-500">
                      Không tìm thấy kết quả nào
                    </div>
                    <div className="text-neutral-400 text-sm mt-1">
                      Hãy thử sử dụng các từ khóa khác xem sao.
                    </div>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}

export default function SearchPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen bg-white">
          <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
            <div className="flex items-center justify-center py-24">
              <div className="h-8 w-8 animate-spin rounded-full border-4 border-neutral-200 border-t-neutral-900"></div>
            </div>
          </main>
        </div>
      }
    >
      <SearchContent />
    </Suspense>
  );
}
