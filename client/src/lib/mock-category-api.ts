export async function getCategoryDetail(slug: string) {
  // Simulate API delay
  await new Promise((resolve) => setTimeout(resolve, 500));

  // Parse ID from slug: name-cat.ID
  // Example: thoi-trang-nam-cat.11035567 -> 11035567
  const parts = slug.split(".");
  const id = parts[parts.length - 1]; // Naive parsing

  // If slug doesn't match pattern, return null (handle 404 in page)
  if (!slug.includes("-cat.")) {
    return null;
  }

  return {
    id: parseInt(id),
    name: "Thời Trang Nam",
    path: [{ name: "Thời Trang Nam", path: `/thoi-trang-nam-cat.${id}` }],
    children: [
      { id: 1, name: "Áo Khoác", slug: `ao-khoac-cat.${id}.1` },
      { id: 2, name: "Áo Vest", slug: `ao-vest-cat.${id}.2` },
      { id: 3, name: "Áo Thun", slug: `ao-thun-cat.${id}.3` },
      { id: 4, name: "Quần Jeans", slug: `quan-jeans-cat.${id}.4` },
      { id: 5, name: "Quần Short", slug: `quan-short-cat.${id}.5` },
      { id: 6, name: "Đồ Lót", slug: `do-lot-cat.${id}.6` },
      { id: 7, name: "Phụ Kiện", slug: `phu-kien-cat.${id}.7` },
    ],
  };
}

export async function getProducts(categoryId: number, filters: any) {
  await new Promise((resolve) => setTimeout(resolve, 800));

  const products = Array.from({ length: 20 }).map((_, i) => ({
    id: i + 1,
    name:
      i % 3 === 0
        ? "Áo Khoác Dù Nam 2 Lớp Chống Nước, Chống Nắng"
        : i % 2 === 0
        ? "Quần Jeans Nam Ống Đứng Cao Cấp"
        : "Áo Thun Nam Cotton Co Giãn 4 Chiều",
    image: "https://placehold.co/300x300",
    price: 150000 + Math.random() * 500000,
    originalPrice: 300000 + Math.random() * 500000,
    rating: 4 + Math.random(),
    sold: Math.floor(Math.random() * 5000),
    discountBadge: Math.random() > 0.5 ? "50%" : undefined,
    location: Math.random() > 0.5 ? "TP. Hồ Chí Minh" : "Hà Nội",
    slug: `product-${i + 1}`,
  }));

  return {
    data: products,
    total: 100,
    page: filters.page || 1,
    last_page: 5,
  };
}
