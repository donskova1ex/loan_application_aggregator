package migrations

import (
	"app_aggregator/internal/config"
	"app_aggregator/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Up(cfg *config.Config) error {
	db, err := gorm.Open(postgres.Open(cfg.PGdb.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed creating extension \"uuid-ossp\": %w", err)
	}

	err = initMigrations(db)
	if err != nil {
		return fmt.Errorf("failed initialising migrations: %w", err)
	}
	return nil
}

func initMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Organization{})
	if err != nil {
		err := db.Migrator().DropTable(&models.Organization{})
		if err != nil {
			return fmt.Errorf("failed dropping table organizations: %w", err)
		}
		return fmt.Errorf("failed creating table organizations: %w", err)
	}

	err = db.AutoMigrate(&models.LoanApplication{})
	if err != nil {
		err := db.Migrator().DropTable(&models.LoanApplication{})
		if err != nil {
			return fmt.Errorf("failed dropping table loan_applications: %w", err)
		}
		return fmt.Errorf("failed creating table loan_applications: %w", err)
	}

	err = db.AutoMigrate(&models.Settings{})
	if err != nil {
		err := db.Migrator().DropTable(&models.Settings{})
		if err != nil {
			return fmt.Errorf("failed dropping table settings: %w", err)
		}
		return fmt.Errorf("failed creating table settings: %w", err)
	}
	return nil
}
