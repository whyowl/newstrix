package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"newstrix/internal/api"
	"newstrix/internal/config"
	"newstrix/internal/embedding"
	"newstrix/internal/search"
	"newstrix/internal/storage"
	"newstrix/internal/storage/postgres"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()
	}()

	cfg := config.Load()

	pool, err := pgxpool.Connect(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storageFacade := newStorageFacade(pool)

	embedder, err := embedding.NewEmbedder(cfg.EmbedderURL)
	if err != nil {
		log.Fatalf("error connect to embed-service: %v", err)
	}

	//todo unit tests

	searchEngine := search.NewSearchEngine(ctx, embedder, storageFacade)

	router := api.SetupRouter(searchEngine)

	log.Printf("Starting API server at %s...", cfg.ApiAddress)
	err = router.Run(cfg.ApiAddress)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
