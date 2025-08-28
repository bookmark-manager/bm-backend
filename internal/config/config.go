package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Port     string `env:"PORT" env-default:"5432"`
	Name     string `env:"DB_NAME"`
}

func MustLoad() *DBConfig {
	var cfg DBConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	return &cfg
}
