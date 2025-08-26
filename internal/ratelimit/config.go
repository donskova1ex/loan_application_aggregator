package ratelimit

import "time"

type Config struct {
	RequestsPerMinute int
	WindowSize        time.Duration
	CleanupInterval   time.Duration
	BlockDuration     time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		RequestsPerMinute: 100,
		WindowSize:        1 * time.Minute,
		CleanupInterval:   1 * time.Minute,
		BlockDuration:     1 * time.Minute,
	}
}
