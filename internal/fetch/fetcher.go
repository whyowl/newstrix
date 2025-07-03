package fetch

import (
	"context"
	"fmt"
	"log"
	"newstrix/internal/embedding"
	"newstrix/internal/fetch/sources"
)

type Fetcher struct {
	sources  []sources.Source
	embedder *embedding.Embedder
}

func NewFetcher(s []sources.Source, e *embedding.Embedder) *Fetcher {
	return &Fetcher{
		sources:  s,
		embedder: e,
	}
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
				//fmt.Printf("→ %s [%s]\n", item.Title, item.Publisher)
				vector, err := f.embedder.Vectorize(ctx, item.Title+" "+item.Description) //+ " " + item.FullText)
				if err != nil {
					fmt.Printf("Error vectorize: %v\n", err) // TODO обработка ошибок
					continue
				}
				fmt.Println(vector)
				// TODO: отправить в БД
			}
		}
	}

	return nil
}
