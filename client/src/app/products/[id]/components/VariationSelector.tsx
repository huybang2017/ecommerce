"use client";

import { useState, useEffect, useMemo } from "react";
import { useProductVariations } from "@/hooks/useVariations";
import { useProductItems } from "@/hooks/useProductItems";
import type { ProductItem } from "@/types/product-item";
import type { SelectedVariations } from "@/types/variation";

interface VariationSelectorProps {
  productId: number;
  selectedSKU: ProductItem | null;
  onSKUChange: (sku: ProductItem | null) => void;
  isLoading?: boolean;
}

export default function VariationSelector({
  productId,
  selectedSKU,
  onSKUChange,
  isLoading = false,
}: VariationSelectorProps) {
  const { data: variations, isLoading: isLoadingVariations } =
    useProductVariations(productId);
  const { data: itemsData, isLoading: isLoadingItems } =
    useProductItems(productId);

  const [selected, setSelected] = useState<SelectedVariations>({});

  const items = useMemo(() => itemsData?.items || [], [itemsData]);

  // Find matching SKU based on selected variations
  const findMatchingSKU = (
    selectedVars: SelectedVariations
  ): ProductItem | null => {
    if (!variations || variations.length === 0) {
      // No variations - return first available item
      return items.find((item) => item.status === "ACTIVE") || null;
    }

    // Check if all variations are selected
    const allSelected = variations.every(
      (v) => selectedVars[v.id] !== undefined
    );
    if (!allSelected) return null;

    // Get selected option IDs
    const selectedOptionIds = Object.values(selectedVars).sort((a, b) => a - b);

    // Find SKU that matches ALL selected options
    const matchingSKU = items.find((item) => {
      if (!item.variation_option_ids || item.variation_option_ids.length === 0)
        return false;

      const itemOptions = [...item.variation_option_ids].sort((a, b) => a - b);

      // Must match exactly
      if (itemOptions.length !== selectedOptionIds.length) return false;

      return itemOptions.every((id, index) => id === selectedOptionIds[index]);
    });

    return matchingSKU || null;
  };

  // Handle option click
  const handleOptionClick = (variationId: number, optionId: number) => {
    const newSelected = {
      ...selected,
      [variationId]: optionId,
    };
    setSelected(newSelected);

    // Find and set matching SKU
    const matchedSKU = findMatchingSKU(newSelected);
    onSKUChange(matchedSKU);
  };

  // Compute auto-selected state
  const autoSelectedState = useMemo(() => {
    if (
      !variations ||
      !items ||
      variations.length === 0 ||
      items.length === 0 ||
      selectedSKU
    ) {
      return null;
    }

    // Try to find first available item and extract its variation options
    const firstAvailable = items.find(
      (item) => item.status === "ACTIVE" && item.qty_in_stock > 0
    );

    if (firstAvailable && firstAvailable.variation_option_ids) {
      // Build selected map from first SKU's options
      const autoSelected: SelectedVariations = {};

      firstAvailable.variation_option_ids.forEach((optionId) => {
        // Find which variation this option belongs to
        for (const variation of variations) {
          const option = variation.options.find((opt) => opt.id === optionId);
          if (option) {
            autoSelected[variation.id] = optionId;
            break;
          }
        }
      });

      return { selected: autoSelected, sku: firstAvailable };
    }

    return null;
  }, [variations, items, selectedSKU]);

  // Auto-select first available combination on load
  useEffect(() => {
    if (autoSelectedState && !selectedSKU) {
      setSelected(autoSelectedState.selected);
      onSKUChange(autoSelectedState.sku);
    }
  }, [autoSelectedState, selectedSKU, onSKUChange]);

  // Check if option is available (has stock)
  const isOptionAvailable = (
    variationId: number,
    optionId: number
  ): boolean => {
    const tempSelected = { ...selected, [variationId]: optionId };

    // Check if any SKU exists with this option combination
    return items.some((item) => {
      if (
        !item.variation_option_ids ||
        item.status !== "ACTIVE" ||
        item.qty_in_stock === 0
      )
        return false;

      // Must include this option
      if (!item.variation_option_ids.includes(optionId)) return false;

      // Check if matches other selected variations
      for (const [varId, optId] of Object.entries(tempSelected)) {
        if (Number(varId) !== variationId) {
          if (!item.variation_option_ids.includes(optId)) return false;
        }
      }

      return true;
    });
  };

  if (isLoadingVariations || isLoadingItems || isLoading) {
    return (
      <div className="grid grid-cols-[110px_1fr] gap-4">
        <div className="text-neutral-500 text-sm">Loading...</div>
        <div className="flex gap-2">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="w-20 h-9 bg-neutral-200 animate-pulse rounded-sm"
            />
          ))}
        </div>
      </div>
    );
  }

  if (!variations || variations.length === 0) {
    return null; // No variations
  }

  return (
    <div className="space-y-6">
      {variations.map((variation) => (
        <div key={variation.id} className="grid grid-cols-[110px_1fr] gap-4">
          <div className="text-neutral-500 text-sm">{variation.name}</div>
          <div className="flex flex-wrap gap-2">
            {variation.options.map((option) => {
              const isSelected = selected[variation.id] === option.id;
              const isAvailable = isOptionAvailable(variation.id, option.id);

              return (
                <button
                  key={option.id}
                  onClick={() => {
                    if (isAvailable) {
                      handleOptionClick(variation.id, option.id);
                    }
                  }}
                  disabled={!isAvailable}
                  className={`
                    px-5 py-2 border rounded-sm min-w-20 text-sm transition-all
                    ${
                      isSelected
                        ? "border-[#ee4d2d] text-[#ee4d2d] bg-[#fff6f5]"
                        : isAvailable
                        ? "border-neutral-300 hover:border-[#ee4d2d] cursor-pointer"
                        : "border-neutral-200 bg-neutral-50 text-neutral-400 cursor-not-allowed opacity-60"
                    }
                  `}
                >
                  {option.value}
                </button>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}
