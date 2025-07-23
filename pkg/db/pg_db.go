package db

import (
	"app_aggregator/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func InitDB(cfg *config.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.PGdb.DSN))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &DB{db}, nil
}
