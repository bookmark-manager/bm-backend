package config

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type HttpConfig struct {
	Host        string        `env:"HOST" env-default:"0.0.0.0"`
	Port        int           `env:"PORT" env-default:"8080"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
}

func (c *HttpConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func (c *HttpConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", c.Port)
	}

	return nil
}
