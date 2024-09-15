package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	// PostgresConn     string `env:"POSTGRES_CONN"`
	// PostgresJdbcUrl  string `env:"POSTGRES_JDBC_URL"`
	PostgresUsername string `env:"POSTGRES_USERNAME"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDatabase string `env:"POSTGRES_DATABASE"`
}

func New() (*Config, error) {
	var err error

	err = godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Config parsed")

	return cfg, nil
}
