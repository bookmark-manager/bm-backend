package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB     DBConfig     `env-prefix:"BM_DB_"`
	Http   HttpConfig   `env-prefix:"BM_HTTP_"`
	Logger LoggerConfig `env-prefix:"BM_LOGGER"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Port     int    `env:"PORT" env-default:"5432"`
	Name     string `env:"NAME" env-default:"bookmarks"`
}

type HttpConfig struct {
	Host        string        `env:"HOST" env-default:"0.0.0.0"`
	Port        int           `env:"PORT" env-default:"8080"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
}

type LoggerConfig struct {
	NoColor bool `env:"NO_COLOR" env-default:"false"`
	Debug   bool `env:"DEBUG" env-default:"true"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := cfg.validateConfig(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) DSN() url.URL {
	return url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DB.User, c.DB.Password),
		Host:     net.JoinHostPort(c.DB.Host, strconv.Itoa(c.DB.Port)),
		Path:     c.DB.Name,
		RawQuery: "sslmode=disable",
	}
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Http.Host, strconv.Itoa(c.Http.Port))
}

func (c *Config) validateConfig() error {
	if err := validatePort(c.DB.Port); err != nil {
		return fmt.Errorf("invalid database port (BM_DB_PORT): %w", err)
	}
	if err := validatePort(c.Http.Port); err != nil {
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
