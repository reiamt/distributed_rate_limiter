package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr string
	RateLimit int
	Port      string
	Mode      string
	FailOpen  bool
}

func Load() *Config {
	godotenv.Load()
	cfg := &Config{
		RedisAddr: "localhost:6379",
		RateLimit: 5,
		Port:      ":8080",
		Mode:      "redis",
		FailOpen:  false,
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

	if mode := os.Getenv("MODE"); mode != "" {
		cfg.Mode = mode
	}

	if failOpenTmp, exists := os.LookupEnv("FAIL_OPEN"); exists {
		failOpen, err := strconv.ParseBool(failOpenTmp)
		if err == nil {
			cfg.FailOpen = failOpen
		}
	}

	return cfg
}
