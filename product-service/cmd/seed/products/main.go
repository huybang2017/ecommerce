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
			Price:       129000,
			BasePrice:   159000,
			SKU:         "AOTHUN-NAM-001",
			CategoryID:  &aoThunNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       200,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Thun Nam Polo Trơn Cao Cấp",
			Description: "Áo thun polo nam, chất liệu cotton mềm mại, không xù lông",
			Price:       149000,
			BasePrice:   199000,
			SKU:         "AOTHUN-NAM-002",
			CategoryID:  &aoThunNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       180,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Thun Nam Tay Lỡ Form Rộng Streetwear",
			Description: "Áo thun oversize phong cách Hàn Quốc, chất liệu cotton 4 chiều",
			Price:       159000,
			BasePrice:   229000,
			SKU:         "AOTHUN-NAM-003",
			CategoryID:  &aoThunNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       220,
			IsActive:    true,
		},

		// Áo Sơ Mi Nam (22)
		{
			ShopID:      defaultShopID,
			Name:        "Áo Sơ Mi Nam Dài Tay Công Sở",
			Description: "Áo sơ mi nam dài tay, chống nhăn, phù hợp đi làm",
			Price:       199000,
			BasePrice:   299000,
			SKU:         "AOSOMI-NAM-001",
			CategoryID:  &aoSoMiNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       150,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Sơ Mi Nam Ngắn Tay Trẻ Trung",
			Description: "Áo sơ mi nam ngắn tay, form fitted hiện đại",
			Price:       169000,
			BasePrice:   249000,
			SKU:         "AOSOMI-NAM-002",
			CategoryID:  &aoSoMiNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       170,
			IsActive:    true,
		},

		// Áo Khoác Nam (23)
		{
			ShopID:      defaultShopID,
			Name:        "Áo Khoác Nam Bomber Jacket 2 Lớp Chống Nước",
			Description: "Áo khoác bomber 2 lớp, chống nước, nhiều màu sắc",
			Price:       459000,
			BasePrice:   599000,
			SKU:         "KHOAC-NAM-001",
			CategoryID:  &aoKhoacNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       80,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Khoác Nam Dù Nhẹ Chống Tia UV",
			Description: "Áo khoác dù siêu nhẹ, chống tia UV, gấp gọn tiện lợi",
			Price:       299000,
			BasePrice:   449000,
			SKU:         "KHOAC-NAM-002",
			CategoryID:  &aoKhoacNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       120,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Khoác Nam Hoodie Nỉ Ngoại Có Mũ",
			Description: "Áo hoodie nỉ ngoại dày dặn, giữ ấm tốt",
			Price:       349000,
			BasePrice:   499000,
			SKU:         "KHOAC-NAM-003",
			CategoryID:  &aoKhoacNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       95,
			IsActive:    true,
		},

		// Quần Jeans Nam (24)
		{
			ShopID:      defaultShopID,
			Name:        "Quần Jeans Nam Ống Rộng Suông Baggy",
			Description: "Quần jean nam ống rộng, chất liệu denim cao cấp",
			Price:       299000,
			BasePrice:   399000,
			SKU:         "JEAN-NAM-001",
			CategoryID:  &quanJeansNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       150,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Quần Jeans Nam Ống Đứng Slimfit",
			Description: "Quần jean nam ống đứng, form slimfit ôm vừa vặn",
			Price:       319000,
			BasePrice:   429000,
			SKU:         "JEAN-NAM-002",
			CategoryID:  &quanJeansNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       160,
			IsActive:    true,
		},

		// Quần Short Nam (25)
		{
			ShopID:      defaultShopID,
			Name:        "Quần Short Nam Kaki Túi Hộp Thể Thao",
			Description: "Quần short kaki nam, túi hộp tiện dụng, thoáng mát",
			Price:       159000,
			BasePrice:   229000,
			SKU:         "SHORT-NAM-001",
			CategoryID:  &quanShortNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       200,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Quần Short Nam Jeans Rách Cá Tính",
			Description: "Quần short jeans rách, phong cách năng động trẻ trung",
			Price:       189000,
			BasePrice:   279000,
			SKU:         "SHORT-NAM-002",
			CategoryID:  &quanShortNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       175,
			IsActive:    true,
		},
	}

	for _, product := range products {
		// Check if product already exists
		existing, err := productRepo.GetBySKU(product.SKU)
		if err == nil && existing != nil {
			log.Printf("⏭️  Product already exists: %s (SKU: %s)", existing.Name, existing.SKU)
			continue
		}

		// Create product
		err = productRepo.Create(product)
		if err != nil {
			log.Printf("❌ Failed to create product %s: %v", product.Name, err)
			continue
		}

		log.Printf("✅ Created product: %s (CategoryID: %d, SKU: %s)", product.Name, *product.CategoryID, product.SKU)
	}

	log.Println("\n🎉 Seed completed!")
}
