"use client";

interface QuantitySelectorProps {
  value: number;
  onChange: (value: number) => void;
  max?: number;
  disabled?: boolean;
}

export default function QuantitySelector({
  value,
  onChange,
  max,
  disabled = false,
}: QuantitySelectorProps) {
  const handleDecrease = () => {
    if (value > 1) {
      onChange(value - 1);
    }
  };

  const handleIncrease = () => {
    if (!max || value < max) {
      onChange(value + 1);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = parseInt(e.target.value) || 1;
    if (newValue >= 1 && (!max || newValue <= max)) {
      onChange(newValue);
    }
  };

  return (
    <div className="flex items-center gap-4">
      <div className="flex items-center border border-neutral-300 rounded-sm">
        <button
          onClick={handleDecrease}
          disabled={disabled || value <= 1}
          className="w-8 h-8 flex items-center justify-center border-r hover:bg-neutral-50 text-neutral-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          -
        </button>
        <input
          type="number"
          value={value}
          onChange={handleInputChange}
          disabled={disabled}
          min={1}
          max={max}
          className="w-12 h-8 text-center text-neutral-900 border-none outline-none disabled:bg-neutral-50"
        />
        <button
          onClick={handleIncrease}
          disabled={disabled || (max !== undefined && value >= max)}
          className="w-8 h-8 flex items-center justify-center border-l hover:bg-neutral-50 text-neutral-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          +
        </button>
      </div>

      {max !== undefined && (
        <div className="text-neutral-500 text-sm">
          {max > 0 ? (
            <span>Còn {max} sản phẩm</span>
          ) : (
            <span className="text-red-500">Hết hàng</span>
          )}
        </div>
      )}
    </div>
  );
}
