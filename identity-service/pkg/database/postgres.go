package database

import (
	"fmt"
	"identity-service/config"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// GetDB returns the singleton PostgreSQL database connection
func GetDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var err error

	once.Do(func() {
		dsn := cfg.GetDSN()

		gormLogger := logger.Default.LogMode(logger.Silent)
		if cfg.SSLMode == "disable" {
			gormLogger = logger.Default.LogMode(logger.Info)
		}

		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})

		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return
		}

		sqlDB, err := dbInstance.DB()
		if err != nil {
			log.Printf("Failed to get sql.DB: %v", err)
			return
		}

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


