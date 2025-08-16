package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"newstrix/internal/config"
	"newstrix/internal/embedding"
	"newstrix/internal/fetch"
	"newstrix/internal/fetch/sources"
	"newstrix/internal/models"
	"newstrix/internal/storage"
	"newstrix/internal/storage/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	srcs := []models.Source{
		sources.NewLenta(),
		sources.NewRia(),
		sources.NewTass(),
		// TODO sources.NewMeduza(), и т.д.
	}

	embedder, err := embedding.NewEmbedder(cfg.EmbedderURL)
	if err != nil {
		log.Fatalf("error connect to embed-service: %v", err) // TODO test try
	}

	pool, err := pgxpool.Connect(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storageFacade := newStorageFacade(pool)

	f := fetch.NewFetcher(srcs, embedder, storageFacade, cfg.MaxWorkers)

	go func() {
		log.Printf("Starting Fetcher with interval %s...", cfg.FetchInterval)
		err := f.Start(ctx, cfg.FetchInterval)
		if err != nil {
			log.Fatalf("Fetcher failed: %v", err)
		} else {
			log.Println("Fetcher stopped gracefully")
		}
	}()

	<-ctx.Done()
	log.Println("Received termination signal, shutting down...")
	time.Sleep(10 * time.Second) // Give some time for fetcher to finish current tasks

}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
