'use client';

import Link from 'next/link';
import React from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useCartContext as useCart } from '@/contexts/CartContext';

export default function Header() {
  const { isAuthenticated, user, logout } = useAuth();
  const { itemCount } = useCart();

  return (
    <header className="sticky top-0 z-50 border-b border-neutral-200 bg-white/95 backdrop-blur-sm">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-20 items-center justify-between">
          <Link
            href="/"
            className="text-2xl font-semibold tracking-tight text-neutral-900 transition-opacity hover:opacity-80"
          >
            Ecommerce
          </Link>
          <nav className="flex items-center gap-8">
            <Link
              href="/"
              className="text-sm font-medium text-neutral-600 transition-colors hover:text-neutral-900"
            >
              Home
            </Link>
            <Link
              href="/products"
              className="text-sm font-medium text-neutral-600 transition-colors hover:text-neutral-900"
            >
              Products
            </Link>
            <Link
              href="/search"
              className="text-sm font-medium text-neutral-600 transition-colors hover:text-neutral-900"
            >
              Search
            </Link>
            <Link
              href="/cart"
              className="relative text-sm font-medium text-neutral-600 transition-colors hover:text-neutral-900"
            >
              Cart
              {itemCount > 0 && (
                <span className="absolute -right-2 -top-2 flex h-5 w-5 items-center justify-center rounded-full bg-neutral-900 text-xs font-semibold text-white">
                  {itemCount > 9 ? '9+' : itemCount}
                </span>
              )}
            </Link>
            <Link
              href="/orders"
              className="text-sm font-medium text-neutral-600 transition-colors hover:text-neutral-900"
            >
              Orders
            </Link>
            {isAuthenticated ? (
              <>
                <span className="text-sm text-neutral-500">
                  {user?.email}
                </span>
                <button
                  onClick={logout}
                  className="rounded-lg px-4 py-2 text-sm font-medium text-neutral-600 transition-colors hover:bg-neutral-100 hover:text-neutral-900"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link
                  href="/login"
                  className="rounded-lg px-4 py-2 text-sm font-medium text-neutral-600 transition-colors hover:bg-neutral-100 hover:text-neutral-900"
                >
                  Login
                </Link>
                <Link
                  href="/register"
                  className="rounded-lg bg-neutral-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-neutral-800"
                >
                  Register
                </Link>
              </>
            )}
          </nav>
        </div>
      </div>
    </header>
  );
}
