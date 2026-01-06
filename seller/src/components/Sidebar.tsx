import { useState } from "react";
import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  Package,
  ShoppingBag,
  FileText,
  Settings,
  ChevronDown,
  ChevronRight,
  Store,
} from "lucide-react";

interface MenuItem {
  title: string;
  icon: any;
  path?: string;
  subItems?: { title: string; path: string }[];
}

const menuItems: MenuItem[] = [
  {
    title: "Vận chuyển",
    icon: Package,
    subItems: [
      { title: "Quản Lý Vận Chuyển", path: "/shipment/mass-ship" },
      { title: "Giao Hàng Loạt", path: "/shipment/shipping" },
      { title: "Cài Đặt Vận Chuyển", path: "/shipment/settings" },
    ],
  },
  {
    title: "Quản Lý Đơn Hàng",
    icon: FileText,
    subItems: [
      { title: "Tất Cả", path: "/orders/all" },
      { title: "Đơn Hủy", path: "/orders/cancellation" },
      { title: "Trả Hàng/Hoàn Tiền", path: "/orders/return" },
    ],
  },
  {
    title: "Quản Lý Sản Phẩm",
    icon: ShoppingBag,
    subItems: [
      { title: "Tất Cả Sản Phẩm", path: "/products/list" },
      { title: "Thêm Sản Phẩm", path: "/products/add" },
      { title: "Cài Đặt Sản Phẩm", path: "/products/settings" },
    ],
  },
  {
    title: "Kênh Marketing",
    icon: Store,
    path: "/marketing",
  },
  {
    title: "Hồ Sơ Shop",
    icon: Settings,
    path: "/shop/profile",
  },
];

export default function Sidebar() {
  const [openMenus, setOpenMenus] = useState<string[]>([
    "Vận chuyển",
    "Quản Lý Đơn Hàng",
    "Quản Lý Sản Phẩm",
  ]);

  const toggleMenu = (title: string) => {
    setOpenMenus((prev) =>
      prev.includes(title)
        ? prev.filter((item) => item !== title)
        : [...prev, title]
    );
  };

  return (
    <aside className="w-[256px] bg-white border-r border-[#e5e5e5] h-screen overflow-y-auto fixed left-0 top-0 pt-16 pb-10 z-10 flex flex-col">
      <div className="px-4 py-4">
        {menuItems.map((item) => (
          <div key={item.title} className="mb-2">
            {item.subItems ? (
              <div>
                <button
                  onClick={() => toggleMenu(item.title)}
                  className="flex items-center justify-between w-full px-2 py-2 text-sm font-bold text-neutral-700 hover:bg-neutral-50 rounded-sm"
                >
                  <div className="flex items-center gap-2">
                    <item.icon size={18} className="text-neutral-500" />
                    <span>{item.title}</span>
                  </div>
                  {openMenus.includes(item.title) ? (
                    <ChevronDown size={14} />
                  ) : (
                    <ChevronRight size={14} />
                  )}
                </button>

                {openMenus.includes(item.title) && (
                  <div className="ml-9 mt-1 space-y-1">
                    {item.subItems.map((sub) => (
                      <NavLink
                        key={sub.path}
                        to={sub.path}
                        className={({ isActive }) =>
                          `block text-sm py-1.5 px-2 rounded-sm hover:text-[#ee4d2d] transition-colors ${
                            isActive
                              ? "text-[#ee4d2d] font-medium"
                              : "text-neutral-600"
                          }`
                        }
                      >
                        {sub.title}
                      </NavLink>
                    ))}
                  </div>
                )}
              </div>
            ) : (
              <NavLink
                to={item.path || "#"}
                className={({ isActive }) =>
                  `flex items-center gap-2 px-2 py-2 text-sm font-bold text-neutral-700 hover:bg-neutral-50 rounded-sm ${
                    isActive ? "text-[#ee4d2d]" : ""
                  }`
                }
              >
                <item.icon size={18} className="text-neutral-500" />
                <span>{item.title}</span>
              </NavLink>
            )}
          </div>
        ))}
      </div>
    </aside>
  );
}
