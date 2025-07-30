package repository

import (
	"app_aggregator/pkg/db"
	"gorm.io/gorm"
)

type Repository struct {
	db        *gorm.DB
	kassaDb   *gorm.DB
	doverixDb *gorm.DB
	deDb      *gorm.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{
		db:        db.PGDB,
		kassaDb:   db.KassaDB,
		doverixDb: db.DoverixDB,
		deDb:      db.DEDB,
	}
}
