package config

import (
	"fmt"
	"url-shortener/internal/domain/models"

	"github.com/caarlos0/env"
)

type Config struct {
	HTTP_port int         `env:"HTTP_PORT" envDefault:"8080"`
	IsProd    bool        `env:"IS_PROD" envDefault:"false"`
	MongoAddr string      `env:"MONGO_ADDR" envDefault:"mongodb://localhost:27017"`
	RedisAddr string      `env:"REDIS_ADDR" envDefault:"127.0.0.1:6379"`
	Mode      models.Mode `env:"MODE" envDefault:"cached"`
}

var config Config = Config{}

func GetConfig() (*Config, error) {
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("read logger configuration failed: %w", err)
	}
	return &config, nil
}
