"use client";

interface PriceSectionProps {
  basePrice: number;
  currentPrice?: number;
  isLoading?: boolean;
}

export default function PriceSection({
  basePrice,
  currentPrice,
  isLoading,
}: PriceSectionProps) {
  const displayPrice = currentPrice || basePrice;
  const originalPrice = displayPrice * 1.3; // Mock discount calculation
  const discount = Math.round(
    ((originalPrice - displayPrice) / originalPrice) * 100
  );

  if (isLoading) {
    return (
      <div className="bg-[#fafafa] p-4 rounded-sm animate-pulse">
        <div className="h-8 bg-neutral-200 rounded w-48"></div>
      </div>
    );
  }

  return (
    <div className="bg-[#fafafa] p-4 flex items-center gap-4 rounded-sm">
      <div className="text-neutral-400 line-through text-base">
        ₫{new Intl.NumberFormat("vi-VN").format(originalPrice)}
      </div>
      <div className="text-3xl font-medium text-[#ee4d2d]">
        ₫{new Intl.NumberFormat("vi-VN").format(displayPrice)}
      </div>
      {discount > 0 && (
        <div className="bg-[#ee4d2d] text-white text-xs font-bold px-1 py-0.5 rounded-sm uppercase">
          Giảm {discount}%
        </div>
      )}
    </div>
  );
}
