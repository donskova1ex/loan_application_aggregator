package config

type PGConfig struct {
	DSN string
}

type Config struct {
	PGdb PGConfig
}
