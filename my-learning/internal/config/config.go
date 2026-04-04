package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DbUrl string
}

func LoadConfig() *Config {
	// Load .env file from the project root (two levels up from internal/config/)
	godotenv.Load("../../.env")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	dbUrl := os.Getenv("DB_URL")

	if port == "" || dbUrl == "" {
		log.Fatal("PORT and DB_URL must be set in the .env file")
	}

	return &Config{
		Port:  port,
		DbUrl: dbUrl,
	}
}