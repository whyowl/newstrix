package storage

import (
	"context"
	"newstrix/internal/models"
	"newstrix/internal/storage/postgres"
	"time"
)

type Facade interface {
	AddNews(ctx context.Context, news *[]models.NewsItem, source string, updateAt time.Time) error
	GetByID(ctx context.Context, id string) (*models.NewsItem, error)
	SearchByFilters(ctx context.Context, opt models.SearchParams) ([]models.NewsItem, error)
	GetSourceLastParsed(ctx context.Context, source string) (time.Time, error)
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

func (f *StorageFacade) AddNews(ctx context.Context, news *[]models.NewsItem, source string, updateAt time.Time) error {
	return f.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		if err := f.pgRepository.AddNews(ctxTx, *news); err != nil {
			return err
		}

		if err := f.pgRepository.UpdateSourceLastParsed(ctxTx, source, updateAt); err != nil {
			return err
		}

		return nil
	})
}

func (f *StorageFacade) GetByID(ctx context.Context, id string) (*models.NewsItem, error) {
	return f.pgRepository.GetByID(ctx, id)
}

func (f *StorageFacade) SearchByFilters(ctx context.Context, opt models.SearchParams) ([]models.NewsItem, error) {
	return f.pgRepository.SearchByFilters(ctx, opt)
}

func (f *StorageFacade) GetSourceLastParsed(ctx context.Context, source string) (time.Time, error) {
	return f.pgRepository.GetSourceLastParsed(ctx, source)
}
