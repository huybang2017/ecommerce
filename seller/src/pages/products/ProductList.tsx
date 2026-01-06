import {
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Edit,
  Trash2,
  Eye,
} from "lucide-react";

export default function ProductList() {
  const products = [
    {
      id: 1,
      image: "https://placehold.co/50x50",
      name: "Áo thun nam Cotton Compact CAO CẤP, phông trơn cổ tròn ngắn tay, form Regular",
      sku: "SKU-001",
      price: 159000,
      stock: 120,
      sold: 45,
      status: "active",
    },
    {
      id: 2,
      image: "https://placehold.co/50x50",
      name: "Quần Jeans Nam Slimfit Co Giãn Nhẹ, Màu Xanh Basic",
      sku: "SKU-002",
      price: 350000,
      stock: 50,
      sold: 12,
      status: "active",
    },
    {
      id: 3,
      image: "https://placehold.co/50x50",
      name: "Giày Sneaker Nam Thể Thao Cổ Thấp",
      sku: "SKU-003",
      price: 550000,
      stock: 0,
      sold: 120,
      status: "soldout",
    },
    {
      id: 4,
      image: "https://placehold.co/50x50",
      name: "Balo Laptop Chống Nước 15.6 inch",
      sku: "SKU-004",
      price: 299000,
      stock: 20,
      sold: 5,
      status: "disabled",
    },
  ];

  return (
    <div className="bg-white rounded-sm shadow-sm">
      {/* Header */}
      <div className="p-4 border-b border-neutral-100 flex justify-between items-center">
        <div>
          <h1 className="text-xl font-bold text-neutral-800">Sản Phẩm</h1>
          <p className="text-sm text-neutral-500 mt-1">
            Quản lý các sản phẩm của shop
          </p>
        </div>
        <button className="bg-[#ee4d2d] text-white px-4 py-2 rounded-sm text-sm font-medium flex items-center gap-2 hover:bg-[#d73211]">
          <Plus size={16} /> Thêm 1 sản phẩm mới
        </button>
      </div>

      {/* Tabs */}
      <div className="flex text-sm font-medium text-neutral-500 border-b border-neutral-100">
        <div className="px-6 py-3 text-[#ee4d2d] border-b-2 border-[#ee4d2d] cursor-pointer">
          Tất cả
        </div>
        <div className="px-6 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Đang hoạt động
        </div>
        <div className="px-6 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Hết hàng
        </div>
        <div className="px-6 py-3 hover:text-[#ee4d2d] cursor-pointer">
          Đã ẩn
        </div>
      </div>

      {/* Filters */}
      <div className="p-4 flex gap-4 bg-neutral-50 text-sm">
        <div className="relative">
          <select className="border border-neutral-300 rounded-sm px-3 py-2 pr-8 w-40 focus:outline-none focus:border-[#ee4d2d]">
            <option>Tên sản phẩm</option>
            <option>SKU</option>
            <option>Mã sản phẩm</option>
          </select>
        </div>
        <div className="relative flex-1">
          <input
            type="text"
            placeholder="Tìm sản phẩm..."
            className="w-full border border-neutral-300 rounded-sm px-3 py-2 pl-9 focus:outline-none focus:border-[#ee4d2d]"
          />
          <Search
            size={16}
            className="absolute left-3 top-2.5 text-neutral-400"
          />
        </div>
        <button className="bg-[#ee4d2d] text-white px-6 py-2 rounded-sm hover:bg-[#d73211]">
          Tìm
        </button>
        <button className="bg-white border border-neutral-300 px-6 py-2 rounded-sm hover:bg-neutral-50 text-neutral-600">
          Đặt lại
        </button>
      </div>

      {/* Table */}
      <div className="p-4">
        <div className="border border-neutral-200 rounded-sm overflow-hidden">
          <table className="w-full text-sm text-left">
            <thead className="bg-[#f6f6f6] text-neutral-600 font-medium">
              <tr>
                <th className="py-3 px-4 w-10">
                  <input type="checkbox" />
                </th>
                <th className="py-3 px-4">Tên sản phẩm</th>
                <th className="py-3 px-4 w-32">SKU phân loại</th>
                <th className="py-3 px-4 w-32 text-right">Giá</th>
                <th className="py-3 px-4 w-24 text-center">Kho hàng</th>
                <th className="py-3 px-4 w-24 text-center">Đã bán</th>
                <th className="py-3 px-4 w-24 text-center">Trạng thái</th>
                <th className="py-3 px-4 w-24 text-center">Thao tác</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-neutral-100">
              {products.map((product) => (
                <tr key={product.id} className="hover:bg-neutral-50">
                  <td className="py-3 px-4">
                    <input type="checkbox" />
                  </td>
                  <td className="py-3 px-4">
                    <div className="flex gap-3">
                      <img
                        src={product.image}
                        alt=""
                        className="w-12 h-12 object-cover rounded-sm border border-neutral-200"
                      />
                      <div className="flex-1">
                        <div className="line-clamp-2 mb-1">{product.name}</div>
                        <div className="flex gap-2 text-xs">
                          <button className="text-blue-500 hover:underline flex items-center gap-1">
                            <Edit size={12} /> Sửa
                          </button>
                          <button className="text-neutral-500 hover:underline flex items-center gap-1">
                            <Eye size={12} /> Xem
                          </button>
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="py-3 px-4 text-neutral-500">{product.sku}</td>
                  <td className="py-3 px-4 text-right">
                    ₫{product.price.toLocaleString("vi-VN")}
                  </td>
                  <td className="py-3 px-4 text-center">{product.stock}</td>
                  <td className="py-3 px-4 text-center text-neutral-500">
                    {product.sold}
                  </td>
                  <td className="py-3 px-4 text-center">
                    {product.status === "active" && (
                      <span className="text-green-600 bg-green-50 px-2 py-1 rounded-xs text-xs">
                        Đang hoạt động
                      </span>
                    )}
                    {product.status === "soldout" && (
                      <span className="text-orange-600 bg-orange-50 px-2 py-1 rounded-xs text-xs">
                        Hết hàng
                      </span>
                    )}
                    {product.status === "disabled" && (
                      <span className="text-neutral-500 bg-neutral-100 px-2 py-1 rounded-xs text-xs">
                        Đã ẩn
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-4 text-center">
                    <div className="flex justify-center gap-2">
                      <button className="text-blue-500 hover:bg-blue-50 p-1.5 rounded-sm">
                        <Edit size={16} />
                      </button>
                      <button className="text-red-500 hover:bg-red-50 p-1.5 rounded-sm">
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="p-4 border-t border-neutral-100 flex justify-end">
        <div className="flex gap-2 text-sm text-neutral-600 items-center">
          <span>1-4 trong 10 sản phẩm</span>
          <button
            className="w-8 h-8 flex items-center justify-center border border-neutral-200 rounded-sm hover:bg-neutral-50 disabled:opacity-50"
            disabled
          >
            &lt;
          </button>
          <button className="w-8 h-8 flex items-center justify-center border border-neutral-200 rounded-sm hover:bg-neutral-50 bg-[#ee4d2d] text-white border-[#ee4d2d]">
            1
          </button>
          <button className="w-8 h-8 flex items-center justify-center border border-neutral-200 rounded-sm hover:bg-neutral-50">
            2
          </button>
          <button className="w-8 h-8 flex items-center justify-center border border-neutral-200 rounded-sm hover:bg-neutral-50">
            &gt;
          </button>
        </div>
      </div>
    </div>
  );
}
