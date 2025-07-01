package fetch

import (
	"context"
	"fmt"
	"log"
	"newstrix/internal/fetch/sources"
)

type Fetcher struct {
	sources []sources.Source
}

func NewFetcher(s []sources.Source) *Fetcher {
	return &Fetcher{sources: s}
}

func (f *Fetcher) Run(ctx context.Context) error {
	for _, source := range f.sources {
		log.Printf("Parse %s...\n", source.Name())
		items, err := source.Fetch(ctx)
		if err != nil {
			log.Printf("Error source %s: %v\n", source.Name(), err)
			continue
		}

		for _, item := range items {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fmt.Printf("→ %s [%s]\n", item.Title, item.Publisher)
				// TODO: отправить в БД / векторизатор
			}
		}
	}

	return nil
}
