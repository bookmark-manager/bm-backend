package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBHost     string `env:"DB_HOST" env-default:"localhost"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBPort     int    `env:"DB_PORT" env-default:"5432"`
	DBName     string `env:"DB_NAME" env-default:"bookmarks"`

	HttpAddress     string        `env:"HTTP_ADDRESS" env-default:"0.0.0.0:8080"`
	HttpIdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	HttpTimeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"4s"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	return &cfg
}
