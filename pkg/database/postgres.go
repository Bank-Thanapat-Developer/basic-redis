package database

import (
	"fmt"

	"github.com/Bank-Thanapat-Developer/basic-redis/config"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg config.PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.RefItemType{},
		&entities.Item{},
	)
}
