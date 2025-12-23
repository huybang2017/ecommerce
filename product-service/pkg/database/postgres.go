package database

import (
	"fmt"
	"log"
	"product-service/config"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// dbInstance is the singleton database connection
	dbInstance *gorm.DB
	// once ensures the connection is created only once
	once sync.Once
)

// GetDB returns the singleton PostgreSQL database connection
// This implements the Singleton pattern to ensure only one DB connection pool exists
// Connection pooling is handled by GORM and the underlying driver
func GetDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var err error

	once.Do(func() {
		dsn := cfg.GetDSN()

		// Configure GORM logger (adjust based on environment)
		gormLogger := logger.Default.LogMode(logger.Silent)
		if cfg.SSLMode == "disable" { // Development mode
			gormLogger = logger.Default.LogMode(logger.Info)
		}

		// Open connection with connection pool settings
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})

		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return
		}

		// Get underlying sql.DB to configure connection pool
		sqlDB, err := dbInstance.DB()
		if err != nil {
			log.Printf("Failed to get sql.DB: %v", err)
			return
		}

		// Set connection pool parameters
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		log.Println("Database connection established successfully")
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return dbInstance, nil
}

// CloseDB closes the database connection
// This should be called during graceful shutdown
func CloseDB() error {
	if dbInstance == nil {
		return nil
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// AutoMigrate runs database migrations for all domain models
// This should be called at application startup
func AutoMigrate(db *gorm.DB) error {
	// Import domain models here to avoid circular dependencies
	// For now, we'll return nil - migrations should be handled in main.go
	return nil
}

