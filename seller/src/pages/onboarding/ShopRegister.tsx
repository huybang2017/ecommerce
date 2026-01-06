import { useState } from "react";

export default function ShopRegister() {
  const steps = [
    { id: 1, title: "Thông tin Shop" },
    { id: 2, title: "Cài đặt vận chuyển" },
    { id: 3, title: "Thông tin định danh" },
    { id: 4, title: "Thông tin thuế" },
    { id: 5, title: "Hoàn tất" },
  ];

  const [currentStep, setCurrentStep] = useState(1);

  return (
    <div className="bg-white rounded-sm shadow-sm min-h-[500px]">
      {/* Stepper */}
      <div className="pt-8 pb-8 px-10 border-b border-neutral-100">
        <div className="flex items-center justify-between">
          {steps.map((step, index) => (
            <div
              key={step.id}
              className="flex flex-col items-center relative flex-1"
            >
              <div className="flex items-center w-full">
                {/* Line Before */}
                <div
                  className={`h-[2px] w-full ${
                    index === 0 ? "bg-transparent" : "bg-neutral-200"
                  } ${currentStep >= step.id ? "bg-[#ee4d2d]/50" : ""}`}
                ></div>

                {/* Dot */}
                <div
                  className={`relative z-10 w-4 h-4 rounded-full flex items-center justify-center ${
                    currentStep >= step.id ? "bg-[#ee4d2d]" : "bg-neutral-200"
                  }`}
                ></div>

                {/* Line After */}
                <div
                  className={`h-[2px] w-full ${
                    index === steps.length - 1
                      ? "bg-transparent"
                      : "bg-neutral-200"
                  } ${currentStep > step.id ? "bg-[#ee4d2d]/50" : ""}`}
                ></div>
              </div>
              <div
                className={`mt-3 text-sm font-medium ${
                  currentStep >= step.id
                    ? "text-neutral-800"
                    : "text-neutral-400"
                }`}
              >
                {step.title}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Form Content */}
      <div className="px-32 py-10">
        <div className="space-y-8">
          {/* Shop Name */}
          <div className="grid grid-cols-12 gap-6 items-center">
            <div className="col-span-3 text-right">
              <label className="text-sm text-neutral-600 font-medium">
                <span className="text-red-500 mr-1">*</span>Tên Shop
              </label>
            </div>
            <div className="col-span-7">
              <div className="relative">
                <input
                  type="text"
                  className="w-full border border-neutral-300 rounded-sm px-3 py-2 text-sm focus:outline-none focus:border-neutral-500"
                  defaultValue="shopwolf"
                />
                <div className="absolute right-3 top-2.5 text-xs text-neutral-400 border-l border-neutral-200 pl-2">
                  8/30
                </div>
              </div>
            </div>
          </div>

          {/* Pickup Address */}
          <div className="grid grid-cols-12 gap-6 items-start">
            <div className="col-span-3 text-right pt-1">
              <label className="text-sm text-neutral-600 font-medium">
                <span className="text-red-500 mr-1">*</span>Địa chỉ lấy hàng
              </label>
            </div>
            <div className="col-span-7">
              <div className="text-sm text-neutral-800 leading-relaxed">
                nguyen duc huy | 84899383561
                <br />
                24 Mai Thị Dũng
                <br />
                Phường Tân Lập
                <br />
                Thành Phố Nha Trang
                <br />
                Khánh Hòa
              </div>
              <button className="text-[#ee4d2d] text-sm hover:underline mt-1 font-medium">
                Chỉnh sửa
              </button>
            </div>
          </div>

          {/* Email */}
          <div className="grid grid-cols-12 gap-6 items-center">
            <div className="col-span-3 text-right">
              <label className="text-sm text-neutral-600 font-medium">
                <span className="text-red-500 mr-1">*</span>Email
              </label>
            </div>
            <div className="col-span-7">
              <input
                type="text"
                className="w-full border border-neutral-200 bg-neutral-100 rounded-sm px-3 py-2 text-sm text-neutral-500 focus:outline-none"
                value="ndhbang2017@gmail.com"
                disabled
              />
            </div>
          </div>

          {/* Phone Number */}
          <div className="grid grid-cols-12 gap-6 items-center">
            <div className="col-span-3 text-right">
              <label className="text-sm text-neutral-600 font-medium">
                <span className="text-red-500 mr-1">*</span>Số điện thoại
              </label>
            </div>
            <div className="col-span-7">
              <div className="flex border border-neutral-300 rounded-sm overflow-hidden">
                <div className="bg-neutral-50 px-3 py-2 text-sm text-neutral-500 border-r border-neutral-200">
                  +84
                </div>
                <input
                  type="text"
                  className="flex-1 px-3 py-2 text-sm focus:outline-none"
                  placeholder="Nhập vào"
                />
              </div>
            </div>
          </div>

          {/* Verification Code */}
          <div className="grid grid-cols-12 gap-6 items-center">
            <div className="col-span-3 text-right"></div>
            <div className="col-span-7">
              <div className="flex gap-4">
                <input
                  type="text"
                  className="w-full border border-red-300 rounded-sm px-3 py-2 text-sm focus:outline-none bg-red-50"
                  placeholder="Nhập vào"
                />
                <button className="whitespace-nowrap px-6 py-2 border border-neutral-200 text-neutral-400 rounded-sm text-sm font-medium hover:bg-neutral-50">
                  Gửi
                </button>
              </div>
              <div className="text-xs text-red-500 mt-1">
                Vui lòng nhập mã xác minh
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Footer Actions */}
      <div className="border-t border-neutral-100 px-10 py-4 flex justify-end gap-3 rounded-b-sm">
        <button className="px-8 py-2 border border-neutral-200 rounded-sm text-sm text-neutral-600 hover:bg-neutral-50 transition-colors">
          Lưu
        </button>
        <button className="px-8 py-2 bg-[#ee4d2d] text-white rounded-sm text-sm hover:bg-[#d73211] transition-colors shadow-sm">
          Tiếp theo
        </button>
      </div>
    </div>
  );
}
