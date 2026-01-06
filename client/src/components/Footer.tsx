export default function Footer() {
  return (
    <footer className="border-t-[4px] border-[#ee4d2d] bg-[#fbfbfb] text-sm text-neutral-600">
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4 lg:grid-cols-5">
          <div>
            <h3 className="mb-4 text-xs font-bold uppercase text-neutral-700">
              Chăm sóc khách hàng
            </h3>
            <ul className="space-y-2 text-xs">
              <li>Trung tâm trợ giúp</li>
              <li>Shopee Blog</li>
              <li>Shopee Mall</li>
              <li>Hướng dẫn mua hàng</li>
              <li>Hướng dẫn bán hàng</li>
              <li>Thanh toán</li>
              <li>Shopee Xu</li>
              <li>Vận chuyển</li>
              <li>Trả hàng & Hoàn tiền</li>
              <li>Chăm sóc khách hàng</li>
              <li>Chính sách bảo hành</li>
            </ul>
          </div>
          <div>
            <h3 className="mb-4 text-xs font-bold uppercase text-neutral-700">
              Về Shopee
            </h3>
            <ul className="space-y-2 text-xs">
              <li>Giới thiệu về Shopee Việt Nam</li>
              <li>Tuyển dụng</li>
              <li>Điều Khoản Shopee</li>
              <li>Chính sách bảo mật</li>
              <li>Chính Trường Uy Tín</li>
              <li>Kênh Người Bán</li>
              <li>Flash Sales</li>
              <li>Chương trình Tiếp thị liên kết Shopee</li>
              <li>Liên Hệ Với Truyền Thông</li>
            </ul>
          </div>
          <div>
            <h3 className="mb-4 text-xs font-bold uppercase text-neutral-700">
              Thanh toán
            </h3>
            <div className="flex flex-wrap gap-2">
              <div className="h-6 w-10 bg-white shadow-sm border rounded"></div>
              <div className="h-6 w-10 bg-white shadow-sm border rounded"></div>
              <div className="h-6 w-10 bg-white shadow-sm border rounded"></div>
              <div className="h-6 w-10 bg-white shadow-sm border rounded"></div>
            </div>
            <h3 className="mt-8 mb-4 text-xs font-bold uppercase text-neutral-700">
              Đơn vị vận chuyển
            </h3>
            <div className="flex flex-wrap gap-2">
              <div className="h-6 w-10 bg-white shadow-sm border rounded"></div>
              <div className="h-6 w-10 bg-white shadow-sm border rounded"></div>
              <div className="h-6 w-10 bg-white shadow-sm border rounded"></div>
            </div>
          </div>
          <div>
            <h3 className="mb-4 text-xs font-bold uppercase text-neutral-700">
              Theo dõi chúng tôi
            </h3>
            <ul className="space-y-2 text-xs">
              <li>Facebook</li>
              <li>Instagram</li>
              <li>LinkedIn</li>
            </ul>
          </div>
          <div>
            <h3 className="mb-4 text-xs font-bold uppercase text-neutral-700">
              Tải ứng dụng
            </h3>
            <div className="flex gap-2">
              <div className="h-20 w-20 bg-gray-200">QR Code</div>
              <div className="flex flex-col gap-1 justify-center">
                <div className="h-6 w-20 bg-gray-200">App Store</div>
                <div className="h-6 w-20 bg-gray-200">Google Play</div>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-12 border-t pt-8 text-center text-xs text-neutral-500">
          <div className="mb-4 flex justify-center gap-4">
            <span>CHÍNH SÁCH BẢO MẬT</span>
            <span>QUY CHẾ HOẠT ĐỘNG</span>
            <span>CHÍNH SÁCH VẬN CHUYỂN</span>
            <span>CHÍNH SÁCH TRẢ HÀNG VÀ HOÀN TIỀN</span>
          </div>
          <p className="mb-2">Công ty TNHH Shopee</p>
          <p>
            &copy; {new Date().getFullYear()} Shopee. Tất cả các quyền được bảo
            lưu.
          </p>
        </div>
      </div>
    </footer>
  );
}
