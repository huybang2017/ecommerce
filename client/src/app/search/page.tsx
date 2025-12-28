'use client';

import { useState, useEffect, useCallback } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { searchProductsAdvanced, getCategories } from '@/lib/api';
import ProductCard from '@/components/ProductCard';
import { ProductCardSkeleton } from '@/components/Loading';
import Error from '@/components/Error';
import { Product, Category, SearchResponse } from '@/lib/types';

export default function SearchPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  
  // State
  const [query, setQuery] = useState(searchParams.get('q') || '');
  const [categoryId, setCategoryId] = useState<number | undefined>(
    searchParams.get('category_id') ? parseInt(searchParams.get('category_id')!, 10) : undefined
  );
  const [minPrice, setMinPrice] = useState<string>(
    searchParams.get('min_price') || ''
  );
  const [maxPrice, setMaxPrice] = useState<string>(
    searchParams.get('max_price') || ''
  );
  const [status, setStatus] = useState<string>(
    searchParams.get('status') || ''
  );
  const [sortField, setSortField] = useState<'price' | 'name' | 'created_at'>(
    (searchParams.get('sort_field') as 'price' | 'name' | 'created_at') || 'created_at'
  );
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>(
    (searchParams.get('sort_order') as 'asc' | 'desc') || 'desc'
  );
  const [page, setPage] = useState(parseInt(searchParams.get('page') || '1', 10));
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

      const response: SearchResponse = await searchProductsAdvanced(params);
      setProducts(response.products || []);
      setTotal(response.total || 0);
    } catch (err: any) {
      setError(err?.message || 'Failed to search products');
      setProducts([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [query, categoryId, minPrice, maxPrice, status, sortField, sortOrder, page, limit]);

  // Update URL when filters change
  useEffect(() => {
    const params = new URLSearchParams();
    if (query) params.set('q', query);
    if (categoryId) params.set('category_id', categoryId.toString());
    if (minPrice) params.set('min_price', minPrice);
    if (maxPrice) params.set('max_price', maxPrice);
    if (status) params.set('status', status);
    if (sortField) params.set('sort_field', sortField);
    if (sortOrder) params.set('sort_order', sortOrder);
    if (page > 1) params.set('page', page.toString());

    router.replace(`/search?${params.toString()}`, { scroll: false });
  }, [query, categoryId, minPrice, maxPrice, status, sortField, sortOrder, page, router]);

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
    setQuery('');
    setCategoryId(undefined);
    setMinPrice('');
    setMaxPrice('');
    setStatus('');
    setSortField('created_at');
    setSortOrder('desc');
    setPage(1);
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="min-h-screen bg-white">
      <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="mb-10 text-center">
          <h1 className="text-4xl font-semibold tracking-tight text-neutral-900 sm:text-5xl">
            Search Products
          </h1>
          <p className="mt-4 text-lg text-neutral-600">
            Find exactly what you're looking for
          </p>
        </div>

        {/* Search Form */}
        <form onSubmit={handleSearch} className="mb-10 rounded-2xl border border-neutral-200 bg-neutral-50 p-8">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
            {/* Search Query */}
            <div className="lg:col-span-2">
              <label htmlFor="query" className="mb-2 block text-sm font-medium text-neutral-700">
                Search
              </label>
              <input
                id="query"
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search products..."
                className="w-full rounded-lg border border-neutral-300 bg-white px-4 py-3 text-neutral-900 placeholder:text-neutral-400 focus:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-neutral-200"
              />
            </div>

            {/* Category Filter */}
            <div>
              <label htmlFor="category" className="mb-2 block text-sm font-medium text-neutral-700">
                Category
              </label>
              <select
                id="category"
                value={categoryId || ''}
                onChange={(e) => setCategoryId(e.target.value ? parseInt(e.target.value, 10) : undefined)}
                className="w-full rounded-lg border border-neutral-300 bg-white px-4 py-3 text-neutral-900 focus:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-neutral-200"
              >
                <option value="">All Categories</option>
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Status Filter */}
            <div>
              <label htmlFor="status" className="mb-2 block text-sm font-medium text-neutral-700">
                Status
              </label>
              <select
                id="status"
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                className="w-full rounded-lg border border-neutral-300 bg-white px-4 py-3 text-neutral-900 focus:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-neutral-200"
              >
                <option value="">All Status</option>
                <option value="ACTIVE">Active</option>
                <option value="INACTIVE">Inactive</option>
              </select>
            </div>

            {/* Price Range */}
            <div>
              <label htmlFor="min_price" className="mb-2 block text-sm font-medium text-neutral-700">
                Min Price
              </label>
              <input
                id="min_price"
                type="number"
                value={minPrice}
                onChange={(e) => setMinPrice(e.target.value)}
                placeholder="0"
                min="0"
                step="0.01"
                className="w-full rounded-lg border border-neutral-300 bg-white px-4 py-3 text-neutral-900 placeholder:text-neutral-400 focus:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-neutral-200"
              />
            </div>

            <div>
              <label htmlFor="max_price" className="mb-2 block text-sm font-medium text-neutral-700">
                Max Price
              </label>
              <input
                id="max_price"
                type="number"
                value={maxPrice}
                onChange={(e) => setMaxPrice(e.target.value)}
                placeholder="No limit"
                min="0"
                step="0.01"
                className="w-full rounded-lg border border-neutral-300 bg-white px-4 py-3 text-neutral-900 placeholder:text-neutral-400 focus:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-neutral-200"
              />
            </div>

            {/* Sort Field */}
            <div>
              <label htmlFor="sort_field" className="mb-2 block text-sm font-medium text-neutral-700">
                Sort By
              </label>
              <select
                id="sort_field"
                value={sortField}
                onChange={(e) => setSortField(e.target.value as 'price' | 'name' | 'created_at')}
                className="w-full rounded-lg border border-neutral-300 bg-white px-4 py-3 text-neutral-900 focus:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-neutral-200"
              >
                <option value="created_at">Date</option>
                <option value="price">Price</option>
                <option value="name">Name</option>
              </select>
            </div>

            {/* Sort Order */}
            <div>
              <label htmlFor="sort_order" className="mb-2 block text-sm font-medium text-neutral-700">
                Order
              </label>
              <select
                id="sort_order"
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
                className="w-full rounded-lg border border-neutral-300 bg-white px-4 py-3 text-neutral-900 focus:border-neutral-400 focus:outline-none focus:ring-2 focus:ring-neutral-200"
              >
                <option value="desc">Descending</option>
                <option value="asc">Ascending</option>
              </select>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="mt-6 flex gap-3">
            <button
              type="submit"
              className="rounded-lg bg-neutral-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-neutral-800 focus:outline-none focus:ring-2 focus:ring-neutral-200"
            >
              Search
            </button>
            <button
              type="button"
              onClick={handleReset}
              className="rounded-lg border border-neutral-300 bg-white px-6 py-3 text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50 hover:border-neutral-400"
            >
              Reset
            </button>
          </div>
        </form>

        {/* Results */}
        {error ? (
          <Error message={error} />
        ) : (
          <>
            <div className="mb-6 flex items-center justify-between">
              <p className="text-sm font-medium text-neutral-600">
                {loading ? 'Searching...' : `Found ${total} product${total !== 1 ? 's' : ''}`}
              </p>
            </div>

            {loading ? (
              <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {[...Array(8)].map((_, i) => (
                  <ProductCardSkeleton key={i} />
                ))}
              </div>
            ) : products.length === 0 ? (
              <div className="rounded-xl border border-neutral-200 bg-neutral-50 p-16 text-center">
                <p className="text-lg font-medium text-neutral-500">
                  No products found. Try adjusting your search criteria.
                </p>
              </div>
            ) : (
              <>
                <div className="mb-10 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                  {products.map((product) => (
                    <ProductCard key={product.id} product={product} />
                  ))}
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="flex items-center justify-center gap-3">
                    <button
                      onClick={() => setPage(Math.max(1, page - 1))}
                      disabled={page === 1}
                      className="rounded-lg border border-neutral-200 bg-white px-5 py-2.5 text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50 hover:border-neutral-300 disabled:opacity-40 disabled:cursor-not-allowed"
                    >
                      Previous
                    </button>
                    <span className="px-5 py-2.5 text-sm font-medium text-neutral-600">
                      Page {page} of {totalPages}
                    </span>
                    <button
                      onClick={() => setPage(Math.min(totalPages, page + 1))}
                      disabled={page === totalPages}
                      className="rounded-lg border border-neutral-200 bg-white px-5 py-2.5 text-sm font-medium text-neutral-700 transition-colors hover:bg-neutral-50 hover:border-neutral-300 disabled:opacity-40 disabled:cursor-not-allowed"
                    >
                      Next
                    </button>
                  </div>
                )}
              </>
            )}
          </>
        )}
      </main>
    </div>
  );
}


