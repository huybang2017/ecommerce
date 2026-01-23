"use client";

import { useCartContext as useCart } from "@/contexts/CartContext";
import type { Cart, CartItem } from "@/hooks/useCart";
import { useAuth } from "@/contexts/AuthContext";
import { createOrder } from "@/services/order.service";
import { useAddresses, useSetDefaultAddress } from "@/hooks/useAddresses";
import { CreateOrderRequest } from "@/lib/types";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";

export default function CheckoutPage() {
  const { cart, loading: cartLoading, refreshCart } = useCart();
  const { user } = useAuth();
  const router = useRouter();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    shipping_name: "",
    shipping_phone: "",
    shipping_address: "",
    shipping_city: "",
    shipping_province: "",
    shipping_postal_code: "",
    shipping_country: "VN",
  });

  // Address modal state - use server-backed addresses (React Query)
  const [showAddressModal, setShowAddressModal] = useState(false);
  const { data: addressesData } = useAddresses();
  const setDefaultAddress = useSetDefaultAddress();
  const [selectedAddressId, setSelectedAddressId] = useState<number | null>(
    null,
  );

  useEffect(() => {
    if (!addressesData || addressesData.length === 0) return;

    // Prefer address marked as default, otherwise first address
    const defaultAddr =
      addressesData.find((a) => a.is_default) || addressesData[0];

    if (!selectedAddressId) {
      setSelectedAddressId(defaultAddr.id);
    }

    const sel =
      (addressesData || []).find((a) => a.id === selectedAddressId) ||
      defaultAddr;

    if (sel) {
      setFormData((prev) => ({
        ...prev,
        shipping_name: sel.recipient_name,
        shipping_phone: sel.phone_number,
        shipping_address: sel.address_line,
        shipping_city: sel.city || "",
        shipping_province: sel.district || "",
      }));
    }
  }, [selectedAddressId, addressesData]);

  useEffect(() => {
    if (
      !cartLoading &&
      (!cart ||
        !cart.items ||
        (Array.isArray(cart.items) && cart.items.length === 0))
    ) {
      router.push("/cart");
    }
  }, [cart, cartLoading, router]);

  const handleSubmit = async (e?: React.FormEvent) => {
    e?.preventDefault();
    setError(null);
    setLoading(true);

    try {
      if (
        !cart ||
        !cart.items ||
        (Array.isArray(cart.items) && cart.items.length === 0)
      ) {
        setError("Cart is empty");
        setLoading(false);
        return;
      }

      const sessionId =
        typeof window !== "undefined"
          ? localStorage.getItem("session_id") || ""
          : "";

      const itemsForRequest: CartItem[] = Array.isArray(cart.items)
        ? cart.items
        : (Object.values(cart.items) as CartItem[]);

      const orderRequest: CreateOrderRequest = {
        user_id: user?.id,
        session_id: sessionId,
        ...formData,
        shipping_fee: 0,
        tax: 0,
        discount: 0,
        // Attach items in legacy format if backend expects items
        items: itemsForRequest.map((it: CartItem) => ({
          product_id:
            it.product_item_id ??
            (it as unknown as { product_id?: number }).product_id,
          quantity: it.quantity,
        })),
      };

      const order = await createOrder(orderRequest);

      await refreshCart();
      router.push(`/orders/${order.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create order");
    } finally {
      setLoading(false);
    }
  };

  // form change handler removed (not used) — form fields are filled from selected address

  if (cartLoading) {
    return (
      <div className="min-h-screen bg-white py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">Loading...</div>
      </div>
    );
  }

  if (!cart || !cart.items) return null;

  const cartItems: CartItem[] = Array.isArray(cart.items)
    ? cart.items
    : (Object.values(cart.items) as CartItem[]);
  const subtotal = (cart as Cart).total_price || 0;
  const shippingFee = 0;
  const tax = 0;
  const discount = 0;
  const total = subtotal + shippingFee + tax - discount;

  return (
    <div className="min-h-screen bg-white py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <div className="mb-6 rounded-lg border border-neutral-200 bg-white p-4">
          <div className="flex items-start justify-between">
            <div className="flex items-start gap-4">
              <div className="text-[#ee4d2d] mt-1">📍</div>
              <div>
                <div className="text-sm font-semibold text-neutral-900">
                  Địa Chỉ Nhận Hàng
                </div>
                <div className="mt-1 text-sm text-neutral-700">
                  {formData.shipping_name} {formData.shipping_phone}
                </div>
                <div className="text-sm text-neutral-500">
                  {formData.shipping_address}
                </div>
              </div>
            </div>
            <div className="text-sm">
              <button
                onClick={() => setShowAddressModal(true)}
                className="text-[#1677ff]"
              >
                Thay Đổi
              </button>
            </div>
          </div>

          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
              <p className="text-sm text-red-800">{error}</p>
            </div>
          )}
        </div>

        <div className="grid gap-8">
          <div>
            <div className="rounded-lg border border-neutral-200 bg-white">
              <div className="flex items-center justify-between p-4 border-b border-neutral-100">
                <div className="flex items-center gap-3">
                  <input type="checkbox" className="h-4 w-4" />
                  <span className="text-sm font-medium">Sản phẩm</span>
                </div>
                <div className="hidden sm:flex gap-6 text-sm text-neutral-600">
                  <span className="w-28 text-right">Đơn giá</span>
                  <span className="w-24 text-center">Số lượng</span>
                  <span className="w-28 text-right">Thành tiền</span>
                  <span className="w-24 text-right">Thao tác</span>
                </div>
              </div>

              <div className="p-4">
                <div className="mb-4 flex items-center gap-3">
                  <span className="rounded bg-[#fff2ef] px-2 py-0.5 text-xs font-medium text-[#ee4d2d]">
                    Yêu thích
                  </span>
                  <span className="text-sm text-neutral-700">
                    Nhựa_Việt_Nhật
                  </span>
                  <button className="ml-4 text-sm text-[#00a86b]">
                    Chat ngay
                  </button>
                </div>

                <div className="space-y-3">
                  {cartItems.map((item: CartItem) => (
                    <div
                      key={
                        item.product_item_id ??
                        (item as unknown as { product_id?: number }).product_id
                      }
                      className="flex items-center gap-4 border-b border-neutral-100 py-3"
                    >
                      {item.image && (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                          src={item.image}
                          alt={item.name}
                          className="w-16 h-16 object-cover rounded"
                        />
                      )}
                      <div className="flex-1 text-sm">
                        <div className="font-medium text-neutral-900">
                          {item.name}
                        </div>
                        <div className="text-neutral-500 mt-1">
                          Phân loại:{" "}
                          {(item as unknown as { variant?: string }).variant ||
                            "-"}
                        </div>
                      </div>
                      <div className="w-28 text-right text-sm">
                        {item.price.toLocaleString("vi-VN")}đ
                      </div>
                      <div className="w-24 text-center text-sm">
                        {item.quantity}
                      </div>
                      <div className="w-28 text-right text-sm font-semibold text-[#ee4d2d]">
                        {(item.price * item.quantity).toLocaleString("vi-VN")}đ
                      </div>
                      <div className="w-24 text-right text-sm text-[#ee4d2d]">
                        <button>Xóa</button>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="mt-4 border-t border-neutral-100 pt-4 text-sm text-neutral-700">
                  <div className="flex items-center justify-between">
                    <div>
                      <div className="mb-2">
                        Voucher của Shop ·{" "}
                        <button className="text-[#1677ff]">Chọn Voucher</button>
                      </div>
                      <div className="text-neutral-500">
                        Lời nhắn:{" "}
                        <input
                          className="ml-2 rounded border px-2 py-1 text-sm"
                          placeholder="Lưu ý cho Người bán..."
                        />
                      </div>
                    </div>
                    <div className="text-right">
                      <div>Phương thức vận chuyển</div>
                      <div className="text-[#1677ff]">
                        Nhận từ 25 Th01 - 27 Th01 ·{" "}
                        <button className="text-[#1677ff]">Thay Đổi</button>
                      </div>
                      <div className="mt-2 text-sm text-neutral-900 font-semibold">
                        Tổng số tiền ({cartItems.length} sản phẩm):{" "}
                        <span className="text-[#ee4d2d] ml-2">
                          {total.toLocaleString("vi-VN")}đ
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-6 space-y-4">
          <div className="rounded-lg border border-neutral-200 bg-white p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="text-[#ee4d2d]">🎟️</span>
                <span className="text-sm font-medium">Shopee Voucher</span>
              </div>
              <button className="text-[#1677ff] text-sm">Chọn Voucher</button>
            </div>
          </div>

          <div className="rounded-lg border border-neutral-200 bg-white p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="text-[#f5a623]">◕</span>
                <span className="text-sm font-medium">Shopee Xu</span>
                <span className="text-sm text-neutral-500">
                  Không thể sử dụng Xu
                </span>
              </div>
              <div className="text-sm text-neutral-500">[-0₫]</div>
            </div>
          </div>
        </div>

        <div className="mt-6 rounded-lg border border-neutral-200 bg-white p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold">Phương thức thanh toán</h3>
            <button className="text-[#1677ff] text-sm">Thay Đổi</button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="text-sm text-neutral-600">
              <div className="flex justify-between py-2">
                <span>Tổng tiền hàng</span>
                <span>{subtotal.toLocaleString("vi-VN")}đ</span>
              </div>
              <div className="flex justify-between py-2">
                <span>Tổng tiền phí vận chuyển</span>
                <span>{shippingFee.toLocaleString("vi-VN")}đ</span>
              </div>
              {discount > 0 && (
                <div className="flex justify-between py-2">
                  <span>Giảm</span>
                  <span className="text-green-600">
                    -{discount.toLocaleString("vi-VN")}đ
                  </span>
                </div>
              )}
            </div>

            <div className="flex flex-col justify-end">
              <div className="flex justify-between items-center">
                <div className="text-sm text-neutral-600">Tổng thanh toán</div>
                <div className="text-2xl font-bold text-[#ee4d2d]">
                  {total.toLocaleString("vi-VN")}đ
                </div>
              </div>
              <div className="mt-4 flex gap-3">
                <button
                  onClick={() => handleSubmit()}
                  disabled={loading}
                  className="ml-auto rounded-sm bg-[#ee4d2d] px-6 py-3 text-sm font-medium text-white"
                >
                  Đặt hàng
                </button>
              </div>
            </div>
          </div>

          <div className="mt-3 text-xs text-neutral-500">
            Nhấn &quot;Đặt hàng&quot; đồng nghĩa với việc bạn đồng ý tuân theo
            Điều khoản Shopee
          </div>
        </div>
      </div>

      {/* Address modal */}
      {showAddressModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="mx-auto w-full max-w-3xl rounded bg-white p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">Chọn địa chỉ</h3>
              <button
                onClick={() => setShowAddressModal(false)}
                className="text-neutral-500"
              >
                Đóng
              </button>
            </div>

            <div className="space-y-3">
              {addressesData?.map((a) => (
                <div
                  key={a.id}
                  className={`flex items-start justify-between rounded border p-3 ${selectedAddressId === a.id ? "border-[#ee4d2d]" : "border-neutral-200"}`}
                >
                  <div>
                    <div className="font-medium">
                      {a.recipient_name}{" "}
                      <span className="text-sm text-neutral-500">
                        {a.phone_number}
                      </span>
                    </div>
                    <div className="text-sm text-neutral-600">
                      {a.address_line}
                    </div>
                    {a.label && (
                      <div className="text-xs text-neutral-500">{a.label}</div>
                    )}
                  </div>
                  <div className="flex flex-col items-end gap-2">
                    <div className="flex flex-col items-end gap-2">
                      <button
                        onClick={() => setSelectedAddressId(a.id)}
                        className="text-[#1677ff]"
                      >
                        Chọn
                      </button>
                      <button
                        onClick={() => setDefaultAddress.mutate(a.id)}
                        className="text-sm text-neutral-500"
                      >
                        Đặt làm mặc định
                      </button>
                    </div>
                    <div className="text-xs text-neutral-500">
                      {a.is_default ? "Mặc định" : ""}
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <div className="mt-4 flex justify-end gap-3">
              <button
                onClick={() => setShowAddressModal(false)}
                className="rounded border px-4 py-2"
              >
                Hủy
              </button>
              <button
                onClick={() => setShowAddressModal(false)}
                className="rounded bg-[#ee4d2d] px-4 py-2 text-white"
              >
                Hoàn tất
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
