"use client";

import { getProduct } from "@/lib/api";
import Error from "@/components/Error";
import { notFound, useRouter } from "next/navigation";
import Link from "next/link";
import { useCartContext as useCart } from "@/contexts/CartContext";
import { useState, useEffect } from "react";
import { Product } from "@/lib/types";

interface ProductDetailPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default function ProductDetailPage({ params }: ProductDetailPageProps) {
  const router = useRouter();
  const { addItem } = useCart();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [productId, setProductId] = useState<number | null>(null);
  const [activeImage, setActiveImage] = useState<string>("");
  const [quantity, setQuantity] = useState(1);

  const handleQuantityUpdate = (val: number) => {
    if (val < 1) return;
    if (product && val > product.stock) return;
    setQuantity(val);
  };

  useEffect(() => {
    async function loadProduct() {
      const resolvedParams = await params;
      const id = parseInt(resolvedParams.id, 10);

      if (isNaN(id)) {
        setError("Invalid product ID");
        setLoading(false);
        return;
      }

      setProductId(id);

      try {
        const productData = await getProduct(id);
        setProduct(productData);
        if (productData.images && productData.images.length > 0) {
          setActiveImage(productData.images[0]);
        }
      } catch (err: any) {
        if (err?.message?.includes("404")) {
          notFound();
        }
        setError(err?.message || "Failed to load product");
      } finally {
        setLoading(false);
      }
    }

    loadProduct();
  }, [params]);

  const handleAddToCart = async () => {
    if (!product || product.status !== "ACTIVE" || product.stock === 0) return;

    setAdding(true);
    try {
      await addItem({
        product_id: product.id,
        name: product.name,
        price: product.price,
        quantity: 1,
        image:
          product.images &&
          Array.isArray(product.images) &&
          product.images.length > 0
            ? product.images[0]
            : undefined,
        sku: product.sku,
      });
    } catch (err) {
      console.error("Failed to add to cart:", err);
    } finally {
      setAdding(false);
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(price);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <main className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
          <div className="flex items-center justify-center py-24">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-neutral-200 border-t-neutral-900"></div>
          </div>
        </main>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-white">
        <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
          <Error message={error} />
        </main>
      </div>
    );
  }

  if (!product) {
    notFound();
  }

  const images =
    product.images && Array.isArray(product.images) && product.images.length > 0
      ? product.images
      : ["/placeholder-product.jpg"];

  return (
    <div className="min-h-screen bg-[#f5f5f5] pb-10">
      <main className="mx-auto max-w-7xl px-4 pt-6 sm:px-6 lg:px-8">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-neutral-600 mb-4">
          <Link href="/">Shopee</Link>
          <span>&gt;</span>
          <Link href="/products">{product.category?.name || "Danh mục"}</Link>
          <span>&gt;</span>
          <span className="truncate max-w-xs">{product.name}</span>
        </div>

        {/* Main Product Section */}
        <div className="bg-white rounded shadow-sm p-4 grid grid-cols-1 md:grid-cols-12 gap-8">
          {/* Gallery (Left - 5 cols) */}
          <div className="md:col-span-5 flex flex-col gap-4">
            <div className="relative w-full pt-[100%] border border-neutral-100 rounded-sm overflow-hidden group hover:shadow-sm">
              <img
                src={activeImage || images[0]}
                alt={product.name}
                className="absolute top-0 left-0 w-full h-full object-contain"
              />
            </div>
            <div className="flex gap-2 overflow-x-auto pb-2">
              {images.map((img, idx) => (
                <button
                  key={idx}
                  onMouseEnter={() => setActiveImage(img)}
                  className={`flex-shrink-0 w-20 h-20 border-2 rounded-sm overflow-hidden ${
                    activeImage === img
                      ? "border-[#ee4d2d]"
                      : "border-transparent hover:border-[#ee4d2d]"
                  }`}
                >
                  <img
                    src={img}
                    alt="thumbnail"
                    className="w-full h-full object-cover"
                  />
                </button>
              ))}
            </div>
            <div className="flex justify-center gap-6 text-sm text-neutral-600 mt-2">
              <button className="flex items-center gap-1 hover:text-[#ee4d2d]">
                <span>Chia sẻ:</span>
                {/* Mock Icons */}
                <div className="w-5 h-5 bg-blue-600 rounded-full"></div>
                <div className="w-5 h-5 bg-blue-400 rounded-full"></div>
                <div className="w-5 h-5 bg-red-500 rounded-full"></div>
              </button>
              <div className="w-[1px] bg-neutral-300 h-5"></div>
              <button className="flex items-center gap-1 hover:text-[#ee4d2d]">
                <span>❤️ Đã thích (2,1k)</span>
              </button>
            </div>
          </div>

          {/* Info (Right - 7 cols) */}
          <div className="md:col-span-7 flex flex-col gap-6">
            <div>
              <h1 className="text-xl font-medium text-neutral-800 line-clamp-2 leading-tight">
                <span className="inline-block bg-[#ee4d2d] text-white text-[10px] px-1 mr-2 rounded-sm align-middle mb-0.5">
                  Yêu Thích+
                </span>
                {product.name}
              </h1>
              <div className="flex items-center gap-4 mt-3 text-sm">
                <span className="text-[#ee4d2d] border-b border-[#ee4d2d] px-1 cursor-pointer">
                  4.9
                </span>
                <div className="flex text-[#ee4d2d] text-xs">★★★★★</div>
                <div className="w-[1px] bg-neutral-300 h-4"></div>
                <div className="flex gap-1">
                  <span className="border-b border-black text-neutral-700 font-medium">
                    10,2k
                  </span>
                  <span className="text-neutral-500">Đánh giá</span>
                </div>
                <div className="w-[1px] bg-neutral-300 h-4"></div>
                <div className="flex gap-1">
                  <span className="text-neutral-700 font-medium">29,1k</span>
                  <span className="text-neutral-500">Đã bán</span>
                </div>
              </div>
            </div>

            {/* Price Box */}
            <div className="bg-[#fafafa] p-4 flex items-center gap-4 rounded-sm">
              <div className="text-neutral-400 line-through">
                ₫{new Intl.NumberFormat("vi-VN").format(product.price * 1.3)}
              </div>
              <div className="text-3xl font-medium text-[#ee4d2d]">
                ₫{new Intl.NumberFormat("vi-VN").format(product.price)}
              </div>
              <div className="bg-[#ee4d2d] text-white text-xs font-bold px-1 py-0.5 rounded-sm uppercase">
                Giảm 30%
              </div>
            </div>

            {/* Options/Variations Mock */}
            <div className="grid grid-cols-[110px_1fr] gap-4 text-sm items-center">
              <div className="text-neutral-500">Bảo Hành</div>
              <div>Bảo hiểm Thời trang &gt;</div>

              <div className="text-neutral-500">Vận Chuyển</div>
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <img
                    src="https://deo.shopeemobile.com/shopee/shopee-pcmall-live-sg/productdetailspage/d9e992985b18d96aab90.png"
                    width={20}
                    alt="freeship"
                  />
                  <span>Miễn phí vận chuyển</span>
                </div>
              </div>

              <div className="text-neutral-500 mt-2 self-start pt-2">
                Màu Sắc
              </div>
              <div className="flex flex-wrap gap-2">
                {["Trắng", "Đen", "Xám", "Xanh Dương"].map((color) => (
                  <button
                    key={color}
                    className="border border-neutral-200 px-4 py-2 hover:border-[#ee4d2d] hover:text-[#ee4d2d] bg-white rounded-sm min-w-[80px]"
                  >
                    {color}
                  </button>
                ))}
              </div>

              <div className="text-neutral-500 mt-2 self-start pt-2">
                Kích Cỡ
              </div>
              <div className="flex flex-wrap gap-2">
                {["S", "M", "L", "XL", "2XL"].map((size) => (
                  <button
                    key={size}
                    className="border border-neutral-200 px-4 py-2 hover:border-[#ee4d2d] hover:text-[#ee4d2d] bg-white rounded-sm min-w-[60px]"
                  >
                    {size}
                  </button>
                ))}
              </div>

              <div className="text-neutral-500 mt-4">Số Lượng</div>
              <div className="flex items-center gap-4 mt-4">
                <div className="flex items-center border border-neutral-300 rounded-sm">
                  <button
                    onClick={() => handleQuantityUpdate(quantity - 1)}
                    className="w-8 h-8 flex items-center justify-center border-r hover:bg-neutral-50 text-neutral-600"
                  >
                    -
                  </button>
                  <input
                    type="text"
                    value={quantity}
                    readOnly
                    className="w-12 h-8 text-center text-neutral-900 border-none outline-none"
                  />
                  <button
                    onClick={() => handleQuantityUpdate(quantity + 1)}
                    className="w-8 h-8 flex items-center justify-center border-l hover:bg-neutral-50 text-neutral-600"
                  >
                    +
                  </button>
                </div>
                <div className="text-neutral-500 text-sm">
                  {product.stock} sản phẩm có sẵn
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-4 mt-4">
              <button
                onClick={() => handleAddToCart()}
                className="flex-1 max-w-[200px] border border-[#ee4d2d] bg-[#ff57221a] text-[#ee4d2d] h-12 flex items-center justify-center gap-2 rounded-sm hover:bg-[#ff57222a]"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="w-5 h-5"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 00-16.536-1.84M7.5 14.25L5.106 5.272M6 20.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm12.75 0a.75.75 0 11-1.5 0 .75.75 0 011.5 0z"
                  />
                </svg>
                Thêm Vào Giỏ Hàng
              </button>
              <button className="flex-1 max-w-[200px] bg-[#ee4d2d] text-white h-12 rounded-sm hover:bg-[#d73211] shadow-sm">
                Mua Ngay
              </button>
              {adding && (
                <span className="text-sm text-green-600 self-center">
                  Đã thêm!
                </span>
              )}
            </div>
          </div>
        </div>

        {/* Shop Info Card */}
        <div className="bg-white rounded shadow-sm p-4 mt-4 flex items-center gap-6">
          <div className="relative">
            <div className="w-20 h-20 rounded-full bg-neutral-200 border border-neutral-300 overflow-hidden">
              <img
                src="https://down-cvs-vn.img.susercontent.com/vn-11134259-7r98o-lz4b6r1q1z1z8d_tn"
                alt="shop"
                className="w-full h-full object-cover"
              />
            </div>
            <div className="absolute bottom-0 right-0 bg-[#ee4d2d] text-white text-[10px] px-1 rounded-sm">
              Yêu Thích+
            </div>
          </div>
          <div className="flex-1 border-r border-neutral-200 pr-6">
            <div className="font-medium text-neutral-800">
              Official Store Việt Nam
            </div>
            <div className="text-sm text-neutral-500 mb-2">
              Online 15 phút trước
            </div>
            <div className="flex gap-2">
              <button className="border border-[#ee4d2d] bg-[#ff57221a] text-[#ee4d2d] px-3 py-1.5 text-xs rounded-sm flex items-center gap-1">
                Chat Ngay
              </button>
              <button className="border border-neutral-300 text-neutral-600 px-3 py-1.5 text-xs rounded-sm hover:bg-neutral-50 flex items-center gap-1">
                Xem Shop
              </button>
            </div>
          </div>
          <div className="flex gap-12 px-6 text-sm text-neutral-600">
            <div className="space-y-3">
              <div className="flex gap-2">
                <span className="text-neutral-500">Đánh Giá</span>
                <span className="text-[#ee4d2d]">123,4k</span>
              </div>
              <div className="flex gap-2">
                <span className="text-neutral-500">Sản Phẩm</span>
                <span className="text-[#ee4d2d]">213</span>
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex gap-2">
                <span className="text-neutral-500">Tỉ Lệ Phản Hồi</span>
                <span className="text-[#ee4d2d]">98%</span>
              </div>
              <div className="flex gap-2">
                <span className="text-neutral-500">Thời Gian Phản Hồi</span>
                <span className="text-[#ee4d2d]">trong vài giờ</span>
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex gap-2">
                <span className="text-neutral-500">Tham Gia</span>
                <span className="text-[#ee4d2d]">4 năm trước</span>
              </div>
              <div className="flex gap-2">
                <span className="text-neutral-500">Người Theo Dõi</span>
                <span className="text-[#ee4d2d]">543,2k</span>
              </div>
            </div>
          </div>
        </div>

        {/* Description Section */}
        <div className="bg-white rounded shadow-sm p-6 mt-4">
          <h2 className="bg-neutral-50 p-3 text-lg font-medium text-neutral-800 uppercase mb-4">
            Chi tiết sản phẩm
          </h2>
          <div className="grid grid-cols-[140px_1fr] gap-y-3 text-sm mb-6">
            <div className="text-neutral-500">Danh Mục</div>
            <div className="text-blue-600 flex gap-1">
              <Link href="/">Shopee</Link> &gt;{" "}
              <Link href="/products">{product.category?.name}</Link> &gt;{" "}
              <span>{product.name}</span>
            </div>
            <div className="text-neutral-500">Kho hàng</div>
            <div>{product.stock}</div>
            <div className="text-neutral-500">Gửi từ</div>
            <div>Hà Nội</div>
          </div>

          <h2 className="bg-neutral-50 p-3 text-lg font-medium text-neutral-800 uppercase mb-4">
            Mô tả sản phẩm
          </h2>
          <div className="text-sm leading-7 text-neutral-700 whitespace-pre-line px-2">
            {product.description || "Chưa có mô tả cho sản phẩm này."}
            {/* Mock description content if empty */}
            <br />
            <br />
            ✨ ĐẶC ĐIỂM NỔI BẬT:
            <br />
            - Thiết kế hiện đại, trẻ trung.
            <br />
            - Chất liệu cao cấp, bền đẹp theo thời gian.
            <br />
            - Phù hợp đi chơi, đi làm, dự tiệc.
            <br />
            <br />
            🔰 HƯỚNG DẪN BẢO QUẢN:
            <br />
            - Giặt ở nhiệt độ thường.
            <br />
            - Không dùng hóa chất tẩy.
            <br />
            - Ủi ở nhiệt độ thích hợp.
            <br />
            <br />
            CAM KẾT SHOP:
            <br />
            ✅ Hàng chính hãng 100%
            <br />
            ✅ Hoàn tiền nếu sản phẩm không giống mô tả
            <br />✅ Hỗ trợ đổi trả theo quy định của Shopee
          </div>
        </div>
      </main>
    </div>
  );
}
