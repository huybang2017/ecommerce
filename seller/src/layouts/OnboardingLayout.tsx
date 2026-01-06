import { Outlet } from "react-router-dom";
import { Bell, HelpCircle, Grid } from "lucide-react";

export default function OnboardingLayout() {
  return (
    <div className="min-h-screen bg-[#f5f5f5] flex flex-col font-sans">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-neutral-200">
        <div className="max-w-[1200px] mx-auto px-4 h-[60px] flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="text-[#ee4d2d] text-3xl font-bold tracking-tighter flex items-center gap-1">
                <ShoppingBagIcon /> Shopee
              </span>
              <span className="text-xl text-neutral-800">
                Đăng ký trở thành Người bán Shopee
              </span>
            </div>
          </div>
          <div className="flex items-center gap-6 text-neutral-600">
            <div className="flex items-center gap-1 cursor-pointer hover:text-[#ee4d2d]">
              <img
                src="https://placehold.co/30x30"
                alt="Avatar"
                className="w-8 h-8 rounded-full border border-neutral-200"
              />
              <span className="text-sm font-medium">1pe_rbl3lz</span>
            </div>
            <div className="flex items-center gap-4">
              <Bell size={20} className="cursor-pointer hover:text-[#ee4d2d]" />
              <HelpCircle
                size={20}
                className="cursor-pointer hover:text-[#ee4d2d]"
              />
              <Grid size={20} className="cursor-pointer hover:text-[#ee4d2d]" />
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 py-6">
        <div className="max-w-[1200px] mx-auto px-4">
          <Outlet />
        </div>
      </main>
    </div>
  );
}

function ShoppingBagIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="currentColor"
      className="w-8 h-8"
    >
      <path
        fillRule="evenodd"
        d="M7.5 6v.75H5.513c-.96 0-1.764.724-1.865 1.679l-1.263 12A1.875 1.875 0 004.25 22.5h15.5a1.875 1.875 0 001.865-2.071l-1.263-12a1.875 1.875 0 00-1.865-1.679H16.5V6a4.5 4.5 0 10-9 0zM12 3a3 3 0 00-3 3v.75h6V6a3 3 0 00-3-3zm-3 8.25a3 3 0 106 0v-.75a.75.75 0 011.5 0v.75a4.5 4.5 0 11-9 0v-.75a.75.75 0 011.5 0v.75z"
        clipRule="evenodd"
      />
    </svg>
  );
}
