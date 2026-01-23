package main

import (
	"encoding/json"
	"log"
	"product-service/config"
	"product-service/internal/domain"
	"product-service/internal/repository/postgres"
	"product-service/pkg/database"

	"gorm.io/datatypes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.GetDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	// Initialize repository
	productRepo := postgres.NewProductRepository(db)

	log.Println("Starting to seed products (child categories)...")

	createImagesJSON := func(images []string) datatypes.JSON {
		if len(images) == 0 {
			return nil
		}
		jsonBytes, _ := json.Marshal(images)
		return datatypes.JSON(jsonBytes)
	}

	defaultShopID := uint(1)

	// Category IDs (from API response)
	aoThunNamID := uint(21)
	aoSoMiNamID := uint(22)
	aoKhoacNamID := uint(23)
	quanJeansNamID := uint(24)
	quanShortNamID := uint(25)

	products := []*domain.Product{
		// Áo Thun Nam (21)
		{
			ShopID:      defaultShopID,
			Name:        "Áo Thun Nam Cotton Compact Form Rộng Unisex",
			Description: "Áo thun nam cotton 100%, form rộng thoải mái, nhiều màu",
			BasePrice:   159000,
			CategoryID:  &aoThunNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Thun Nam Polo Trơn Cao Cấp",
			Description: "Áo thun polo nam, chất liệu cotton mềm mại, không xù lông",
			BasePrice:   199000,
			CategoryID:  &aoThunNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Thun Nam Tay Lỡ Form Rộng Streetwear",
			Description: "Áo thun oversize phong cách Hàn Quốc, chất liệu cotton 4 chiều",
			BasePrice:   229000,
			CategoryID:  &aoThunNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},

		// Áo Sơ Mi Nam (22)
		{
			ShopID:      defaultShopID,
			Name:        "Áo Sơ Mi Nam Dài Tay Công Sở",
			Description: "Áo sơ mi nam dài tay, chống nhăn, phù hợp đi làm",
			BasePrice:   299000,
			CategoryID:  &aoSoMiNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Sơ Mi Nam Ngắn Tay Trẻ Trung",
			Description: "Áo sơ mi nam ngắn tay, form fitted hiện đại",
			BasePrice:   249000,
			CategoryID:  &aoSoMiNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},

		// Áo Khoác Nam (23)
		{
			ShopID:      defaultShopID,
			Name:        "Áo Khoác Nam Bomber Jacket 2 Lớp Chống Nước",
			Description: "Áo khoác bomber 2 lớp, chống nước, nhiều màu sắc",
			BasePrice:   599000,
			CategoryID:  &aoKhoacNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Khoác Nam Dù Nhẹ Chống Tia UV",
			Description: "Áo khoác dù siêu nhẹ, chống tia UV, gấp gọn tiện lợi",
			BasePrice:   449000,
			CategoryID:  &aoKhoacNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Khoác Nam Hoodie Nỉ Ngoại Có Mũ",
			Description: "Áo hoodie nỉ ngoại dày dặn, giữ ấm tốt",
			BasePrice:   499000,
			CategoryID:  &aoKhoacNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},

		// Quần Jeans Nam (24)
		{
			ShopID:      defaultShopID,
			Name:        "Quần Jeans Nam Ống Rộng Suông Baggy",
			Description: "Quần jean nam ống rộng, chất liệu denim cao cấp",
			BasePrice:   399000,
			CategoryID:  &quanJeansNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Quần Jeans Nam Ống Đứng Slimfit",
			Description: "Quần jean nam ống đứng, form slimfit ôm vừa vặn",
			BasePrice:   429000,
			CategoryID:  &quanJeansNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},

		// Quần Short Nam (25)
		{
			ShopID:      defaultShopID,
			Name:        "Quần Short Nam Kaki Túi Hộp Thể Thao",
			Description: "Quần short kaki nam, túi hộp tiện dụng, thoáng mát",
			BasePrice:   229000,
			CategoryID:  &quanShortNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Quần Short Nam Jeans Rách Cá Tính",
			Description: "Quần short jeans rách, phong cách năng động trẻ trung",
			BasePrice:   279000,
			CategoryID:  &quanShortNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			IsActive:    true,
		},
	}

	// helper: find product by exact name using ListProducts search filter
	findProductByName := func(repo domain.ProductRepository, name string) (*domain.Product, error) {
		filters := map[string]interface{}{"search": name}
		items, _, err := repo.ListProducts(filters, 1, 1)
		if err != nil {
			return nil, err
		}
		if len(items) == 0 {
			return nil, nil
		}
		// prefer exact match
		for _, p := range items {
			if p.Name == name {
				return p, nil
			}
		}
		return items[0], nil
	}

	for _, product := range products {
		// Check if product already exists by name
		existing, err := findProductByName(productRepo, product.Name)
		if err == nil && existing != nil {
			log.Printf("⏭️  Product already exists: %s", existing.Name)
			continue
		}

		// Create product
		err = productRepo.Create(product)
		if err != nil {
			log.Printf("❌ Failed to create product %s: %v", product.Name, err)
			continue
		}

		if product.CategoryID != nil {
			log.Printf("✅ Created product: %s (CategoryID: %d)", product.Name, *product.CategoryID)
		} else {
			log.Printf("✅ Created product: %s", product.Name)
		}
	}

	log.Println("\n🎉 Seed completed!")
}
