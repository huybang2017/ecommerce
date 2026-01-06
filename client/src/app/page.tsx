import { getProducts } from "@/lib/api";
import ProductCard from "@/components/ProductCard";
import Error from "@/components/Error";
import { Product } from "@/lib/types";
import CategorySection from "@/components/home/CategorySection";

export default async function Home() {
  let products: Product[] = [];
  let error: string | null = null;

  try {
    const response = await getProducts({
      page: 1,
      limit: 20,
      status: "ACTIVE",
    });
    products = response.products || [];
  } catch (err: any) {
    error = err?.message || "Failed to load products";
  }

  const flashSaleProducts = products.slice(0, 8);
  const recommendedProducts = products;

  return (
    <div className="min-h-screen bg-[#f5f5f5]">
      <main className="mx-auto max-w-7xl px-4 pb-10 pt-4 sm:px-6 lg:px-8">
        {/* Hero banner + khuyến mãi bên phải */}
        <section className="grid gap-4 lg:grid-cols-[2fr,1fr]">
          <div className="overflow-hidden rounded-lg bg-white shadow-sm">
            <div className="flex h-52 items-center justify-between bg-gradient-to-r from-orange-500 to-pink-500 px-8 text-white">
              <div>
                <p className="text-xs font-medium uppercase tracking-[0.2em] text-white/80">
                  Miễn phí vận chuyển
                </p>
                <h1 className="mt-2 text-2xl font-semibold sm:text-3xl">
                  Săn deal mỗi ngày
                </h1>
                <p className="mt-2 text-sm text-white/85">
                  Hàng nghìn sản phẩm chính hãng. Giảm giá sốc, voucher ngập
                  tràn.
                </p>
              </div>
              <div className="hidden h-40 w-40 rounded-full bg-white/15 sm:block" />
            </div>
          </div>

          <div className="hidden flex-col gap-3 lg:flex">
            <div className="h-24 rounded-lg bg-white p-4 shadow-sm">
              <p className="text-xs font-semibold uppercase tracking-wide text-orange-500">
                Voucher hôm nay
              </p>
              <p className="mt-1 text-sm text-neutral-700">
                Sưu tầm mã giảm đến 50K cho mọi đơn.
              </p>
            </div>
            <div className="h-24 rounded-lg bg-white p-4 shadow-sm">
              <p className="text-xs font-semibold uppercase tracking-wide text-orange-500">
                Freeship Xtra
              </p>
              <p className="mt-1 text-sm text-neutral-700">
                Miễn phí vận chuyển cho đơn từ 0đ.
              </p>
            </div>
          </div>
        </section>

        {/* Danh mục nổi bật */}
        <CategorySection />

        {/* Flash sale */}
        <section className="mt-6 rounded-lg bg-white shadow-sm">
          <div className="flex items-center justify-between border-b px-4 py-3">
            <div className="flex items-baseline gap-2">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-orange-500">
                Flash Sale
              </h2>
              <span className="rounded bg-orange-100 px-1.5 py-0.5 text-[11px] font-semibold text-orange-600">
                ĐANG DIỄN RA
              </span>
            </div>
            <a
              href="/products"
              className="text-xs font-medium text-orange-500 hover:text-orange-600"
            >
              Xem tất cả
            </a>
          </div>

          {error ? (
            <div className="px-4 py-6">
              <Error message={error} />
            </div>
          ) : flashSaleProducts.length === 0 ? (
            <div className="px-4 py-8 text-center text-sm text-neutral-500">
              Chưa có sản phẩm.
            </div>
          ) : (
            <div className="flex gap-3 overflow-x-auto px-4 py-4">
              {flashSaleProducts.map((product) => (
                <div key={product.id} className="w-44 flex-shrink-0">
                  <ProductCard product={product} />
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Gợi ý hôm nay */}
        <section className="mt-6 rounded-lg bg-white shadow-sm">
          <div className="border-b px-4 py-3">
            <h2 className="text-sm font-semibold uppercase tracking-wide text-orange-500">
              Gợi ý hôm nay
            </h2>
          </div>

          {error ? (
            <div className="px-4 py-6">
              <Error message={error} />
            </div>
          ) : recommendedProducts.length === 0 ? (
            <div className="px-4 py-10 text-center text-sm text-neutral-500">
              Chưa có sản phẩm.
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-3 px-2 py-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
              {recommendedProducts.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
