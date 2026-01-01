// TypeScript types for E-commerce Marketplace
// NOTE: Re-export from hooks for convenience

// Re-export from hooks
export type {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
} from "@/hooks/useAuth";

export type {
  Product,
  Category,
  ProductsResponse,
  SearchParams,
} from "@/hooks/useProducts";

export type {
  Cart,
  CartItem,
  AddToCartRequest,
  UpdateCartRequest,
} from "@/hooks/useCart";

export type {
  Order,
  OrderItem,
  CreateOrderRequest,
  OrdersResponse,
} from "@/hooks/useOrders";

// ==========================================
// SHOP TYPES (từ db-diagram.db)
// ==========================================

export interface Shop {
  id: number;
  owner_user_id: number;
  name: string;
  description: string;
  logo_url: string;
  cover_url: string;
  is_official: boolean;
  rating: number;
  response_rate: number;
  status: string; // ACTIVE, SUSPENDED
  created_at: string;
  updated_at: string;
}

// ==========================================
// SKU TYPES (từ db-diagram.db)
// ==========================================

export interface Variation {
  id: number;
  product_id: number;
  name: string; // "Size", "Color"
}

export interface VariationOption {
  id: number;
  variation_id: number;
  value: string; // "M", "L", "Red", "Blue"
}

export interface ProductItem {
  id: number;
  product_id: number;
  sku_code: string;
  image_url: string;
  price: number;
  qty_in_stock: number;
  status: string; // ACTIVE, OUT_OF_STOCK, DISABLED
}

export interface SKUConfiguration {
  product_item_id: number;
  variation_option_id: number;
}

// ==========================================
// ATTRIBUTE TYPES (EAV - từ db-diagram.db)
// ==========================================

export interface CategoryAttribute {
  id: number;
  category_id: number;
  attribute_name: string;
  input_type: string; // text, number, select, checkbox
  is_mandatory: boolean;
  is_filterable: boolean;
}

export interface ProductAttributeValue {
  id: number;
  product_id: number;
  attribute_id: number;
  value: string;
}

// ==========================================
// COMMON API TYPES
// ==========================================

export interface ApiError {
  error: string;
  message?: string;
}

export type OrderStatus =
  | "pending"
  | "paid"
  | "processing"
  | "shipped"
  | "delivered"
  | "cancelled";
