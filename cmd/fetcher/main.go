package main

import (
	"context"
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

	// Отлов сигнала Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()
	}()

	//cfg := config.Load()

	srcs := []sources.Source{
		sources.NewLenta(),
		//sources.NewInterfax(),
		// sources.NewMeduza(), и т.д.
	}

	f := fetch.NewFetcher(srcs)
	if err := f.Run(ctx); err != nil {
		panic(err)
	}
}
