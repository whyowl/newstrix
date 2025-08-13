package storage

import (
	"context"
	"newstrix/internal/models"
	"newstrix/internal/storage/postgres"
)

type Facade interface {
	AddNews(ctx context.Context, news []models.NewsItem) error
	GetByID(ctx context.Context, id string) (*models.NewsItem, error)
	SearchByFilters(ctx context.Context, opt models.SearchParams) ([]models.NewsItem, error)
}

type StorageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository *postgres.PgRepository
}

func NewStorageFacade(txManager postgres.TransactionManager, pgRepository *postgres.PgRepository) Facade {
	return &StorageFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
	}
}

func (f *StorageFacade) AddNews(ctx context.Context, news []models.NewsItem) error {
	return f.pgRepository.AddNews(ctx, news)
}

func (f *StorageFacade) GetByID(ctx context.Context, id string) (*models.NewsItem, error) {
	return f.pgRepository.GetByID(ctx, id)
}

func (f *StorageFacade) SearchByFilters(ctx context.Context, opt models.SearchParams) ([]models.NewsItem, error) {
	return f.pgRepository.SearchByFilters(ctx, opt)
}
