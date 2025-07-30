package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type PGConfig struct {
	DSN string
}

type SQLConfig struct {
	DsnKassa   string
	DsnDoverix string
	DsnDe      string
}

type Config struct {
	PGdb PGConfig
	SQL  SQLConfig
}

func InitConfig() (*Config, error) {
	err := godotenv.Load(".env.local")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	sqlServer := os.Getenv("SQL_SERVER")
	sqlUser := os.Getenv("SQL_USER")
	sqlPassword := os.Getenv("SQL_PASSWORD")
	sqlKassaDb := os.Getenv("SQL_DB_KASSA")
	sqlDoverixDb := os.Getenv("SQL_DB_DOVERIX")
	sqlDeDb := os.Getenv("SQL_DB_DE")

	kassaDSN := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;", sqlServer, sqlUser, sqlPassword, sqlKassaDb)
	doverixDSN := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;", sqlServer, sqlUser, sqlPassword, sqlDoverixDb)
	deDSN := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;", sqlServer, sqlUser, sqlPassword, sqlDeDb)
	config := &Config{
		PGdb: PGConfig{
			DSN: os.Getenv("POSTGRES_DSN"),
		},
		SQL: SQLConfig{
			DsnKassa:   kassaDSN,
			DsnDoverix: doverixDSN,
			DsnDe:      deDSN,
		},
	}
	return config, nil
}
