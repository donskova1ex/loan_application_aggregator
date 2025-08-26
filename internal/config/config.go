package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
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

func buildMSSQLDSN(server, user, password, database string) string {
	if server == "" || user == "" || password == "" || database == "" {
		return ""
	}
	return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable;", server, user, password, database)
}

func InitConfig() (*Config, error) {

	_ = godotenv.Load(".env.local", ".env")

	sqlServer := os.Getenv("SQL_SERVER")
	sqlUser := os.Getenv("SQL_USER")
	sqlPassword := os.Getenv("SQL_PASSWORD")
	sqlKassaDb := os.Getenv("SQL_DB_KASSA")
	sqlDoverixDb := os.Getenv("SQL_DB_DOVERIX")
	sqlDeDb := os.Getenv("SQL_DB_DE")

	kassaDSN := buildMSSQLDSN(sqlServer, sqlUser, sqlPassword, sqlKassaDb)
	doverixDSN := buildMSSQLDSN(sqlServer, sqlUser, sqlPassword, sqlDoverixDb)
	deDSN := buildMSSQLDSN(sqlServer, sqlUser, sqlPassword, sqlDeDb)

	pgDSN := os.Getenv("POSTGRES_DSN")
	if pgDSN == "" {
		return nil, errors.New("POSTGRES_DSN is required but not set")
	}
	config := &Config{
		PGdb: PGConfig{
			DSN: pgDSN,
		},
		SQL: SQLConfig{
			DsnKassa:   kassaDSN,
			DsnDoverix: doverixDSN,
			DsnDe:      deDSN,
		},
	}
	return config, nil
}
