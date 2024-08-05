package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	DBName     string `env:"DB_NAME" required:"true"`
	DBUser     string `env:"DB_USER" required:"true"`
	DBPassword string `env:"DB_PASSWORD" required:"true"`
	DBHost     string `env:"DB_HOST" required:"true"`
	DBPort     string `env:"DB_PORT" required:"true"`

	RedisPort string `env:"REDIS_PORT" required:"true"`

	JWTSecret             string `env:"JWT_SECRET" required:"true"`
	PostServiceServerPort string `env:"POST_SERVICE_SERVER_PORT" required:"true"`
	UserServiceAddress    string `env:"USER_SERVICE_ADDRESS" required:"true"`
}

func Load() *Config {
	path := "./.env"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
