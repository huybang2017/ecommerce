import Link from "next/link";
import { ChevronRight } from "lucide-react";

interface BreadcrumbProps {
  items: {
    name: string;
    path: string;
  }[];
}

export function Breadcrumb({ items }: BreadcrumbProps) {
  return (
    <nav className="text-sm text-neutral-600 mb-4 flex items-center gap-2">
      <Link href="/" className="hover:text-[#ee4d2d] transition-colors">
        Trang chủ
      </Link>
      {items.map((item, index) => (
        <div key={item.path} className="flex items-center gap-2">
          <ChevronRight size={14} className="text-neutral-400" />
          <Link
            href={item.path}
            className={`hover:text-[#ee4d2d] transition-colors ${
              index === items.length - 1 ? "font-medium text-neutral-800" : ""
            }`}
          >
            {item.name}
          </Link>
        </div>
      ))}
    </nav>
  );
}
