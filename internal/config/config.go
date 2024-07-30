package config

import "log"

type Config struct {
	DatabaseURL string `yaml:"database_url" required:"true"`
	JWTSecret   string `yaml:"jwt_secret" required:"true"`
	ServerPort  string `yaml:"server_port" required:"true"`
}

func Load() *Config {
	path := "./config/local.yaml"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
