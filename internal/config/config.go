package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

func LoadEnv() *Config {
	// Load from .env file if it exists, ignore error in production environments
	_ = godotenv.Load()

	return &Config{
		Port: os.Getenv("PORT"),
		Dsn: os.Getenv("DSN"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

}