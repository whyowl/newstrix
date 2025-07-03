package main

import (
	"context"
	"log"
	"newstrix/internal/embedding"
	"newstrix/internal/fetch"
	"newstrix/internal/fetch/sources"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()
	}()

	//cfg := config.Load() //TODO config

	srcs := []sources.Source{
		sources.NewLenta(),
		sources.NewRia(),
		sources.NewTass(),
		// sources.NewMeduza(), и т.д.
	}

	embedder, err := embedding.NewEmbedder("localhost:50051")
	if err != nil {
		log.Fatalf("Ошибка подключения к embed-сервису: %v", err)
	}

	f := fetch.NewFetcher(srcs, embedder)
	if err := f.Run(ctx); err != nil {
		panic(err)
	}
}
