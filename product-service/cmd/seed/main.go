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
	meBeID := createdCategories[3].ID
	thietBiDienTuID := createdCategories[4].ID
	nhaCuaID := createdCategories[5].ID
	sacDepID := createdCategories[6].ID
	sucKhoeID := createdCategories[7].ID
	giayNamID := createdCategories[8].ID
	giayNuID := createdCategories[9].ID
	tuiNamID := createdCategories[10].ID
	tuiNuID := createdCategories[11].ID
	dongHoID := createdCategories[12].ID
	theThaoID := createdCategories[13].ID
	xeID := createdCategories[14].ID

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
		// Electronics
		{
			ShopID:      defaultShopID,
			Name:        "iPhone 15 Pro",
			Description: "Latest iPhone with A17 Pro chip, 48MP camera, and titanium design",
			Price:       999.99,
			BasePrice:   999.99,
			SKU:         "IPH15P-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/iphone15pro.jpg"}),
			Stock:       50,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Samsung Galaxy S24 Ultra",
			Description: "Flagship Android phone with S Pen, 200MP camera, and Snapdragon 8 Gen 3",
			Price:       1199.99,
			BasePrice:   1199.99,
			SKU:         "SGS24U-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/galaxy-s24.jpg"}),
			Stock:       30,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "MacBook Pro 16-inch",
			Description: "Apple M3 Pro chip, 16GB RAM, 512GB SSD",
			Price:       2499.99,
			BasePrice:   2499.99,
			SKU:         "MBP16-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/macbook-pro.jpg"}),
			Stock:       20,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Sony WH-1000XM5 Headphones",
			Description: "Industry-leading noise canceling wireless headphones",
			Price:       399.99,
			BasePrice:   399.99,
			SKU:         "SONY-XM5-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/sony-headphones.jpg"}),
			Stock:       75,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "iPad Air",
			Description: "10.9-inch display, M2 chip, 64GB storage",
			Price:       599.99,
			BasePrice:   599.99,
			SKU:         "IPAD-AIR-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/ipad-air.jpg"}),
			Stock:       40,
			IsActive:    true,
		},
		// Clothing
		{
			ShopID:      defaultShopID,
			Name:        "Nike Air Max 90",
			Description: "Classic running shoes with Air cushioning",
			Price:       129.99,
			BasePrice:   129.99,
			SKU:         "NIKE-AM90-001",
			CategoryID:  &clothingID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/nike-shoes.jpg"}),
			Stock:       100,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Levi's 501 Jeans",
			Description: "Original fit straight leg jeans",
			Price:       89.99,
			BasePrice:   89.99,
			SKU:         "LEVI-501-001",
			CategoryID:  &clothingID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/levis-jeans.jpg"}),
			Stock:       80,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Adidas Originals T-Shirt",
			Description: "Classic three-stripe design cotton t-shirt",
			Price:       29.99,
			BasePrice:   29.99,
			SKU:         "ADIDAS-TS-001",
			CategoryID:  &clothingID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/adidas-tshirt.jpg"}),
			Stock:       150,
			IsActive:    true,
		},
		// Books
		{
			ShopID:      defaultShopID,
			Name:        "The Clean Code",
			Description: "A Handbook of Agile Software Craftsmanship by Robert C. Martin",
			Price:       49.99,
			BasePrice:   49.99,
			SKU:         "BOOK-CC-001",
			CategoryID:  &booksID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/clean-code.jpg"}),
			Stock:       200,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Designing Data-Intensive Applications",
			Description: "The Big Ideas Behind Reliable, Scalable, and Maintainable Systems",
			Price:       59.99,
			BasePrice:   59.99,
			SKU:         "BOOK-DDIA-001",
			CategoryID:  &booksID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/ddia.jpg"}),
			Stock:       150,
			IsActive:    true,
		},
		// Home & Kitchen
		{
			ShopID:      defaultShopID,
			Name:        "KitchenAid Stand Mixer",
			Description: "5-quart stand mixer with 10 speeds",
			Price:       379.99,
			BasePrice:   379.99,
			SKU:         "KA-MIXER-001",
			CategoryID:  &homeID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/kitchenaid.jpg"}),
			Stock:       25,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Instant Pot Duo",
			Description: "7-in-1 electric pressure cooker, slow cooker, rice cooker",
			Price:       99.99,
			BasePrice:   99.99,
			SKU:         "IP-DUO-001",
			CategoryID:  &homeID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/instant-pot.jpg"}),
			Stock:       60,
			IsActive:    true,
		},
		// Sports & Outdoors
		{
			ShopID:      defaultShopID,
			Name:        "Yoga Mat Premium",
			Description: "Non-slip exercise mat, 6mm thick, eco-friendly",
			Price:       39.99,
			BasePrice:   39.99,
			SKU:         "YOGA-MAT-001",
			CategoryID:  &sportsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/yoga-mat.jpg"}),
			Stock:       120,
			IsActive:    true,
		},
		{
			ShopID:      defaultShopID,
			Name:        "Dumbbell Set 20kg",
			Description: "Adjustable dumbbells with stand",
			Price:       199.99,
			BasePrice:   199.99,
			SKU:         "DUMB-20KG-001",
			CategoryID:  &sportsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/dumbbells.jpg"}),
			Stock:       35,
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
	iphone, _ := productRepo.GetBySKU("IPH15P-001")
	nikeshoes, _ := productRepo.GetBySKU("NIKE-AM90-001")
	tshirt, _ := productRepo.GetBySKU("ADIDAS-TS-001")

	if iphone == nil || nikeshoes == nil || tshirt == nil {
		log.Println("⚠️  Required products not found, skipping product items seeding")
		return
	}

	// ============= iPhone 15 Pro - Storage + Color variations =============
	log.Printf("\n--- Creating variations for: %s ---", iphone.Name)

	// Create Variations
	storageVar := &domain.Variation{ProductID: iphone.ID, Name: "Storage"}
	variationRepo.Create(storageVar)

	colorVar := &domain.Variation{ProductID: iphone.ID, Name: "Color"}
	variationRepo.Create(colorVar)

	// Create Variation Options - Storage
	storage128 := &domain.VariationOption{VariationID: storageVar.ID, Value: "128GB"}
	storage256 := &domain.VariationOption{VariationID: storageVar.ID, Value: "256GB"}
	storage512 := &domain.VariationOption{VariationID: storageVar.ID, Value: "512GB"}
	variationOptRepo.Create(storage128)
	variationOptRepo.Create(storage256)
	variationOptRepo.Create(storage512)

	// Create Variation Options - Color
	colorNatural := &domain.VariationOption{VariationID: colorVar.ID, Value: "Natural Titanium"}
	colorBlue := &domain.VariationOption{VariationID: colorVar.ID, Value: "Blue Titanium"}
	colorBlack := &domain.VariationOption{VariationID: colorVar.ID, Value: "Black Titanium"}
	variationOptRepo.Create(colorNatural)
	variationOptRepo.Create(colorBlue)
	variationOptRepo.Create(colorBlack)

	// Create Product Items (SKUs) - combinations
	iphoneSKUs := []struct {
		storage *domain.VariationOption
		color   *domain.VariationOption
		sku     string
		price   float64
		stock   int
	}{
		{storage128, colorNatural, "IPH15P-128-NAT", 999.99, 20},
		{storage128, colorBlue, "IPH15P-128-BLU", 999.99, 15},
		{storage128, colorBlack, "IPH15P-128-BLK", 999.99, 18},
		{storage256, colorNatural, "IPH15P-256-NAT", 1099.99, 12},
		{storage256, colorBlue, "IPH15P-256-BLU", 1099.99, 10},
		{storage512, colorNatural, "IPH15P-512-NAT", 1299.99, 8},
		{storage512, colorBlack, "IPH15P-512-BLK", 1299.99, 5},
	}

	for _, item := range iphoneSKUs {
		productItem := &domain.ProductItem{
			ProductID:  iphone.ID,
			SKUCode:    item.sku,
			ImageURL:   "https://example.com/iphone15pro.jpg",
			Price:      item.price,
			QtyInStock: item.stock,
			Status:     "ACTIVE",
		}

		if err := productItemRepo.Create(productItem); err != nil {
			log.Printf("⏭️  SKU %s already exists", item.sku)
			continue
		}

		// Link with variation options
		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.storage.ID,
		})
		skuConfigRepo.Create(&domain.SKUConfiguration{
			ProductItemID:     productItem.ID,
			VariationOptionID: item.color.ID,
		})

		log.Printf("✅ Created SKU: %s - %s %s ($%.2f, stock: %d)",
			item.sku, item.storage.Value, item.color.Value, item.price, item.stock)
	}

	// ============= Nike Air Max 90 - Size + Color variations =============
	log.Printf("\n--- Creating variations for: %s ---", nikeshoes.Name)

	sizeVar := &domain.Variation{ProductID: nikeshoes.ID, Name: "Size"}
	variationRepo.Create(sizeVar)

	shoeColorVar := &domain.Variation{ProductID: nikeshoes.ID, Name: "Color"}
	variationRepo.Create(shoeColorVar)

	// Sizes
	size8 := &domain.VariationOption{VariationID: sizeVar.ID, Value: "US 8"}
	size9 := &domain.VariationOption{VariationID: sizeVar.ID, Value: "US 9"}
	size10 := &domain.VariationOption{VariationID: sizeVar.ID, Value: "US 10"}
	size11 := &domain.VariationOption{VariationID: sizeVar.ID, Value: "US 11"}
	variationOptRepo.Create(size8)
	variationOptRepo.Create(size9)
	variationOptRepo.Create(size10)
	variationOptRepo.Create(size11)

	// Colors
	colorWhite := &domain.VariationOption{VariationID: shoeColorVar.ID, Value: "White"}
	colorRed := &domain.VariationOption{VariationID: shoeColorVar.ID, Value: "Red"}
	colorNavy := &domain.VariationOption{VariationID: shoeColorVar.ID, Value: "Navy"}
	variationOptRepo.Create(colorWhite)
	variationOptRepo.Create(colorRed)
	variationOptRepo.Create(colorNavy)

	nikeSKUs := []struct {
		size  *domain.VariationOption
		color *domain.VariationOption
		sku   string
		stock int
	}{
		{size8, colorWhite, "NIKE-AM90-8-WHT", 15},
		{size9, colorWhite, "NIKE-AM90-9-WHT", 20},
		{size9, colorRed, "NIKE-AM90-9-RED", 18},
		{size10, colorWhite, "NIKE-AM90-10-WHT", 25},
		{size10, colorRed, "NIKE-AM90-10-RED", 22},
		{size10, colorNavy, "NIKE-AM90-10-NAV", 20},
		{size11, colorWhite, "NIKE-AM90-11-WHT", 12},
		{size11, colorNavy, "NIKE-AM90-11-NAV", 10},
	}

	for _, item := range nikeSKUs {
		productItem := &domain.ProductItem{
			ProductID:  nikeshoes.ID,
			SKUCode:    item.sku,
			ImageURL:   "https://example.com/nike-shoes.jpg",
			Price:      129.99, // Same price for all variants
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

		log.Printf("✅ Created SKU: %s - %s %s (stock: %d)",
			item.sku, item.size.Value, item.color.Value, item.stock)
	}

	// ============= Adidas T-Shirt - Size + Color variations =============
	log.Printf("\n--- Creating variations for: %s ---", tshirt.Name)

	tshirtSizeVar := &domain.Variation{ProductID: tshirt.ID, Name: "Size"}
	variationRepo.Create(tshirtSizeVar)

	tshirtColorVar := &domain.Variation{ProductID: tshirt.ID, Name: "Color"}
	variationRepo.Create(tshirtColorVar)

	// Sizes
	sizeS := &domain.VariationOption{VariationID: tshirtSizeVar.ID, Value: "S"}
	sizeM := &domain.VariationOption{VariationID: tshirtSizeVar.ID, Value: "M"}
	sizeL := &domain.VariationOption{VariationID: tshirtSizeVar.ID, Value: "L"}
	sizeXL := &domain.VariationOption{VariationID: tshirtSizeVar.ID, Value: "XL"}
	variationOptRepo.Create(sizeS)
	variationOptRepo.Create(sizeM)
	variationOptRepo.Create(sizeL)
	variationOptRepo.Create(sizeXL)

	// Colors
	tBlack := &domain.VariationOption{VariationID: tshirtColorVar.ID, Value: "Black"}
	tWhite := &domain.VariationOption{VariationID: tshirtColorVar.ID, Value: "White"}
	tGray := &domain.VariationOption{VariationID: tshirtColorVar.ID, Value: "Gray"}
	variationOptRepo.Create(tBlack)
	variationOptRepo.Create(tWhite)
	variationOptRepo.Create(tGray)

	tshirtSKUs := []struct {
		size  *domain.VariationOption
		color *domain.VariationOption
		sku   string
		stock int
	}{
		{sizeS, tBlack, "ADIDAS-TS-S-BLK", 30},
		{sizeS, tWhite, "ADIDAS-TS-S-WHT", 25},
		{sizeM, tBlack, "ADIDAS-TS-M-BLK", 40},
		{sizeM, tWhite, "ADIDAS-TS-M-WHT", 35},
		{sizeM, tGray, "ADIDAS-TS-M-GRY", 30},
		{sizeL, tBlack, "ADIDAS-TS-L-BLK", 35},
		{sizeL, tWhite, "ADIDAS-TS-L-WHT", 30},
		{sizeL, tGray, "ADIDAS-TS-L-GRY", 25},
		{sizeXL, tBlack, "ADIDAS-TS-XL-BLK", 20},
		{sizeXL, tWhite, "ADIDAS-TS-XL-WHT", 18},
	}

	for _, item := range tshirtSKUs {
		productItem := &domain.ProductItem{
			ProductID:  tshirt.ID,
			SKUCode:    item.sku,
			ImageURL:   "https://example.com/adidas-tshirt.jpg",
			Price:      29.99,
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

		log.Printf("✅ Created SKU: %s - Size %s Color %s (stock: %d)",
			item.sku, item.size.Value, item.color.Value, item.stock)
	}

	log.Println("\n✅ Product items seeding completed!")

	// ============= Seed specific ID 101 for testing cart =============
	log.Println("\n--- Creating test product item ID 101 ---")

	// Check if ID 101 already exists
	existing, _ := productItemRepo.GetByID(101)
	if existing != nil {
		log.Println("⏭️  Product item ID 101 already exists, skipping")
	} else {
		// Create a simple product item with ID 101 (using Nike shoes product)
		testItem := &domain.ProductItem{
			ID:         101,
			ProductID:  nikeshoes.ID,
			SKUCode:    "NIKE-AM90-TEST-101",
			ImageURL:   "https://example.com/nike-test-101.jpg",
			Price:      119.99,
			QtyInStock: 100,
			Status:     "active",
		}

		if err := productItemRepo.Create(testItem); err != nil {
			log.Printf("❌ Failed to create product item ID 101: %v", err)
		} else {
			log.Printf("✅ Created product item ID 101: %s (Price: $%.2f, Stock: %d)",
				testItem.SKUCode, testItem.Price, testItem.QtyInStock)
		}
	}
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
