"use client";

import { useCategories } from "@/hooks/useProducts";
import Link from "next/link";

export default function CategorySection() {
  const { data: categories, isLoading, isError } = useCategories();

  if (isError) return null;

  return (
    <section className="mt-6 rounded-lg bg-white p-4 shadow-sm">
      <div className="mb-3 flex items-center justify-between">
        <h2 className="text-sm font-semibold uppercase tracking-wide text-neutral-800">
          Danh mục
        </h2>
      </div>

      <div className="flex overflow-x-auto gap-3 pb-2 scrollbar-hide">
        {isLoading
          ? Array.from({ length: 15 }).map((_, i) => (
              <div
                key={i}
                className="flex min-w-[120px] shrink-0 flex-col items-center gap-2 rounded-md border border-neutral-100 bg-neutral-50 px-2 py-3 animate-pulse"
              >
                <div className="h-14 w-14 rounded-full bg-neutral-200" />
                <div className="h-3 w-20 bg-neutral-200 rounded" />
              </div>
            ))
          : categories
              ?.filter((cat) => cat.is_active)
              .map((category) => (
                <Link
                  key={category.id}
                  href={`/${category.slug}-cat.${category.id}`}
                  className="flex min-w-[120px] shrink-0 flex-col items-center gap-2 rounded-md border border-neutral-100 bg-neutral-50 px-2 py-3 text-center text-[13px] text-neutral-700 hover:border-orange-400 hover:bg-orange-50 transition-colors"
                >
                  {category.image_url ? (
                    <img
                      src={category.image_url}
                      alt={category.name}
                      className="h-14 w-14 rounded-full object-cover border border-neutral-100"
                    />
                  ) : (
                    <span className="flex h-14 w-14 items-center justify-center rounded-full bg-orange-100 text-lg font-medium text-orange-600">
                      {category.name.charAt(0)}
                    </span>
                  )}
                  <span className="line-clamp-2 leading-tight h-8 flex items-center justify-center">
                    {category.name}
                  </span>
                </Link>
              ))}
      </div>
    </section>
  );
}
