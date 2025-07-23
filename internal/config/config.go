package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type PGConfig struct {
	DSN string
}

type Config struct {
	PGdb PGConfig
}

func InitConfig() (*Config, error) {
	err := godotenv.Load(".env.local")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		PGdb: PGConfig{
			DSN: os.Getenv("POSTGRES_DSN"),
		},
	}
	return config, nil
}
