import { Route, Routes } from "react-router-dom";
import SellerLayout from "./layouts/SellerLayout";
import Dashboard from "./pages/Dashboard";
import ProductList from "./pages/products/ProductList";
import OrderList from "./pages/orders/OrderList";
import OnboardingLayout from "./layouts/OnboardingLayout";
import ShopRegister from "./pages/onboarding/ShopRegister";

function App() {
  return (
    <Routes>
      <Route path="/" element={<SellerLayout />}>
        <Route index element={<Dashboard />} />
        <Route path="products/list" element={<ProductList />} />
        <Route path="orders/all" element={<OrderList />} />
        {/* Add more routes here */}
        <Route
          path="*"
          element={
            <div className="p-8 text-center text-neutral-500">
              Trang chưa được xây dựng
            </div>
          }
        />
      </Route>
      <Route path="/onboarding" element={<OnboardingLayout />}>
        <Route index element={<ShopRegister />} />
      </Route>
    </Routes>
  );
}

export default App;
