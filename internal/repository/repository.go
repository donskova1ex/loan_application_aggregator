package repository

import (
	"app_aggregator/pkg/db"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{
		db: db.DB,
	}
}
