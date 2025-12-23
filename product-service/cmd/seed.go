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

	log.Println("Starting to seed data...")

	// 1. Create Categories
	log.Println("Creating categories...")
	categories := []*domain.Category{
		{
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Electronic devices and gadgets",
		},
		{
			Name:        "Clothing",
			Slug:        "clothing",
			Description: "Fashion and apparel",
		},
		{
			Name:        "Books",
			Slug:        "books",
			Description: "Books and literature",
		},
		{
			Name:        "Home & Kitchen",
			Slug:        "home-kitchen",
			Description: "Home and kitchen appliances",
		},
		{
			Name:        "Sports & Outdoors",
			Slug:        "sports-outdoors",
			Description: "Sports equipment and outdoor gear",
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
	electronicsID := createdCategories[0].ID
	clothingID := createdCategories[1].ID
	booksID := createdCategories[2].ID
	homeID := createdCategories[3].ID
	sportsID := createdCategories[4].ID

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

	products := []*domain.Product{
		// Electronics
		{
			Name:        "iPhone 15 Pro",
			Description: "Latest iPhone with A17 Pro chip, 48MP camera, and titanium design",
			Price:       999.99,
			SKU:         "IPH15P-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/iphone15pro.jpg"}),
			Stock:       50,
			IsActive:    true,
		},
		{
			Name:        "Samsung Galaxy S24 Ultra",
			Description: "Flagship Android phone with S Pen, 200MP camera, and Snapdragon 8 Gen 3",
			Price:       1199.99,
			SKU:         "SGS24U-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/galaxy-s24.jpg"}),
			Stock:       30,
			IsActive:    true,
		},
		{
			Name:        "MacBook Pro 16-inch",
			Description: "Apple M3 Pro chip, 16GB RAM, 512GB SSD",
			Price:       2499.99,
			SKU:         "MBP16-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/macbook-pro.jpg"}),
			Stock:       20,
			IsActive:    true,
		},
		{
			Name:        "Sony WH-1000XM5 Headphones",
			Description: "Industry-leading noise canceling wireless headphones",
			Price:       399.99,
			SKU:         "SONY-XM5-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/sony-headphones.jpg"}),
			Stock:       75,
			IsActive:    true,
		},
		{
			Name:        "iPad Air",
			Description: "10.9-inch display, M2 chip, 64GB storage",
			Price:       599.99,
			SKU:         "IPAD-AIR-001",
			CategoryID:  &electronicsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/ipad-air.jpg"}),
			Stock:       40,
			IsActive:    true,
		},
		// Clothing
		{
			Name:        "Nike Air Max 90",
			Description: "Classic running shoes with Air cushioning",
			Price:       129.99,
			SKU:         "NIKE-AM90-001",
			CategoryID:  &clothingID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/nike-shoes.jpg"}),
			Stock:       100,
			IsActive:    true,
		},
		{
			Name:        "Levi's 501 Jeans",
			Description: "Original fit straight leg jeans",
			Price:       89.99,
			SKU:         "LEVI-501-001",
			CategoryID:  &clothingID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/levis-jeans.jpg"}),
			Stock:       80,
			IsActive:    true,
		},
		{
			Name:        "Adidas Originals T-Shirt",
			Description: "Classic three-stripe design cotton t-shirt",
			Price:       29.99,
			SKU:         "ADIDAS-TS-001",
			CategoryID:  &clothingID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/adidas-tshirt.jpg"}),
			Stock:       150,
			IsActive:    true,
		},
		// Books
		{
			Name:        "The Clean Code",
			Description: "A Handbook of Agile Software Craftsmanship by Robert C. Martin",
			Price:       49.99,
			SKU:         "BOOK-CC-001",
			CategoryID:  &booksID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/clean-code.jpg"}),
			Stock:       200,
			IsActive:    true,
		},
		{
			Name:        "Designing Data-Intensive Applications",
			Description: "The Big Ideas Behind Reliable, Scalable, and Maintainable Systems",
			Price:       59.99,
			SKU:         "BOOK-DDIA-001",
			CategoryID:  &booksID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/ddia.jpg"}),
			Stock:       150,
			IsActive:    true,
		},
		// Home & Kitchen
		{
			Name:        "KitchenAid Stand Mixer",
			Description: "5-quart stand mixer with 10 speeds",
			Price:       379.99,
			SKU:         "KA-MIXER-001",
			CategoryID:  &homeID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/kitchenaid.jpg"}),
			Stock:       25,
			IsActive:    true,
		},
		{
			Name:        "Instant Pot Duo",
			Description: "7-in-1 electric pressure cooker, slow cooker, rice cooker",
			Price:       99.99,
			SKU:         "IP-DUO-001",
			CategoryID:  &homeID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/instant-pot.jpg"}),
			Stock:       60,
			IsActive:    true,
		},
		// Sports & Outdoors
		{
			Name:        "Yoga Mat Premium",
			Description: "Non-slip exercise mat, 6mm thick, eco-friendly",
			Price:       39.99,
			SKU:         "YOGA-MAT-001",
			CategoryID:  &sportsID,
			Status:      "ACTIVE",
			Images:      createImagesJSON([]string{"https://example.com/yoga-mat.jpg"}),
			Stock:       120,
			IsActive:    true,
		},
		{
			Name:        "Dumbbell Set 20kg",
			Description: "Adjustable dumbbells with stand",
			Price:       199.99,
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
	log.Println("\n✅ Data seeding finished!")
}

