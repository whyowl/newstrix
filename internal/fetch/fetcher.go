package fetch

import (
	"context"
	"fmt"
	"log"
	"newstrix/internal/embedding"
	"newstrix/internal/models"
	"newstrix/internal/storage"
	"time"
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

		lastParsed, err := f.storage.GetSourceLastParsed(ctx, source.Name())
		if err != nil {
			log.Printf("error get last parsed time for source %s: %v\n", source.Name(), err)
			continue
		}

		items, err := source.Fetch(ctx, lastParsed)
		if err != nil {
			log.Printf("Error source %s: %v\n", source.Name(), err)
			continue
		}

		for index, _ := range *items {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err = f.Vectorize(ctx, &(*items)[index]); err != nil {
					log.Printf("Error vectorize: %v\n", err) // TODO обработка ошибок
					continue
				}

			}
		}
		if len(*items) == 0 {
			log.Printf("No new items for source %s since %s\n", source.Name(), lastParsed)
			continue
		}
		if err := f.AddNews(ctx, items, source.Name(), time.Now()); err != nil {
			log.Print(err) // TODO обработка ошибок
			continue
		}
	}

	return nil
}

func (f *Fetcher) Start(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := f.Run(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (f *Fetcher) Vectorize(ctx context.Context, item *models.NewsItem) error {
	vector, err := f.embedder.Vectorize(ctx, item.Title+" "+item.Description)
	if err != nil {
		return fmt.Errorf("error vectorize item %s: %w", item.Guid, err)
	}
	item.Vector = vector
	return nil
}

func (f *Fetcher) AddNews(ctx context.Context, items *[]models.NewsItem, source string, lastParsed time.Time) error {
	if err := f.storage.AddNews(ctx, items, source, lastParsed); err != nil {
		return fmt.Errorf("error add news to storage: %w", err)
	}
	log.Printf("Added %d news items to storage", len(*items))
	return nil
}
