package config

import (
	"fmt"
	"net"
	"net/url"
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
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", c.Port)
	}

	return nil
}
