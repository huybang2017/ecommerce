"use client";

import Link from "next/link";
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useCartContext as useCart } from "@/contexts/CartContext";

export default function Header() {
  const { isAuthenticated, user, logout } = useAuth();
  const { itemCount } = useCart();
  const [keyword, setKeyword] = useState("");
  const router = useRouter();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (keyword.trim()) {
      router.push(`/search?q=${encodeURIComponent(keyword)}`);
    }
  };

  return (
    <header className="sticky top-0 z-50 bg-gradient-to-b from-[#f53d2d] to-[#f63] text-white shadow-md">
      {/* Top Utility Bar */}
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-1 text-[11px] sm:px-6 lg:px-8 font-light">
        <div className="flex gap-4">
          <span className="hover:text-white/80 cursor-pointer">
            Kênh Người Bán
          </span>
          <span className="hover:text-white/80 cursor-pointer">
            Tải ứng dụng
          </span>
          <span className="hover:text-white/80 cursor-pointer">Kết nối</span>
        </div>
        <div className="flex gap-4">
          <span className="hover:text-white/80 cursor-pointer">Thông báo</span>
          <span className="hover:text-white/80 cursor-pointer">Hỗ trợ</span>
          {isAuthenticated ? (
            <div className="flex items-center gap-2 group relative cursor-pointer">
              <span className="font-semibold">{user?.email}</span>
              <div className="absolute right-0 top-full pt-2 hidden group-hover:block w-40">
                <div className="bg-white text-black rounded shadow-md overflow-hidden flex flex-col">
                  <Link
                    href="/orders"
                    className="px-4 py-2 hover:bg-gray-100 text-left"
                  >
                    Đơn mua
                  </Link>
                  <button
                    onClick={logout}
                    className="px-4 py-2 hover:bg-gray-100 text-left w-full"
                  >
                    Đăng xuất
                  </button>
                </div>
              </div>
            </div>
          ) : (
            <div className="flex gap-3 font-semibold">
              <Link href="/register" className="hover:text-white/80">
                Đăng Ký
              </Link>
              <Link href="/login" className="hover:text-white/80">
                Đăng Nhập
              </Link>
            </div>
          )}
        </div>
      </div>

      {/* Main Header */}
      <div className="mx-auto flex max-w-7xl items-center gap-8 px-4 pb-3 pt-2 sm:px-6 lg:px-8">
        <Link
          href="/"
          className="flex flex-col items-center text-2xl font-bold tracking-tight text-white"
        >
          <span className="leading-none">Shopee</span>
          <span className="text-[10px] bg-white text-[#f53d2d] px-1 rounded transform -rotate-6 mt-1 self-end font-bold">
            Clone
          </span>
        </Link>

        <div className="flex-1">
          <form
            onSubmit={handleSearch}
            className="relative shadow-sm rounded-sm bg-white p-[3px]"
          >
            <input
              type="text"
              className="w-full border-none px-3 py-2 text-sm text-black outline-none focus:ring-0"
              placeholder="Đăng ký và nhận voucher bạn mới đến 70k!"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
            />
            <button
              type="submit"
              className="absolute right-1 top-1 bottom-1 w-14 rounded-sm bg-[#fb5533] hover:bg-[#fa4722] flex items-center justify-center text-white"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={2}
                stroke="currentColor"
                className="w-4 h-4"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
                />
              </svg>
            </button>
          </form>
          <div className="mt-1 flex gap-4 text-[11px] text-white/90">
            <span className="cursor-pointer">Váy</span>
            <span className="cursor-pointer">Dép Nữ</span>
            <span className="cursor-pointer">Áo Thun</span>
            <span className="cursor-pointer">Túi Xách</span>
          </div>
        </div>

        <Link
          href="/cart"
          className="relative px-4 text-white hover:text-white/90"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="w-7 h-7"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 00-16.536-1.84M7.5 14.25L5.106 5.272M6 20.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm12.75 0a.75.75 0 11-1.5 0 .75.75 0 011.5 0z"
            />
          </svg>
          {itemCount > 0 && (
            <span className="absolute right-2 -top-1 flex h-5 w-fit min-w-[20px] items-center justify-center rounded-full border border-[#ee4d2d] bg-white px-1 text-center text-xs font-bold text-[#ee4d2d]">
              {itemCount > 99 ? "99+" : itemCount}
            </span>
          )}
        </Link>
      </div>
    </header>
  );
}
