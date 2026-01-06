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

	// Initialize repositories
	categoryRepo := postgres.NewCategoryRepository(db)
	productRepo := postgres.NewProductRepository(db)
	variationRepo := postgres.NewVariationRepository(db)
	variationOptRepo := postgres.NewVariationOptionRepository(db)
	productItemRepo := postgres.NewProductItemRepository(db)
	skuConfigRepo := postgres.NewSKUConfigurationRepository(db)
	categoryAttrRepo := postgres.NewCategoryAttributeRepository(db)
	productAttrRepo := postgres.NewProductAttributeValueRepository(db)

	log.Println("Starting to seed data...")

	// 1. Create Categories
	log.Println("Creating categories...")
	categories := []*domain.Category{
		// Thời Trang Nam
		{
			Name:        "Thời Trang Nam",
			Slug:        "thoi-trang-nam",
			Description: "Quần áo và phụ kiện dành cho nam giới",
			IsActive:    true,
		},
		// Thời Trang Nữ
		{
			Name:        "Thời Trang Nữ",
			Slug:        "thoi-trang-nu",
			Description: "Quần áo và phụ kiện dành cho nữ giới",
			IsActive:    true,
		},
		// Điện Thoại & Phụ Kiện
		{
			Name:        "Điện Thoại & Phụ Kiện",
			Slug:        "dien-thoai-phu-kien",
			Description: "Điện thoại di động, máy tính bảng và phụ kiện",
			IsActive:    true,
		},
		// Mẹ & Bé
		{
			Name:        "Mẹ & Bé",
			Slug:        "me-va-be",
			Description: "Sản phẩm cho mẹ và bé",
			IsActive:    true,
		},
		// Thiết Bị Điện Tử
		{
			Name:        "Thiết Bị Điện Tử",
			Slug:        "thiet-bi-dien-tu",
			Description: "Laptop, máy tính, camera, thiết bị âm thanh",
			IsActive:    true,
		},
		// Nhà Cửa & Đời Sống
		{
			Name:        "Nhà Cửa & Đời Sống",
			Slug:        "nha-cua-doi-song",
			Description: "Đồ gia dụng, nội thất, trang trí nhà cửa",
			IsActive:    true,
		},
		// Sắc Đẹp
		{
			Name:        "Sắc Đẹp",
			Slug:        "sac-dep",
			Description: "Mỹ phẩm, chăm sóc da, trang điểm",
			IsActive:    true,
		},
		// Sức Khỏe
		{
			Name:        "Sức Khỏe",
			Slug:        "suc-khoe",
			Description: "Thực phẩm chức năng, thiết bị y tế",
			IsActive:    true,
		},
		// Giày Dép Nam
		{
			Name:        "Giày Dép Nam",
			Slug:        "giay-dep-nam",
			Description: "Giày thể thao, giày tây, dép nam",
			IsActive:    true,
		},
		// Giày Dép Nữ
		{
			Name:        "Giày Dép Nữ",
			Slug:        "giay-dep-nu",
			Description: "Giày cao gót, giày thể thao, sandal nữ",
			IsActive:    true,
		},
		// Túi Ví Nam
		{
			Name:        "Túi Ví Nam",
			Slug:        "tui-vi-nam",
			Description: "Balo, cặp da, ví nam",
			IsActive:    true,
		},
		// Túi Ví Nữ
		{
			Name:        "Túi Ví Nữ",
			Slug:        "tui-vi-nu",
			Description: "Túi xách, ví nữ, clutch",
			IsActive:    true,
		},
		// Đồng Hồ
		{
			Name:        "Đồng Hồ",
			Slug:        "dong-ho",
			Description: "Đồng hồ nam, nữ, trẻ em",
			IsActive:    true,
		},
		// Thể Thao & Du Lịch
		{
			Name:        "Thể Thao & Du Lịch",
			Slug:        "the-thao-du-lich",
			Description: "Dụng cụ thể thao, đồ du lịch",
			IsActive:    true,
		},
		// Ô Tô & Xe Máy & Xe Đạp
		{
			Name:        "Ô Tô & Xe Máy & Xe Đạp",
			Slug:        "o-to-xe-may-xe-dap",
			Description: "Phụ kiện và thiết bị cho xe",
			IsActive:    true,
		},
	}

	var createdCategories []*domain.Category
	for _, cat := range categories {
		// Check if category already exists
		existing, err := categoryRepo.GetBySlug(cat.Slug)
		if err == nil && existing != nil {
			createdCategories = append(createdCategories, existing)
			log.Printf("⏭️  Using existing category: %s (ID: %d)", existing.Name, existing.ID)
			continue
		}

		// Create new category
		err = categoryRepo.Create(cat)
		if err != nil {
			log.Printf("❌ Failed to create category %s: %v", cat.Name, err)
			continue
		}

		// Get the created category to get its ID
		created, err := categoryRepo.GetBySlug(cat.Slug)
		if err != nil {
			log.Printf("⚠️  Created category %s but failed to retrieve it: %v", cat.Name, err)
			continue
		}

		createdCategories = append(createdCategories, created)
		log.Printf("✅ Created category: %s (ID: %d)", created.Name, created.ID)
	}

	if len(createdCategories) == 0 {
		log.Fatal("No categories available. Cannot create products.")
	}

	// Get category IDs
	thoiTrangNamID := createdCategories[0].ID
	thoiTrangNuID := createdCategories[1].ID
	dienThoaiID := createdCategories[2].ID
	// meBeID := createdCategories[3].ID
	thietBiDienTuID := createdCategories[4].ID
	nhaCuaID := createdCategories[5].ID
	sacDepID := createdCategories[6].ID
	// sucKhoeID := createdCategories[7].ID
	giayNamID := createdCategories[8].ID
	// giayNuID := createdCategories[9].ID
	tuiNamID := createdCategories[10].ID
	// tuiNuID := createdCategories[11].ID
	// dongHoID := createdCategories[12].ID
	// theThaoID := createdCategories[13].ID
	// xeID := createdCategories[14].ID

	// 2. Create Products
	log.Println("\nCreating products...")

	// Helper function to create images JSON
	createImagesJSON := func(images []string) datatypes.JSON {
		if len(images) == 0 {
			return nil
		}
		jsonBytes, _ := json.Marshal(images)
		return datatypes.JSON(jsonBytes)
	}

	// Note: ShopID defaults to 1 for all products (can be changed later)
	defaultShopID := uint(1)

	products := []*domain.Product{
		// Thời Trang Nam
		{
			ShopID:      defaultShopID,
			Name:        "Áo Thun Nam Cotton Compact Form Rộng Unisex",
			Description: "Áo thun nam cotton 100%, form rộng thoải mái, nhiều màu",
			Price:       129000,
			BasePrice:   159000,
			SKU:         "AOTHUN-NAM-001",
			CategoryID:  &thoiTrangNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       200,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Quần Jeans Nam Ống Rộng Suông Baggy",
			Description: "Quần jean nam ống rộng, chất liệu denim cao cấp",
			Price:       299000,
			BasePrice:   399000,
			SKU:         "JEAN-NAM-001",
			CategoryID:  &thoiTrangNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       150,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Khoác Nam Bomber Jacket 2 Lớp Chống Nước",
			Description: "Áo khoác bomber 2 lớp, chống nước, nhiều màu sắc",
			Price:       459000,
			BasePrice:   599000,
			SKU:         "KHOAC-NAM-001",
			CategoryID:  &thoiTrangNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       80,
			IsActive:    true,
		},
		// Thời Trang Nữ
		{
			ShopID:      defaultShopID,
			Name:        "Váy Babydoll Hoa Nhí Tay Bồng",
			Description: "Váy babydoll dáng xòe, họa tiết hoa nhí xinh xắn",
			Price:       189000,
			BasePrice:   249000,
			SKU:         "VAY-NU-001",
			CategoryID:  &thoiTrangNuID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       120,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Áo Kiểu Nữ Dài Tay Công Sở",
			Description: "Áo kiểu nữ dài tay, chất liệu lụa mềm mại",
			Price:       159000,
			BasePrice:   199000,
			SKU:         "AOKIEU-NU-001",
			CategoryID:  &thoiTrangNuID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       180,
			IsActive:    true,
		},
		// Điện Thoại
		{
			ShopID:      defaultShopID,
			Name:        "iPhone 15 Pro Max 256GB Chính Hãng VN/A",
			Description: "iPhone 15 Pro Max - Chip A17 Pro, Camera 48MP, Màn hình 6.7 inch",
			Price:       29990000,
			BasePrice:   33990000,
			SKU:         "IPHONE15PM-256",
			CategoryID:  &dienThoaiID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       50,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Samsung Galaxy S24 Ultra 12GB/256GB",
			Description: "Galaxy S24 Ultra - Snapdragon 8 Gen 3, Camera 200MP, S Pen",
			Price:       26990000,
			BasePrice:   31990000,
			SKU:         "SAMSUNG-S24U-256",
			CategoryID:  &dienThoaiID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       60,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Xiaomi Redmi Note 13 Pro 8GB/256GB",
			Description: "Redmi Note 13 Pro - Camera 200MP, Màn hình AMOLED 120Hz",
			Price:       6990000,
			BasePrice:   8990000,
			SKU:         "XIAOMI-RN13P-256",
			CategoryID:  &dienThoaiID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       200,
			IsActive:    true,
		},
		// Thiết Bị Điện Tử
		{
			ShopID:      defaultShopID,
			Name:        "Laptop Dell Inspiron 15 3520 i5-1235U/8GB/512GB",
			Description: "Dell Inspiron 15 - Intel Core i5 Gen 12, RAM 8GB, SSD 512GB",
			Price:       13990000,
			BasePrice:   16990000,
			SKU:         "DELL-INS15-3520",
			CategoryID:  &thietBiDienTuID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       40,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Tai Nghe Bluetooth Sony WH-1000XM5",
			Description: "Tai nghe chống ồn chủ động hàng đầu, pin 30 giờ",
			Price:       7990000,
			BasePrice:   9990000,
			SKU:         "SONY-WH1000XM5",
			CategoryID:  &thietBiDienTuID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       75,
			IsActive:    true,
		},
		// Giày Dép Nam
		{
			ShopID:      defaultShopID,
			Name:        "Giày Sneaker Nam Thể Thao Cổ Thấp",
			Description: "Giày sneaker nam, đế cao su, êm ái thoáng khí",
			Price:       249000,
			BasePrice:   399000,
			SKU:         "GIAY-NAM-001",
			CategoryID:  &giayNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       300,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Dép Quai Ngang Nam Nữ Unisex",
			Description: "Dép quai ngang đế êm, chống trơn trượt",
			Price:       89000,
			BasePrice:   129000,
			SKU:         "DEP-NAM-001",
			CategoryID:  &giayNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       500,
			IsActive:    true,
		},
		// Túi Ví Nam
		{
			ShopID:      defaultShopID,
			Name:        "Balo Laptop 15.6 inch Chống Nước",
			Description: "Balo laptop đa ngăn, chống nước, chống sốc",
			Price:       299000,
			BasePrice:   449000,
			SKU:         "BALO-NAM-001",
			CategoryID:  &tuiNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       100,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Ví Da Nam Cao Cấp Đựng Thẻ ATM",
			Description: "Ví da bò thật, nhiều ngăn đựng thẻ tiện lợi",
			Price:       159000,
			BasePrice:   259000,
			SKU:         "VI-NAM-001",
			CategoryID:  &tuiNamID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       200,
			IsActive:    true,
		},
		// Sắc Đẹp
		{
			ShopID:      defaultShopID,
			Name:        "Kem Chống Nắng Anessa SPF50+ PA++++",
			Description: "Kem chống nắng Nhật Bản, chống nước, lâu trôi",
			Price:       459000,
			BasePrice:   599000,
			SKU:         "ANESSA-SPF50",
			CategoryID:  &sacDepID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       150,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Son Kem Lì 3CE Velvet Lip Tint",
			Description: "Son kem lì Hàn Quốc, lên màu chuẩn, bền màu",
			Price:       249000,
			BasePrice:   329000,
			SKU:         "3CE-VELVET-001",
			CategoryID:  &sacDepID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       250,
			IsActive:    true,
		},
		// Nhà Cửa & Đời Sống
		{
			ShopID:      defaultShopID,
			Name:        "Nồi Cơm Điện Tử Sharp 1.8L",
			Description: "Nồi cơm điện tử công nghệ Nhật, lòng chống dính",
			Price:       1290000,
			BasePrice:   1690000,
			SKU:         "SHARP-RC18",
			CategoryID:  &nhaCuaID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       60,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Đèn LED Thông Minh Xiaomi",
			Description: "Đèn LED điều khiển qua app, 16 triệu màu",
			Price:       299000,
			BasePrice:   499000,
			SKU:         "XIAOMI-LED-001",
			CategoryID:  &nhaCuaID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://placehold.co/400x400"}),
			Stock:       180,
			IsActive:    true,
		},
	}

	createdCount := 0
	skippedCount := 0
	for _, product := range products {
		// Check if product with same SKU already exists
		existing, err := productRepo.GetBySKU(product.SKU)
		if err == nil && existing != nil {
			log.Printf("⏭️  Skipped product (already exists): %s (SKU: %s)", product.Name, product.SKU)
			skippedCount++
			continue
		}

		err = productRepo.Create(product)
		if err != nil {
			log.Printf("❌ Failed to create product %s: %v", product.Name, err)
			continue
		}

		// Get the created product to get its ID
		created, err := productRepo.GetBySKU(product.SKU)
		if err != nil {
			log.Printf("⚠️  Created product %s but failed to retrieve it: %v", product.Name, err)
			createdCount++
			continue
		}

		createdCount++
		log.Printf("✅ Created product: %s (ID: %d, SKU: %s, Price: $%.2f)",
			created.Name, created.ID, created.SKU, created.Price)
	}

	log.Printf("\n=== Seeding Complete ===")
	log.Printf("Categories: %d created/used", len(createdCategories))
	log.Printf("Products: %d created, %d skipped", createdCount, skippedCount)

	// 3. Seed Variations and ProductItems for some products
	log.Println("\n=== Seeding Product Items (SKUs) ===")
	seedProductItems(productRepo, variationRepo, variationOptRepo, productItemRepo, skuConfigRepo)

	// 4. Seed Category Attributes and Product Attribute Values
	log.Println("\n=== Seeding Category Attributes & Product Attributes ===")
	seedCategoryAndProductAttributes(categoryRepo, productRepo, categoryAttrRepo, productAttrRepo, createdCategories)

	log.Println("\n✅ Data seeding finished!")
}

func seedProductItems(
	productRepo domain.ProductRepository,
	variationRepo domain.VariationRepository,
	variationOptRepo domain.VariationOptionRepository,
	productItemRepo domain.ProductItemRepository,
	skuConfigRepo domain.SKUConfigurationRepository,
) {
	// Get some products to add variations
	aoThun, _ := productRepo.GetBySKU("AOTHUN-NAM-001")
	iphone, _ := productRepo.GetBySKU("IPHONE15PM-256")
	giay, _ := productRepo.GetBySKU("GIAY-NAM-001")

	if aoThun == nil || iphone == nil || giay == nil {
		log.Println("⚠️  Required products not found, skipping product items seeding")
		return
	}

	// ============= Áo Thun - Size + Color variations =============
	log.Printf("\n--- Creating variations for: %s ---", aoThun.Name)

	// Create Variations
	sizeVar := &domain.Variation{ProductID: aoThun.ID, Name: "Kích Thước"}
	variationRepo.Create(sizeVar)

	colorVar := &domain.Variation{ProductID: aoThun.ID, Name: "Màu Sắc"}
	variationRepo.Create(colorVar)

	// Create Variation Options - Size
	sizeM := &domain.VariationOption{VariationID: sizeVar.ID, Value: "M"}
	sizeL := &domain.VariationOption{VariationID: sizeVar.ID, Value: "L"}
	sizeXL := &domain.VariationOption{VariationID: sizeVar.ID, Value: "XL"}
	variationOptRepo.Create(sizeM)
	variationOptRepo.Create(sizeL)
	variationOptRepo.Create(sizeXL)

	// Create Variation Options - Color
	colorWhite := &domain.VariationOption{VariationID: colorVar.ID, Value: "Trắng"}
	colorBlack := &domain.VariationOption{VariationID: colorVar.ID, Value: "Đen"}
	colorGray := &domain.VariationOption{VariationID: colorVar.ID, Value: "Xám"}
	variationOptRepo.Create(colorWhite)
	variationOptRepo.Create(colorBlack)
	variationOptRepo.Create(colorGray)

	aoThunSKUs := []struct {
		size  *domain.VariationOption
		color *domain.VariationOption
		sku   string
		price float64
		stock int
	}{
		{sizeM, colorWhite, "AOTHUN-NAM-M-TRANG", 129000, 50},
		{sizeM, colorBlack, "AOTHUN-NAM-M-DEN", 129000, 45},
		{sizeL, colorWhite, "AOTHUN-NAM-L-TRANG", 129000, 60},
		{sizeL, colorBlack, "AOTHUN-NAM-L-DEN", 129000, 55},
		{sizeL, colorGray, "AOTHUN-NAM-L-XAM", 129000, 40},
		{sizeXL, colorWhite, "AOTHUN-NAM-XL-TRANG", 139000, 35},
		{sizeXL, colorBlack, "AOTHUN-NAM-XL-DEN", 139000, 30},
	}

	for _, item := range aoThunSKUs {
		productItem := &domain.ProductItem{
			ProductID:  aoThun.ID,
			SKUCode:    item.sku,
			ImageURL:   "https://placehold.co/400x400",
			Price:      item.price,
			QtyInStock: item.stock,
			Status:     "ACTIVE",
		}

		if err := productItemRepo.Create(productItem); err != nil {
			log.Printf("⏭️  SKU %s already exists", item.sku)
			continue
		}

		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.size.ID,
		})
		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.color.ID,
		})

		log.Printf("✅ Created SKU: %s - Size %s, Màu %s (stock: %d)",
			item.sku, item.size.Value, item.color.Value, item.stock)
	}

	// ============= iPhone - Storage + Color variations =============
	log.Printf("\n--- Creating variations for: %s ---", iphone.Name)

	storageVar := &domain.Variation{ProductID: iphone.ID, Name: "Bộ Nhớ"}
	variationRepo.Create(storageVar)

	iphoneColorVar := &domain.Variation{ProductID: iphone.ID, Name: "Màu Sắc"}
	variationRepo.Create(iphoneColorVar)

	// Storage options
	storage256 := &domain.VariationOption{VariationID: storageVar.ID, Value: "256GB"}
	storage512 := &domain.VariationOption{VariationID: storageVar.ID, Value: "512GB"}
	storage1TB := &domain.VariationOption{VariationID: storageVar.ID, Value: "1TB"}
	variationOptRepo.Create(storage256)
	variationOptRepo.Create(storage512)
	variationOptRepo.Create(storage1TB)

	// Color options
	titaniumNatural := &domain.VariationOption{VariationID: iphoneColorVar.ID, Value: "Titan Tự Nhiên"}
	titaniumBlue := &domain.VariationOption{VariationID: iphoneColorVar.ID, Value: "Titan Xanh"}
	titaniumBlack := &domain.VariationOption{VariationID: iphoneColorVar.ID, Value: "Titan Đen"}
	variationOptRepo.Create(titaniumNatural)
	variationOptRepo.Create(titaniumBlue)
	variationOptRepo.Create(titaniumBlack)

	iphoneSKUs := []struct {
		storage *domain.VariationOption
		color   *domain.VariationOption
		sku     string
		price   float64
		stock   int
	}{
		{storage256, titaniumNatural, "IP15PM-256-NATURAL", 29990000, 10},
		{storage256, titaniumBlue, "IP15PM-256-BLUE", 29990000, 8},
		{storage512, titaniumNatural, "IP15PM-512-NATURAL", 34990000, 5},
		{storage512, titaniumBlack, "IP15PM-512-BLACK", 34990000, 6},
		{storage1TB, titaniumBlue, "IP15PM-1TB-BLUE", 39990000, 3},
	}

	for _, item := range iphoneSKUs {
		productItem := &domain.ProductItem{
			ProductID:  iphone.ID,
			SKUCode:    item.sku,
			ImageURL:   "https://placehold.co/400x400",
			Price:      item.price,
			QtyInStock: item.stock,
			Status:     "ACTIVE",
		}

		if err := productItemRepo.Create(productItem); err != nil {
			log.Printf("⏭️  SKU %s already exists", item.sku)
			continue
		}

		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.storage.ID,
		})
		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.color.ID,
		})

		log.Printf("✅ Created SKU: %s - %s %s (%d VNĐ, stock: %d)",
			item.sku, item.storage.Value, item.color.Value, int(item.price), item.stock)
	}

	// ============= Giày - Size + Color variations =============
	log.Printf("\n--- Creating variations for: %s ---", giay.Name)

	giaySizeVar := &domain.Variation{ProductID: giay.ID, Name: "Kích Thước"}
	variationRepo.Create(giaySizeVar)

	giayColorVar := &domain.Variation{ProductID: giay.ID, Name: "Màu Sắc"}
	variationRepo.Create(giayColorVar)

	// Sizes
	size39 := &domain.VariationOption{VariationID: giaySizeVar.ID, Value: "39"}
	size40 := &domain.VariationOption{VariationID: giaySizeVar.ID, Value: "40"}
	size41 := &domain.VariationOption{VariationID: giaySizeVar.ID, Value: "41"}
	size42 := &domain.VariationOption{VariationID: giaySizeVar.ID, Value: "42"}
	variationOptRepo.Create(size39)
	variationOptRepo.Create(size40)
	variationOptRepo.Create(size41)
	variationOptRepo.Create(size42)

	// Colors
	giayWhite := &domain.VariationOption{VariationID: giayColorVar.ID, Value: "Trắng"}
	giayBlack := &domain.VariationOption{VariationID: giayColorVar.ID, Value: "Đen"}
	giayRed := &domain.VariationOption{VariationID: giayColorVar.ID, Value: "Đỏ"}
	variationOptRepo.Create(giayWhite)
	variationOptRepo.Create(giayBlack)
	variationOptRepo.Create(giayRed)

	giaySKUs := []struct {
		size  *domain.VariationOption
		color *domain.VariationOption
		sku   string
		stock int
	}{
		{size39, giayWhite, "GIAY-39-TRANG", 25},
		{size40, giayWhite, "GIAY-40-TRANG", 30},
		{size40, giayBlack, "GIAY-40-DEN", 28},
		{size41, giayWhite, "GIAY-41-TRANG", 35},
		{size41, giayBlack, "GIAY-41-DEN", 32},
		{size41, giayRed, "GIAY-41-DO", 20},
		{size42, giayBlack, "GIAY-42-DEN", 22},
	}

	for _, item := range giaySKUs {
		productItem := &domain.ProductItem{
			ProductID:  giay.ID,
			SKUCode:    item.sku,
			ImageURL:   "https://placehold.co/400x400",
			Price:      249000,
			QtyInStock: item.stock,
			Status:     "ACTIVE",
		}

		if err := productItemRepo.Create(productItem); err != nil {
			log.Printf("⏭️  SKU %s already exists", item.sku)
			continue
		}

		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.size.ID,
		})
		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.color.ID,
		})

		log.Printf("✅ Created SKU: %s - Size %s, Màu %s (stock: %d)",
			item.sku, item.size.Value, item.color.Value, item.stock)
	}

	log.Println("\n✅ Product items seeding completed!")
}

func seedCategoryAndProductAttributes(
	categoryRepo domain.CategoryRepository,
	productRepo domain.ProductRepository,
	categoryAttrRepo domain.CategoryAttributeRepository,
	productAttrRepo domain.ProductAttributeValueRepository,
	categories []*domain.Category,
) {
	if len(categories) == 0 {
		log.Println("⚠️  No categories found")
		return
	}

	// Map categories by slug for easy access
	categoryMap := make(map[string]*domain.Category)
	for _, cat := range categories {
		categoryMap[cat.Slug] = cat
	}

	// ============= Electronics Category Attributes =============
	if electronics, ok := categoryMap["electronics"]; ok {
		log.Printf("\n--- Creating attributes for category: %s ---", electronics.Name)

		// Define attributes for Electronics
		electronicsAttrs := []*domain.CategoryAttribute{
			{
				CategoryID:    electronics.ID,
				AttributeName: "Brand",
				InputType:     "text",
				IsMandatory:   true,
				IsFilterable:  true,
			},
			{
				CategoryID:    electronics.ID,
				AttributeName: "Screen Size",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  true,
			},
			{
				CategoryID:    electronics.ID,
				AttributeName: "Processor",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  false,
			},
			{
				CategoryID:    electronics.ID,
				AttributeName: "RAM",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  true,
			},
			{
				CategoryID:    electronics.ID,
				AttributeName: "Storage",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  true,
			},
			{
				CategoryID:    electronics.ID,
				AttributeName: "Battery",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  false,
			},
			{
				CategoryID:    electronics.ID,
				AttributeName: "Warranty",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  false,
			},
		}

		// Create attributes
		attrMap := make(map[string]*domain.CategoryAttribute)
		for _, attr := range electronicsAttrs {
			if err := categoryAttrRepo.Create(attr); err != nil {
				log.Printf("⏭️  Attribute %s already exists or error: %v", attr.AttributeName, err)
				// Try to get existing
				existing, _ := categoryAttrRepo.GetByCategoryID(electronics.ID)
				for _, e := range existing {
					if e.AttributeName == attr.AttributeName {
						attrMap[attr.AttributeName] = e
						break
					}
				}
			} else {
				attrMap[attr.AttributeName] = attr
				log.Printf("✅ Created attribute: %s", attr.AttributeName)
			}
		}

		// Add attribute values for iPhone 15 Pro
		iphone, _ := productRepo.GetBySKU("IPH15P-001")
		if iphone != nil && len(attrMap) > 0 {
			log.Printf("\n--- Adding attributes for: %s ---", iphone.Name)

			iphoneAttrs := []*domain.ProductAttributeValue{
				{ProductID: iphone.ID, AttributeID: attrMap["Brand"].ID, Value: "Apple"},
				{ProductID: iphone.ID, AttributeID: attrMap["Screen Size"].ID, Value: "6.1 inch"},
				{ProductID: iphone.ID, AttributeID: attrMap["Processor"].ID, Value: "A17 Pro"},
				{ProductID: iphone.ID, AttributeID: attrMap["RAM"].ID, Value: "8GB"},
				{ProductID: iphone.ID, AttributeID: attrMap["Battery"].ID, Value: "3274 mAh"},
				{ProductID: iphone.ID, AttributeID: attrMap["Warranty"].ID, Value: "1 year"},
			}

			for _, val := range iphoneAttrs {
				if err := productAttrRepo.Create(val); err != nil {
					log.Printf("⏭️  Value already exists: %v", err)
				} else {
					log.Printf("  ✓ %s = %s", attrMap["Brand"].AttributeName, val.Value)
				}
			}
		}

		// Add attribute values for Samsung Galaxy S24 Ultra
		samsung, _ := productRepo.GetBySKU("SGS24U-001")
		if samsung != nil && len(attrMap) > 0 {
			log.Printf("\n--- Adding attributes for: %s ---", samsung.Name)

			samsungAttrs := []*domain.ProductAttributeValue{
				{ProductID: samsung.ID, AttributeID: attrMap["Brand"].ID, Value: "Samsung"},
				{ProductID: samsung.ID, AttributeID: attrMap["Screen Size"].ID, Value: "6.8 inch"},
				{ProductID: samsung.ID, AttributeID: attrMap["Processor"].ID, Value: "Snapdragon 8 Gen 3"},
				{ProductID: samsung.ID, AttributeID: attrMap["RAM"].ID, Value: "12GB"},
				{ProductID: samsung.ID, AttributeID: attrMap["Battery"].ID, Value: "5000 mAh"},
				{ProductID: samsung.ID, AttributeID: attrMap["Warranty"].ID, Value: "1 year"},
			}

			for _, val := range samsungAttrs {
				if err := productAttrRepo.Create(val); err != nil {
					log.Printf("⏭️  Value already exists")
				}
			}
			log.Printf("  ✅ Added attributes for Samsung")
		}

		// Add attribute values for MacBook Pro
		macbook, _ := productRepo.GetBySKU("MBP16-001")
		if macbook != nil && len(attrMap) > 0 {
			log.Printf("\n--- Adding attributes for: %s ---", macbook.Name)

			macbookAttrs := []*domain.ProductAttributeValue{
				{ProductID: macbook.ID, AttributeID: attrMap["Brand"].ID, Value: "Apple"},
				{ProductID: macbook.ID, AttributeID: attrMap["Screen Size"].ID, Value: "16 inch"},
				{ProductID: macbook.ID, AttributeID: attrMap["Processor"].ID, Value: "M3 Pro"},
				{ProductID: macbook.ID, AttributeID: attrMap["RAM"].ID, Value: "16GB"},
				{ProductID: macbook.ID, AttributeID: attrMap["Storage"].ID, Value: "512GB SSD"},
				{ProductID: macbook.ID, AttributeID: attrMap["Warranty"].ID, Value: "1 year"},
			}

			for _, val := range macbookAttrs {
				if err := productAttrRepo.Create(val); err != nil {
					log.Printf("⏭️  Value already exists")
				}
			}
			log.Printf("  ✅ Added attributes for MacBook")
		}
	}

	// ============= Clothing Category Attributes =============
	if clothing, ok := categoryMap["clothing"]; ok {
		log.Printf("\n--- Creating attributes for category: %s ---", clothing.Name)

		clothingAttrs := []*domain.CategoryAttribute{
			{
				CategoryID:    clothing.ID,
				AttributeName: "Brand",
				InputType:     "text",
				IsMandatory:   true,
				IsFilterable:  true,
			},
			{
				CategoryID:    clothing.ID,
				AttributeName: "Material",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  true,
			},
			{
				CategoryID:    clothing.ID,
				AttributeName: "Gender",
				InputType:     "select",
				IsMandatory:   false,
				IsFilterable:  true,
			},
			{
				CategoryID:    clothing.ID,
				AttributeName: "Country of Origin",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  false,
			},
		}

		attrMap := make(map[string]*domain.CategoryAttribute)
		for _, attr := range clothingAttrs {
			if err := categoryAttrRepo.Create(attr); err != nil {
				log.Printf("⏭️  Attribute %s already exists", attr.AttributeName)
				existing, _ := categoryAttrRepo.GetByCategoryID(clothing.ID)
				for _, e := range existing {
					if e.AttributeName == attr.AttributeName {
						attrMap[attr.AttributeName] = e
						break
					}
				}
			} else {
				attrMap[attr.AttributeName] = attr
				log.Printf("✅ Created attribute: %s", attr.AttributeName)
			}
		}

		// Add attribute values for Nike Air Max 90
		nike, _ := productRepo.GetBySKU("NIKE-AM90-001")
		if nike != nil && len(attrMap) > 0 {
			log.Printf("\n--- Adding attributes for: %s ---", nike.Name)

			nikeAttrs := []*domain.ProductAttributeValue{
				{ProductID: nike.ID, AttributeID: attrMap["Brand"].ID, Value: "Nike"},
				{ProductID: nike.ID, AttributeID: attrMap["Material"].ID, Value: "Leather & Mesh"},
				{ProductID: nike.ID, AttributeID: attrMap["Gender"].ID, Value: "Unisex"},
				{ProductID: nike.ID, AttributeID: attrMap["Country of Origin"].ID, Value: "Vietnam"},
			}

			for _, val := range nikeAttrs {
				if err := productAttrRepo.Create(val); err != nil {
					log.Printf("⏭️  Value already exists")
				}
			}
			log.Printf("  ✅ Added attributes for Nike")
		}

		// Add attribute values for Adidas T-Shirt
		adidas, _ := productRepo.GetBySKU("ADIDAS-TS-001")
		if adidas != nil && len(attrMap) > 0 {
			log.Printf("\n--- Adding attributes for: %s ---", adidas.Name)

			adidasAttrs := []*domain.ProductAttributeValue{
				{ProductID: adidas.ID, AttributeID: attrMap["Brand"].ID, Value: "Adidas"},
				{ProductID: adidas.ID, AttributeID: attrMap["Material"].ID, Value: "100% Cotton"},
				{ProductID: adidas.ID, AttributeID: attrMap["Gender"].ID, Value: "Unisex"},
				{ProductID: adidas.ID, AttributeID: attrMap["Country of Origin"].ID, Value: "Bangladesh"},
			}

			for _, val := range adidasAttrs {
				if err := productAttrRepo.Create(val); err != nil {
					log.Printf("⏭️  Value already exists")
				}
			}
			log.Printf("  ✅ Added attributes for Adidas")
		}
	}

	// ============= Books Category Attributes =============
	if books, ok := categoryMap["books"]; ok {
		log.Printf("\n--- Creating attributes for category: %s ---", books.Name)

		booksAttrs := []*domain.CategoryAttribute{
			{
				CategoryID:    books.ID,
				AttributeName: "Author",
				InputType:     "text",
				IsMandatory:   true,
				IsFilterable:  true,
			},
			{
				CategoryID:    books.ID,
				AttributeName: "Publisher",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  true,
			},
			{
				CategoryID:    books.ID,
				AttributeName: "Pages",
				InputType:     "number",
				IsMandatory:   false,
				IsFilterable:  false,
			},
			{
				CategoryID:    books.ID,
				AttributeName: "Language",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  true,
			},
			{
				CategoryID:    books.ID,
				AttributeName: "ISBN",
				InputType:     "text",
				IsMandatory:   false,
				IsFilterable:  false,
			},
		}

		attrMap := make(map[string]*domain.CategoryAttribute)
		for _, attr := range booksAttrs {
			if err := categoryAttrRepo.Create(attr); err != nil {
				log.Printf("⏭️  Attribute %s already exists", attr.AttributeName)
				existing, _ := categoryAttrRepo.GetByCategoryID(books.ID)
				for _, e := range existing {
					if e.AttributeName == attr.AttributeName {
						attrMap[attr.AttributeName] = e
						break
					}
				}
			} else {
				attrMap[attr.AttributeName] = attr
				log.Printf("✅ Created attribute: %s", attr.AttributeName)
			}
		}

		// Add attribute values for Clean Code
		cleanCode, _ := productRepo.GetBySKU("BOOK-CC-001")
		if cleanCode != nil && len(attrMap) > 0 {
			log.Printf("\n--- Adding attributes for: %s ---", cleanCode.Name)

			ccAttrs := []*domain.ProductAttributeValue{
				{ProductID: cleanCode.ID, AttributeID: attrMap["Author"].ID, Value: "Robert C. Martin"},
				{ProductID: cleanCode.ID, AttributeID: attrMap["Publisher"].ID, Value: "Prentice Hall"},
				{ProductID: cleanCode.ID, AttributeID: attrMap["Pages"].ID, Value: "464"},
				{ProductID: cleanCode.ID, AttributeID: attrMap["Language"].ID, Value: "English"},
				{ProductID: cleanCode.ID, AttributeID: attrMap["ISBN"].ID, Value: "978-0132350884"},
			}

			for _, val := range ccAttrs {
				if err := productAttrRepo.Create(val); err != nil {
					log.Printf("⏭️  Value already exists")
				}
			}
			log.Printf("  ✅ Added attributes for Clean Code")
		}

		// Add attribute values for DDIA
		ddia, _ := productRepo.GetBySKU("BOOK-DDIA-001")
		if ddia != nil && len(attrMap) > 0 {
			log.Printf("\n--- Adding attributes for: %s ---", ddia.Name)

			ddiaAttrs := []*domain.ProductAttributeValue{
				{ProductID: ddia.ID, AttributeID: attrMap["Author"].ID, Value: "Martin Kleppmann"},
				{ProductID: ddia.ID, AttributeID: attrMap["Publisher"].ID, Value: "O'Reilly Media"},
				{ProductID: ddia.ID, AttributeID: attrMap["Pages"].ID, Value: "616"},
				{ProductID: ddia.ID, AttributeID: attrMap["Language"].ID, Value: "English"},
				{ProductID: ddia.ID, AttributeID: attrMap["ISBN"].ID, Value: "978-1449373320"},
			}

			for _, val := range ddiaAttrs {
				if err := productAttrRepo.Create(val); err != nil {
					log.Printf("⏭️  Value already exists")
				}
			}
			log.Printf("  ✅ Added attributes for DDIA")
		}
	}

	log.Println("\n✅ Category and product attributes seeding completed!")
}
