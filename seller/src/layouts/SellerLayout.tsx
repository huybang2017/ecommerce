import { Outlet } from "react-router-dom";
import Sidebar from "../components/Sidebar";
import Header from "../components/Header";

export default function SellerLayout() {
  return (
    <div className="min-h-screen bg-[#f5f5f5]">
      <Header />
      <Sidebar />
      <main className="ml-[256px] pt-14 p-6 min-h-[calc(100vh-56px)]">
        <div className="max-w-[1200px] mx-auto">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
