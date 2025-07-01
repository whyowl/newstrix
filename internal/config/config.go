package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
}

func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	var AppConfig Config

	log.Println("Config loaded")
	return AppConfig
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
