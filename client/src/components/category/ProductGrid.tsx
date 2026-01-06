import Link from "next/link";

interface Product {
  id: number;
  name: string;
  image: string;
  price: number;
  originalPrice?: number;
  rating: number;
  sold: number;
  discountBadge?: string;
  location: string;
  slug: string;
}

interface ProductGridProps {
  products: Product[];
}

export function ProductGrid({ products }: ProductGridProps) {
  if (products.length === 0) {
    return (
      <div className="bg-white p-12 text-center rounded-sm">
        <div className="text-neutral-500 mb-2 text-6xl">🔍</div>
        <p className="text-lg text-neutral-600">Không tìm thấy sản phẩm nào</p>
        <p className="text-sm text-neutral-400 mt-1">
          Hãy thử sử dụng các từ khóa khác hoặc bỏ bớt các bộ lọc
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-2">
      {products.map((product) => (
        <Link
          key={product.id}
          href={`/product/${product.slug}`}
          className="bg-white border border-transparent hover:border-[#ee4d2d] hover:shadow-sm hover:z-10 transition-all rounded-sm overflow-hidden relative group block"
        >
          {/* Badge */}
          {product.discountBadge && (
            <div className="absolute top-0 right-0 bg-[#ffd839] text-[#ee4d2d] px-1 pt-1 pb-1 text-xs font-semibold rounded-bl-sm z-10 flex flex-col items-center leading-none w-9 h-10">
              <span>{product.discountBadge}</span>
              <span className="text-white text-[10px] font-normal uppercase">
                giảm
              </span>
            </div>
          )}

          <div className="relative pt-[100%]">
            <img
              src={product.image}
              alt={product.name}
              className="absolute inset-0 w-full h-full object-cover"
            />
          </div>

          <div className="p-2">
            <div className="text-xs line-clamp-2 min-h-[2rem] mb-2 text-neutral-800">
              {product.name}
            </div>

            <div className="flex items-center gap-1 mb-2 flex-wrap h-[16px]">
              {/* Vouchers or Tags could go here */}
              <div className="border border-[#ee4d2d] text-[#ee4d2d] text-[10px] px-1 h-3 flex items-center leading-none">
                Rẻ Vô Địch
              </div>
            </div>

            <div className="flex items-center justify-between mb-2">
              <div className="flex flex-col">
                {product.originalPrice && (
                  <div className="text-xs text-neutral-400 line-through">
                    ₫{product.originalPrice.toLocaleString("vi-VN")}
                  </div>
                )}
                <div className="text-[#ee4d2d] font-semibold text-base">
                  <span className="text-xs">₫</span>
                  {product.price.toLocaleString("vi-VN")}
                </div>
              </div>
            </div>

            <div className="flex items-center justify-end gap-1 text-[10px] text-neutral-500">
              <div className="flex items-center">
                {Array(5)
                  .fill(0)
                  .map((_, i) => (
                    <svg
                      key={i}
                      className={`w-2.5 h-2.5 ${
                        i < Math.floor(product.rating)
                          ? "text-[#ffce3d]"
                          : "text-neutral-300"
                      }`}
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  ))}
              </div>
              <span>
                Đã bán{" "}
                {product.sold > 1000
                  ? `${(product.sold / 1000).toFixed(1)}k`
                  : product.sold}
              </span>
            </div>

            <div className="text-right text-xs text-neutral-400 mt-1">
              {product.location}
            </div>
          </div>

          {/* Hover Find Similar */}
          <div className="absolute top-full left-0 w-full bg-[#ee4d2d] text-white text-center py-2 text-sm opacity-0 group-hover:opacity-100 group-hover:-translate-y-full transition-all duration-300">
            Tìm sản phẩm tương tự
          </div>
        </Link>
      ))}
    </div>
  );
}
