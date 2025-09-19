package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB      DBConfig   `env-prefix:"BM_DB_"`
	HTTP    HttpConfig `env-prefix:"BM_HTTP_"`
	NoColor bool       `env:"BM_NO_COLOR" env-default:"false"`
	Debug   bool       `env:"BM_DEBUG" env-default:"true"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("db validation failed: %w", err)
	}

	if err := c.HTTP.Validate(); err != nil {
		return fmt.Errorf("http validation failed: %w", err)
	}

	return nil
}

func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", port)
	}

	return nil
}
