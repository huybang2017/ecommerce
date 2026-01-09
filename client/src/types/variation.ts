// Variation types - theo backend schema

export interface VariationOption {
  id: number;
  variation_id: number;
  value: string; // "M", "L", "Red", "Blue"
}

export interface Variation {
  id: number;
  name: string; // "Màu Sắc", "Kích Thước"
  options: VariationOption[];
}

export interface VariationsResponse {
  variations: Variation[];
}

// UI State types
export interface SelectedVariations {
  [variationId: number]: number; // variationId -> optionId
}
