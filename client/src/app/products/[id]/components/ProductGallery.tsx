"use client";

import { useEffect, useState } from "react";

interface ProductGalleryProps {
  images: string[];
  name: string;
  activeSkuImage?: string;
}

export default function ProductGallery({
  images,
  name,
  activeSkuImage,
}: ProductGalleryProps) {
  const displayImages =
    images && images.length > 0 ? images : ["/placeholder-product.jpg"];

  // Use activeSkuImage if available, otherwise use first product image
  const activeImage = activeSkuImage || displayImages[0];
  const [selectedImage, setSelectedImage] = useState(activeImage);

  // Sync selected image with SKU image changes
  useEffect(() => {
    setSelectedImage(activeImage);
  }, [activeImage]);

  return (
    <div className="md:col-span-5 flex flex-col gap-4">
      {/* Main Image */}
      <div className="relative w-full pt-[100%] border border-neutral-100 rounded-sm overflow-hidden group hover:shadow-sm">
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={selectedImage}
          alt={name}
          className="absolute top-0 left-0 w-full h-full object-contain"
        />
      </div>

      {/* Thumbnails */}
      <div className="flex gap-2 overflow-x-auto pb-2">
        {displayImages.map((img, idx) => (
          <button
            key={idx}
            onMouseEnter={() => setSelectedImage(img)}
            onClick={() => setSelectedImage(img)}
            className={`shrink-0 w-20 h-20 border-2 rounded-sm overflow-hidden transition-colors ${
              selectedImage === img
                ? "border-[#ee4d2d]"
                : "border-transparent hover:border-[#ee4d2d]"
            }`}
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={img}
              alt={`${name} thumbnail ${idx + 1}`}
              className="w-full h-full object-cover"
            />
          </button>
        ))}
      </div>

      {/* Share & Like Actions */}
      <div className="flex justify-center gap-6 text-sm text-neutral-600 mt-2">
        <button className="flex items-center gap-1 hover:text-[#ee4d2d]">
          <span>Chia sẻ:</span>
          <div className="w-5 h-5 bg-blue-600 rounded-full"></div>
          <div className="w-5 h-5 bg-blue-400 rounded-full"></div>
          <div className="w-5 h-5 bg-red-500 rounded-full"></div>
        </button>
        <div className="w-px bg-neutral-300 h-5"></div>
        <button className="flex items-center gap-1 hover:text-[#ee4d2d]">
          <span>❤️ Đã thích (2,1k)</span>
        </button>
      </div>
    </div>
  );
}
