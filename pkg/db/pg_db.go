package db

import (
	"app_aggregator/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DB struct {
	PGDB      *gorm.DB
	KassaDB   *gorm.DB
	DoverixDB *gorm.DB
	DEDB      *gorm.DB
}

func InitDB(cfg *config.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.PGdb.DSN))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	kassaDb, err := gorm.Open(sqlserver.Open(cfg.SQL.DsnKassa))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kassa database: %w", err)
	}
	doverixDb, err := gorm.Open(sqlserver.Open(cfg.SQL.DsnDoverix))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to doverix database: %w", err)
	}
	deDb, err := gorm.Open(sqlserver.Open(cfg.SQL.DsnDe))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to dened database: %w", err)
	}

	return &DB{
		PGDB:      db,
		KassaDB:   kassaDb,
		DoverixDB: doverixDb,
		DEDB:      deDb,
	}, nil
}
