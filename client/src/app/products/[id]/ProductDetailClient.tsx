"use client";

import { useState } from "react";
import Link from "next/link";
import { Product } from "@/types/product";
import { useCartContext } from "@/contexts/CartContext";
import { useProductDetail } from "@/hooks/useProductDetail";
import PriceSection from "./components/PriceSection";
import VariationSelector from "./components/VariationSelector";
import QuantitySelector from "./components/QuantitySelector";
import ProductGallery from "./components/ProductGallery";

interface ProductDetailClientProps {
  initialProduct: Product;
}

export default function ProductDetailClient({
  initialProduct,
}: ProductDetailClientProps) {
  const { addItem } = useCartContext();
  const [adding, setAdding] = useState(false);

  // Use custom hook for product detail logic
  const {
    productItems,
    selectedSKU,
    quantity,
    isLoadingSKUs,
    setSelectedSKU,
    setQuantity,
    canAddToCart,
  } = useProductDetail(initialProduct.id);

  const handleAddToCart = async () => {
    if (!canAddToCart) {
      alert("Vui lòng chọn phiên bản sản phẩm");
      return;
    }

    setAdding(true);
    try {
      await addItem(initialProduct.id, quantity);
      setQuantity(1);
    } catch (err) {
      console.error("Failed to add to cart:", err);
      alert("Không thể thêm vào giỏ hàng. Vui lòng thử lại.");
    } finally {
      setAdding(false);
    }
  };

  const handleBuyNow = async () => {
    await handleAddToCart();
    window.location.href = "/cart";
  };

  return (
    <div className="min-h-screen bg-[#f5f5f5] pb-10">
      <main className="mx-auto max-w-7xl px-4 pt-6 sm:px-6 lg:px-8">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-neutral-600 mb-4">
          <Link href="/" className="hover:text-[#ee4d2d]">
            Shopee
          </Link>
          <span>&gt;</span>
          <Link href="/products" className="hover:text-[#ee4d2d]">
            {initialProduct.category?.name || "Danh mục"}
          </Link>
          <span>&gt;</span>
          <span className="truncate max-w-xs">{initialProduct.name}</span>
        </div>

        {/* Main Product Section */}
        <div className="bg-white rounded shadow-sm p-4 grid grid-cols-1 md:grid-cols-12 gap-8">
          {/* Gallery (Left - 5 cols) */}
          <ProductGallery
            images={initialProduct.images || []}
            name={initialProduct.name}
            activeSkuImage={selectedSKU?.image_url}
          />

          {/* Info (Right - 7 cols) */}
          <div className="md:col-span-7 flex flex-col gap-6">
            {/* Product Title */}
            <div>
              <h1 className="text-xl font-medium text-neutral-800 line-clamp-2 leading-tight">
                <span className="inline-block bg-[#ee4d2d] text-white text-[10px] px-1 mr-2 rounded-sm align-middle mb-0.5">
                  Yêu Thích+
                </span>
                {initialProduct.name}
              </h1>

              {/* Rating & Stats */}
              <div className="flex items-center gap-4 mt-3 text-sm">
                <span className="text-[#ee4d2d] border-b border-[#ee4d2d] px-1 cursor-pointer">
                  4.9
                </span>
                <div className="flex text-[#ee4d2d] text-xs">★★★★★</div>
                <div className="w-px bg-neutral-300 h-4"></div>
                <div className="flex gap-1">
                  <span className="border-b border-black text-neutral-700 font-medium">
                    10,2k
                  </span>
                  <span className="text-neutral-500">Đánh giá</span>
                </div>
                <div className="w-px bg-neutral-300 h-4"></div>
                <div className="flex gap-1">
                  <span className="text-neutral-700 font-medium">
                    {initialProduct.sold_count || 0}
                  </span>
                  <span className="text-neutral-500">Đã bán</span>
                </div>
              </div>
            </div>

            {/* Price */}
            <PriceSection
              basePrice={initialProduct.base_price}
              currentPrice={selectedSKU?.price}
              isLoading={isLoadingSKUs}
            />

            {/* Shipping & Warranty Mock */}
            <div className="grid grid-cols-[110px_1fr] gap-4 text-sm items-center">
              <div className="text-neutral-500">Bảo Hành</div>
              <div>Bảo hiểm Thời trang &gt;</div>

              <div className="text-neutral-500">Vận Chuyển</div>
              <div className="flex items-center gap-2">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src="https://deo.shopeemobile.com/shopee/shopee-pcmall-live-sg/productdetailspage/d9e992985b18d96aab90.png"
                  width={20}
                  alt="freeship"
                />
                <span>Miễn phí vận chuyển</span>
              </div>
            </div>

            {/* Variation Selector (Shopee-style) */}
            <VariationSelector
              productId={initialProduct.id}
              selectedSKU={selectedSKU}
              onSKUChange={setSelectedSKU}
              isLoading={isLoadingSKUs}
            />

            {/* Quantity Selector */}
            <div className="grid grid-cols-[110px_1fr] gap-4 items-center">
              <div className="text-neutral-500 text-sm">Số Lượng</div>
              <QuantitySelector
                value={quantity}
                onChange={setQuantity}
                max={selectedSKU?.qty_in_stock}
                disabled={!selectedSKU}
              />
            </div>

            {/* Action Buttons */}
            <div className="flex gap-4 mt-4">
              <button
                onClick={handleAddToCart}
                disabled={!canAddToCart || adding}
                className="flex-1 max-w-50 border border-[#ee4d2d] bg-[#ff57221a] text-[#ee4d2d] h-12 flex items-center justify-center gap-2 rounded-sm hover:bg-[#ff57222a] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
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
                {adding ? "Đang thêm..." : "Thêm Vào Giỏ"}
              </button>

              <button
                onClick={handleBuyNow}
                disabled={!canAddToCart || adding}
                className="flex-1 max-w-50 bg-[#ee4d2d] text-white h-12 rounded-sm hover:bg-[#d73211] shadow-sm disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                Mua Ngay
              </button>
            </div>
          </div>
        </div>

        {/* Shop Info */}
        <div className="bg-white rounded shadow-sm p-4 mt-4 flex items-center gap-6">
          <div className="relative">
            <div className="w-20 h-20 rounded-full bg-neutral-200 border border-neutral-300 overflow-hidden">
              {/* eslint-disable-next-line @next/next/no-img-element */}
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
              <button className="border border-[#ee4d2d] bg-[#ff57221a] text-[#ee4d2d] px-3 py-1.5 text-xs rounded-sm flex items-center gap-1 hover:bg-[#ff57222a] transition-colors">
                Chat Ngay
              </button>
              <button className="border border-neutral-300 text-neutral-600 px-3 py-1.5 text-xs rounded-sm hover:bg-neutral-50 flex items-center gap-1 transition-colors">
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

        {/* Description */}
        <div className="bg-white rounded shadow-sm p-6 mt-4">
          <h2 className="bg-neutral-50 p-3 text-lg font-medium text-neutral-800 uppercase mb-4">
            Chi tiết sản phẩm
          </h2>
          <div className="grid grid-cols-[140px_1fr] gap-y-3 text-sm mb-6">
            <div className="text-neutral-500">Danh Mục</div>
            <div className="text-blue-600">
              <Link href="/" className="hover:underline">
                Shopee
              </Link>{" "}
              &gt;{" "}
              <Link href="/products" className="hover:underline">
                {initialProduct.category?.name}
              </Link>
            </div>
            {selectedSKU && (
              <>
                <div className="text-neutral-500">Mã SKU</div>
                <div>{selectedSKU.sku_code}</div>
                <div className="text-neutral-500">Kho hàng</div>
                <div>{selectedSKU.qty_in_stock} sản phẩm</div>
              </>
            )}
            <div className="text-neutral-500">Gửi từ</div>
            <div>Hà Nội</div>
          </div>

          <h2 className="bg-neutral-50 p-3 text-lg font-medium text-neutral-800 uppercase mb-4">
            Mô tả sản phẩm
          </h2>
          <div className="text-sm leading-7 text-neutral-700 whitespace-pre-line px-2">
            {initialProduct.description || "Chưa có mô tả cho sản phẩm này."}
          </div>
        </div>
      </main>
    </div>
  );
}
