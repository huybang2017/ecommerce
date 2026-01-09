/**
 * ProductItem (SKU) types - represents actual sellable variations
 * Backend: product_item table
 */

export interface ProductItem {
  id: number;
  product_id: number;
  sku_code: string;
  image_url?: string;
  price: number;
  qty_in_stock: number;
  status: "ACTIVE" | "OUT_OF_STOCK" | "DISABLED";
  variation_option_ids?: number[]; // [1, 5] = Size M + Color Red (for matching)
}

export interface ProductItemsResponse {
  items: ProductItem[];
  count: number;
}

export interface Variation {
  id: number;
  product_id: number;
  name: string; // "Size", "Color", "Storage"
}

export interface VariationOption {
  id: number;
  variation_id: number;
  value: string; // "M", "L", "Red", "Blue"
}

export interface SKUConfiguration {
  product_item_id: number;
  variation_option_id: number;
}

/**
 * Extended ProductItem with variation details for UI
 */
export interface ProductItemWithVariations extends ProductItem {
  variations?: {
    variation_name: string;
    option_value: string;
  }[];
}
