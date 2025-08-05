package postgres

import (
	"context"
	"fmt"
	"newstrix/internal/models"
	"strings"
)

type PgRepository struct {
	txManager TransactionManager
}

func NewPgRepository(txManager TransactionManager) *PgRepository {
	return &PgRepository{txManager: txManager}
}

func (r *PgRepository) AddNews(ctx context.Context, news []models.NewsItem) error {

	tx := r.txManager.GetQueryEngine(ctx)

	query := "INSERT INTO news (id, title, link, description, published_at, publisher, vector) VALUES "
	values := []interface{}{}
	placeholders := []string{}

	for i, item := range news {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
		values = append(values, item.Guid, item.Title, item.Link, item.Description, item.PublishedAt, item.Publisher, item.Vector)
	}

	query += strings.Join(placeholders, ", ")
	query += " ON CONFLICT (id) DO NOTHING"

	_, err := tx.Exec(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) GetNewsByVector(ctx context.Context, vector []float32) (*[]models.NewsItem, error) {
	// TODO write GetNewsByVector()
	return nil, nil
}

// TODO add another Get functions
