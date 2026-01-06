import { Bell, HelpCircle, Grid, User } from "lucide-react";
import { Link } from "react-router-dom";

export default function Header() {
  return (
    <header className="fixed top-0 left-0 right-0 h-14 bg-white border-b border-[#e5e5e5] z-20 flex items-center justify-between px-6 shadow-sm">
      <div className="flex items-center gap-2">
        <Link to="/" className="flex items-center gap-2">
          <div className="text-[#ee4d2d] font-bold text-xl flex items-center gap-1">
            <span className="text-2xl">Shopee</span>
            <span className="text-neutral-800 font-normal">Kênh Người Bán</span>
          </div>
        </Link>
      </div>

      <div className="flex items-center gap-6 text-neutral-600">
        <button className="flex items-center gap-1 hover:text-[#ee4d2d]">
          <Grid size={20} />
        </button>
        <button className="flex items-center gap-1 hover:text-[#ee4d2d]">
          <Bell size={20} />
        </button>
        <button className="flex items-center gap-1 hover:text-[#ee4d2d]">
          <HelpCircle size={20} />
        </button>
        <div className="flex items-center gap-2 pl-4 border-l border-neutral-200 cursor-pointer hover:text-[#ee4d2d]">
          <div className="w-8 h-8 bg-neutral-100 rounded-full flex items-center justify-center overflow-hidden border border-neutral-200">
            <User size={16} />
          </div>
          <span className="text-sm font-medium">nguyenvana123</span>
        </div>
      </div>
    </header>
  );
}
