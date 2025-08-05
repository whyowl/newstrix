package storage

import (
	"context"
	"newstrix/internal/models"
	"newstrix/internal/storage/postgres"
)

type Facade interface {
	AddNews(ctx context.Context, news []models.NewsItem) error
}

type StorageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository *postgres.PgRepository
}

func (f *StorageFacade) AddNews(ctx context.Context, news []models.NewsItem) error {
	return f.pgRepository.AddNews(ctx, news)
}

func NewStorageFacade(txManager postgres.TransactionManager, pgRepository *postgres.PgRepository) Facade {
	return &StorageFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
	}
}
