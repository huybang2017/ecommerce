import {
  ChevronRight,
  ClipboardList,
  Wallet,
  ShoppingBag,
  Star,
} from "lucide-react";

export default function Dashboard() {
  return (
    <div className="space-y-6">
      {/* Todo List Section */}
      <div className="bg-white rounded-sm shadow-sm p-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="font-bold text-lg text-neutral-800">
            Danh sách cần làm
          </h2>
          <button className="text-sm text-neutral-500 flex items-center hover:text-[#ee4d2d]">
            Những việc cần làm <ChevronRight size={14} />
          </button>
        </div>

        <div className="grid grid-cols-4 gap-4 text-center">
          <div className="flex flex-col items-center gap-1 cursor-pointer hover:bg-neutral-50 py-4 rounded-sm">
            <div className="text-xl font-bold text-blue-600">0</div>
            <div className="text-sm text-neutral-500">Chờ xác nhận</div>
          </div>
          <div className="flex flex-col items-center gap-1 cursor-pointer hover:bg-neutral-50 py-4 rounded-sm">
            <div className="text-xl font-bold text-blue-600">0</div>
            <div className="text-sm text-neutral-500">Chờ lấy hàng</div>
          </div>
          <div className="flex flex-col items-center gap-1 cursor-pointer hover:bg-neutral-50 py-4 rounded-sm">
            <div className="text-xl font-bold text-blue-600">0</div>
            <div className="text-sm text-neutral-500">Đã xử lý</div>
          </div>
          <div className="flex flex-col items-center gap-1 cursor-pointer hover:bg-neutral-50 py-4 rounded-sm">
            <div className="text-xl font-bold text-blue-600">0</div>
            <div className="text-sm text-neutral-500">Đơn hủy</div>
          </div>
        </div>
      </div>

      {/* Business Insights */}
      <div className="bg-white rounded-sm shadow-sm p-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="font-bold text-lg text-neutral-800">
            Phân Tích Bán Hàng
          </h2>
          <button className="text-sm text-neutral-500 flex items-center hover:text-[#ee4d2d]">
            Xem thêm <ChevronRight size={14} />
          </button>
        </div>
        <div className="text-xs text-neutral-400 mb-4">
          Hôm nay 00:00 GMT+7 - 18:00 GMT+7
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div className="border-r border-neutral-100 last:border-0 pl-4 py-2">
            <div className="text-sm text-neutral-500 mb-1 flex items-center gap-1">
              Doanh số đơn hàng
            </div>
            <div className="text-xl font-bold">₫0</div>
            <div className="text-xs text-neutral-400 mt-1">
              Vs hôm qua 0.00%
            </div>
          </div>
          <div className="border-r border-neutral-100 last:border-0 pl-4 py-2">
            <div className="text-sm text-neutral-500 mb-1 flex items-center gap-1">
              Lượt truy cập
            </div>
            <div className="text-xl font-bold">0</div>
            <div className="text-xs text-neutral-400 mt-1">
              Vs hôm qua 0.00%
            </div>
          </div>
          <div className="border-r border-neutral-100 last:border-0 pl-4 py-2">
            <div className="text-sm text-neutral-500 mb-1 flex items-center gap-1">
              Đơn hàng
            </div>
            <div className="text-xl font-bold">0</div>
            <div className="text-xs text-neutral-400 mt-1">
              Vs hôm qua 0.00%
            </div>
          </div>
          <div className="border-r border-neutral-100 last:border-0 pl-4 py-2">
            <div className="text-sm text-neutral-500 mb-1 flex items-center gap-1">
              Tỷ lệ chuyển đổi
            </div>
            <div className="text-xl font-bold">0.00%</div>
            <div className="text-xs text-neutral-400 mt-1">
              Vs hôm qua 0.00%
            </div>
          </div>
        </div>
      </div>

      {/* Marketing Section */}
      <div className="bg-white rounded-sm shadow-sm p-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="font-bold text-lg text-neutral-800">Kênh Marketing</h2>
          <button className="text-sm text-blue-500">Xem tất cả &gt;</button>
        </div>
        <div className="grid grid-cols-3 gap-4">
          <div className="border border-neutral-200 p-4 rounded-sm flex items-center gap-3 cursor-pointer hover:border-[#ee4d2d] hover:bg-orange-50/20">
            <div className="w-10 h-10 bg-orange-100 text-[#ee4d2d] rounded-full flex items-center justify-center">
              <ShoppingBag size={20} />
            </div>
            <div>
              <div className="font-medium text-sm">Mã Giảm Giá Của Tôi</div>
              <div className="text-xs text-neutral-500">
                Công cụ tăng đơn hàng
              </div>
            </div>
          </div>
          <div className="border border-neutral-200 p-4 rounded-sm flex items-center gap-3 cursor-pointer hover:border-[#ee4d2d] hover:bg-orange-50/20">
            <div className="w-10 h-10 bg-orange-100 text-[#ee4d2d] rounded-full flex items-center justify-center">
              <Star size={20} />
            </div>
            <div>
              <div className="font-medium text-sm">Chương Trình Của Tôi</div>
              <div className="text-xs text-neutral-500">Công cụ marketing</div>
            </div>
          </div>
          <div className="border border-neutral-200 p-4 rounded-sm flex items-center gap-3 cursor-pointer hover:border-[#ee4d2d] hover:bg-orange-50/20">
            <div className="w-10 h-10 bg-orange-100 text-[#ee4d2d] rounded-full flex items-center justify-center">
              <ClipboardList size={20} />
            </div>
            <div>
              <div className="font-medium text-sm">Flash Sale Của Shop</div>
              <div className="text-xs text-neutral-500">
                Công cụ tạo đơn nhanh
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
