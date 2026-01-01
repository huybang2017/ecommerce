// Centralized export for all React Query hooks
// Import này thay thế cho các API client cũ

// Auth hooks
export {
  useLogin,
  useRegister,
  useLogout,
  useProfile,
  authApi,
} from "./useAuth";

export type {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  ErrorResponse,
} from "./useAuth";

// Product hooks
export {
  useProducts,
  useProduct,
  useSearchProducts,
  useCategories,
  useCategory,
  useCategoryProducts,
  productApi,
} from "./useProducts";

export type {
  Product,
  Category,
  ProductsResponse,
  SearchParams,
} from "./useProducts";

// Cart hooks
export {
  useCart,
  useAddToCart,
  useUpdateCartItem,
  useRemoveFromCart,
  useClearCart,
  cartApi,
} from "./useCart";

export type {
  Cart,
  CartItem,
  AddToCartRequest,
  UpdateCartRequest,
} from "./useCart";

// Order hooks
export {
  useOrders,
  useOrder,
  useCreateOrder,
  useCancelOrder,
  orderApi,
} from "./useOrders";

export type {
  Order,
  OrderItem,
  CreateOrderRequest,
  OrdersResponse,
} from "./useOrders";
