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
	if err := ValidatePort(c.Port); err != nil {
		return fmt.Errorf("failed to validate DB port")
	}

	return nil
}
