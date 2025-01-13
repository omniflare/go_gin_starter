package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	PORT       string `env:"SERVER_PORT,required"`
	Database   string `env:"DATABASE_URL,required"`
	ResendAPI  string `env:"RESEND_API_KEY,required"`
	JWT_SECRET string `env:"JWT_SECRET,required"`
	CLIENT_URL string `env:"CLIENT_URL,required"`
}

func NewEnvConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error while accessing env file: %s", err.Error())
	}
	config := &Config{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Error parsing env config %e", err)
	}
	return config
}
