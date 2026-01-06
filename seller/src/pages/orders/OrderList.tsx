import { Search, Printer, MessageSquare } from "lucide-react";

export default function OrderList() {
  const orders = [
    {
      id: "230815G7X8Y9Z",
      user: "nguyenvana",
      products: [
        {
          name: "Áo thun nam Cotton Compact CAO CẤP",
          variation: "Trắng, L",
          quantity: 2,
          price: 159000,
          image: "https://placehold.co/50x50",
        },
      ],
      total: 340000,
      status: "pending",
      payment: "COD",
      carrier: "J&T Express",
    },
    {
      id: "230815H1A2B3C",
      user: "tranthib",
      products: [
        {
          name: "Quần Jeans Nam Slimfit Co Giãn Nhẹ",
          variation: "Xanh, 32",
          quantity: 1,
          price: 350000,
          image: "https://placehold.co/50x50",
        },
      ],
      total: 375000,
      status: "shipping",
      payment: "ShopeePay",
      carrier: "Giao Hàng Nhanh",
    },
  ];

  return (
    <div className="bg-white rounded-sm shadow-sm">
      {/* Header */}
      <div className="p-4 border-b border-neutral-100 flex justify-between items-center">
        <div>
          <h1 className="text-xl font-bold text-neutral-800">
            Quản Lý Đơn Hàng
          </h1>
          <p className="text-sm text-neutral-500 mt-1">
            Xem và xử lý đơn hàng của bạn
          </p>
        </div>
        <button className="bg-[#ee4d2d] text-white px-4 py-2 rounded-sm text-sm font-medium hover:bg-[#d73211]">
          Giao hàng loạt
        </button>
      </div>

      {/* Tabs */}
      <div className="flex text-sm font-medium text-neutral-500 border-b border-neutral-100 overflow-x-auto">
        <div className="whitespace-nowrap px-4 py-3 text-[#ee4d2d] border-b-2 border-[#ee4d2d] cursor-pointer">
          Tất cả
        </div>
        <div className="whitespace-nowrap px-4 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Chờ xác nhận
        </div>
        <div className="whitespace-nowrap px-4 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Chờ lấy hàng
        </div>
        <div className="whitespace-nowrap px-4 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Đang giao
        </div>
        <div className="whitespace-nowrap px-4 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Đã giao
        </div>
        <div className="whitespace-nowrap px-4 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Đơn hủy
        </div>
        <div className="whitespace-nowrap px-4 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Trả hàng/Hoàn tiền
        </div>
      </div>

      {/* Search Bar */}
      <div className="p-4 bg-neutral-50 flex gap-2">
        <div className="relative flex-1">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search size={16} className="text-neutral-400" />
          </div>
          <input
            type="text"
            className="block w-full pl-10 pr-3 py-2 border border-neutral-300 rounded-sm leading-5 bg-white placeholder-neutral-500 focus:outline-none focus:placeholder-neutral-400 focus:border-[#ee4d2d] sm:text-sm"
            placeholder="Tìm kiếm Mã đơn hàng, Tên người mua hoặc Sản phẩm"
          />
        </div>
        <button className="bg-[#ee4d2d] text-white px-6 py-2 rounded-sm font-medium hover:bg-[#d73211]">
          Tìm kiếm
        </button>
        <button className="bg-white border border-neutral-300 text-neutral-600 px-6 py-2 rounded-sm font-medium hover:bg-neutral-50">
          Xuất
        </button>
      </div>

      {/* Order List */}
      <div className="p-4 space-y-4">
        {orders.map((order) => (
          <div key={order.id} className="border border-neutral-200 rounded-sm">
            {/* Order Header */}
            <div className="bg-neutral-50 px-4 py-3 border-b border-neutral-200 flex justify-between items-center text-sm">
              <div className="flex items-center gap-4">
                <div className="font-bold text-neutral-700">{order.user}</div>
                <button className="text-[#ee4d2d] border border-[#ee4d2d] bg-white px-2 py-0.5 rounded-sm text-xs flex items-center gap-1">
                  <MessageSquare size={10} /> Chat ngay
                </button>
              </div>
              <div className="flex items-center gap-4">
                <span className="text-neutral-500">
                  Mã đơn hàng:{" "}
                  <span className="text-neutral-800 font-medium">
                    {order.id}
                  </span>
                </span>
              </div>
            </div>

            {/* Order Body */}
            <div className="divide-y divide-neutral-100">
              {order.products.map((product, idx) => (
                <div key={idx} className="p-4 flex gap-4">
                  <img
                    src={product.image}
                    alt=""
                    className="w-20 h-20 object-cover border border-neutral-200 rounded-sm"
                  />
                  <div className="flex-1">
                    <div className="font-medium text-neutral-800">
                      {product.name}
                    </div>
                    <div className="text-neutral-500 text-sm mt-1">
                      Phân loại hàng: {product.variation}
                    </div>
                    <div className="text-neutral-500 text-sm mt-1">
                      x{product.quantity}
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="font-medium text-neutral-800">
                      ₫{product.price.toLocaleString("vi-VN")}
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Order Footer */}
            <div className="px-4 py-3 border-t border-neutral-200 bg-[#fdfdfd] flex justify-between items-center">
              <div className="text-sm text-neutral-600">
                Đơn vị vận chuyển:{" "}
                <span className="font-medium text-neutral-800">
                  {order.carrier}
                </span>
              </div>
              <div className="text-right">
                <div className="flex items-center gap-2 justify-end mb-2">
                  <span className="text-sm text-neutral-600">
                    Tổng đơn hàng:
                  </span>
                  <span className="text-xl font-bold text-[#ee4d2d]">
                    ₫{order.total.toLocaleString("vi-VN")}
                  </span>
                </div>
                <div className="flex gap-2 justify-end">
                  <button className="bg-[#ee4d2d] text-white px-4 py-1.5 rounded-sm text-sm hover:bg-[#d73211]">
                    Chuẩn bị hàng
                  </button>
                  <button className="border border-neutral-300 text-neutral-600 px-4 py-1.5 rounded-sm text-sm hover:bg-neutral-50 flex items-center gap-2">
                    <Printer size={14} /> In phiếu
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
