package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConfig
	HttpConfig
}

type DBConfig struct {
	DBHost     string `env:"BM_DB_HOST" env-default:"localhost"`
	DBUser     string `env:"BM_DB_USER" env-required:"true"`
	DBPassword string `env:"BM_DB_PASSWORD" env-required:"true"`
	DBPort     int    `env:"BM_DB_PORT" env-default:"5432"`
	DBName     string `env:"BM_DB_NAME" env-default:"bookmarks"`
}

type HttpConfig struct {
	HttpHost        string        `env:"BM_HTTP_HOST" env-default:"0.0.0.0"`
	HttpPort        int           `env:"BM_HTTP_PORT" env-default:"8080"`
	HttpIdleTimeout time.Duration `env:"BM_HTTP_IDLE_TIMEOUT" env-default:"60s"`
	HttpTimeout     time.Duration `env:"BM_HTTP_TIMEOUT" env-default:"4s"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.HttpHost, c.HttpPort)
}

func validateConfig(cfg *Config) error {
	if err := validatePort(cfg.DBPort); err != nil {
		return fmt.Errorf("invalid database port (BM_DB_PORT): %w", err)
	}
	if err := validatePort(cfg.HttpPort); err != nil {
		return fmt.Errorf("invalid http server port (BM_HTTP_PORT): %w", err)
	}

	return nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", port)
	}

	return nil
}
