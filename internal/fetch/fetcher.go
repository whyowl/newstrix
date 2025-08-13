package fetch

import (
	"context"
	"fmt"
	"log"
	"newstrix/internal/embedding"
	"newstrix/internal/models"
	"newstrix/internal/storage"
)

type Fetcher struct {
	sources  []models.Source
	embedder *embedding.Embedder
	storage  storage.Facade
}

func NewFetcher(s []models.Source, e *embedding.Embedder, storage storage.Facade) *Fetcher {
	return &Fetcher{
		sources:  s,
		embedder: e,
		storage:  storage,
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

		for index, _ := range items {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err = f.Vectorize(ctx, &items[index]); err != nil {
					log.Printf("Error vectorize: %v\n", err) // TODO обработка ошибок
					continue
				}
			}
		}
		if err := f.AddNews(ctx, items); err != nil {
			log.Print(err) // TODO обработка ошибок
			continue
		}
	}

	return nil
}

func (f *Fetcher) Vectorize(ctx context.Context, item *models.NewsItem) error {
	vector, err := f.embedder.Vectorize(ctx, item.Title+" "+item.Description)
	if err != nil {
		return fmt.Errorf("error vectorize item %s: %w", item.Guid, err)
	}
	item.Vector = vector
	return nil
}

func (f *Fetcher) AddNews(ctx context.Context, items []models.NewsItem) error {
	if err := f.storage.AddNews(ctx, items); err != nil {
		return fmt.Errorf("error add news to storage: %w", err)
	}
	log.Printf("Added %d news items to storage", len(items))
	return nil
}
