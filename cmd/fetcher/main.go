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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	f := fetch.NewFetcher(srcs, embedder, storageFacade)
	if err := f.Run(ctx); err != nil {
		panic(err)
	}
	<-ctx.Done()
}

func newStorageFacade(pool *pgxpool.Pool) storage.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)

	return storage.NewStorageFacade(txManager, pgRepository)
}
