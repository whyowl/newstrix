package main

import (
	"log"
	"newstrix/internal/api"
	"newstrix/internal/config"
)

func main() {
	cfg := config.Load()

	router := api.SetupRouter(cfg)

	log.Printf("Starting API server at %s...", cfg.HTTPPort)
	err := router.Run(":" + cfg.HTTPPort)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
