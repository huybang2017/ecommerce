package database

import (
	"fmt"
	"order-service/config"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB returns a singleton database connection
// This ensures we only have one connection pool per service
func GetDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dsn := cfg.GetDSN()
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return
		}

		// Get underlying sql.DB to configure connection pool
		sqlDB, err2 := db.DB()
		if err2 != nil {
			err = fmt.Errorf("failed to get underlying sql.DB: %w", err2)
			return
		}

		// Set connection pool settings
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		// Test connection
		if err2 = sqlDB.Ping(); err2 != nil {
			err = fmt.Errorf("failed to ping database: %w", err2)
			return
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

