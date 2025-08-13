package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	EmbedderURL string
	PostgresURL string
	GrpcAddress string
	OllamaURL   string
	OllamaModel string
	ApiAddress  string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	var AppConfig Config

	AppConfig = Config{
		EmbedderURL: getEnv("EMBEDDER_URL", "localhost:50051"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://news:password@localhost:5432/newsdb?sslmode=disable"),
		GrpcAddress: getEnv("GRPC_ADDRESS", ":50051"),
		OllamaURL:   getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel: getEnv("OLLAMA_MODEL", "bge-m3:latest"),
		ApiAddress:  getEnv("API_ADDRESS", ":8080"),
	}

	log.Println("Config loaded")
	return &AppConfig
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
