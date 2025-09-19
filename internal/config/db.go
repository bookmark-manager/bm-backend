package config

import (
	"fmt"
	"net"
	"net/url"
	"slices"
	"strconv"
)

type DBConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Port     int    `env:"PORT" env-default:"5432"`
	Name     string `env:"NAME" env-default:"bookmarks"`
	SSLMode  string `env:"SSL_MODE" env-default:"disable"`
}

func (c *DBConfig) DSN() url.URL {
	return url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:     c.Name,
		RawQuery: fmt.Sprintf("sslmode=%s", c.SSLMode),
	}
}

func (c *DBConfig) Validate() error {
	if err := ValidatePort(c.Port); err != nil {
		return fmt.Errorf("failed to validate DB port: %w", err)
	}

	if err := ValidateSLLMode(c.SSLMode); err != nil {
		return fmt.Errorf("failed to validate SSLMode: %w", err)
	}

	return nil
}

func ValidateSLLMode(sslMode string) error {
	options := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}

	if !slices.Contains(options, sslMode) {
		return fmt.Errorf("invalid sslMode value: %s", sslMode)
	}

	return nil
}
