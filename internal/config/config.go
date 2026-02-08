package config

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddr	string
	RateLimit	int
	Port		string
}

func Load() *Config {
	cfg := &Config{
		RedisAddr: "localhost:6379",
		RateLimit: 5,
		Port: ":8080",
	}

	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.RedisAddr = addr
	}

	if ratelimit := os.Getenv("RATELIMIT"); ratelimit != "" {
		val, err := strconv.Atoi(ratelimit)
		if err == nil {
			cfg.RateLimit = val
		}
	}

	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	return cfg
}