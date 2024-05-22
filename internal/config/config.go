package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port string `env:"PORT" env-default:"8080"`
	Host string `env:"HOST" env-default:"localhost"`

	DBPath          string `env:"DB_PATH" env-required:"true"`
	MigrationsPath  string `env:"MIGRATIONS_PATH" env-required:"true"`
	MigrationsTable string `env:"MIGRATIONS_TABLE" env-required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return &cfg, nil
}
