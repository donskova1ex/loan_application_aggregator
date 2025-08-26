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

func (db *DB) Close() error {
	var errors []error

	if db.PGDB != nil {
		if sqlDB, err := db.PGDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close PGDB: %w", err))
			}
		}
	}

	if db.KassaDB != nil {
		if sqlDB, err := db.KassaDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close KassaDB: %w", err))
			}
		}
	}

	if db.DoverixDB != nil {
		if sqlDB, err := db.DoverixDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close DoverixDB: %w", err))
			}
		}
	}

	if db.DEDB != nil {
		if sqlDB, err := db.DEDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close DEDB: %w", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing databases: %v", errors)
	}

	return nil
}
