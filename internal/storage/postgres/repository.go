package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/pgvector/pgvector-go"
	"newstrix/internal/models"
	"strings"
	"time"
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
		values = append(values, item.Guid, item.Title, item.Link, item.Description, item.PublishedAt, item.Publisher, pgvector.NewVector(item.Vector))
	}

	query += strings.Join(placeholders, ", ")
	query += " ON CONFLICT (id) DO NOTHING"

	_, err := tx.Exec(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) UpdateSourceLastParsed(ctx context.Context, source string, lastParsed time.Time) error {

	tx := r.txManager.GetQueryEngine(ctx)

	query := "INSERT INTO source_last_parsed (source_id, last_parsed) VALUES ($1, $2) ON CONFLICT (source_id) DO UPDATE SET last_parsed = EXCLUDED.last_parsed"
	_, err := tx.Exec(ctx, query, source, lastParsed)
	if err != nil {
		return fmt.Errorf("failed to update last parsed time for source %s: %w", source, err)
	}

	return nil
}

func (r *PgRepository) GetByID(ctx context.Context, id string) (*models.NewsItem, error) {

	tx := r.txManager.GetQueryEngine(ctx)

	query := "SELECT id, title, link, description, published_at, publisher, vector FROM news WHERE id = $1"
	row := tx.QueryRow(ctx, query, id)

	var item models.NewsItem
	var v pgvector.Vector
	if err := row.Scan(&item.Guid, &item.Title, &item.Link, &item.Description, &item.PublishedAt, &item.Publisher, &v); err != nil {
		return nil, err
	}
	item.Vector = v.Slice()
	return &item, nil
}

func (r *PgRepository) SearchByFilters(ctx context.Context, opt models.SearchParams) ([]models.NewsItem, error) {
	tx := r.txManager.GetQueryEngine(ctx)

	qb := sq.Select("id", "title", "link", "description", "published_at", "publisher", "vector").
		From("news").
		Limit(uint64(opt.Limit)).
		PlaceholderFormat(sq.Dollar)

	if opt.Keywords != nil && len(*opt.Keywords) > 0 {
		for _, kw := range *opt.Keywords {
			qb = qb.Where(sq.Or{
				sq.Expr("title ILIKE ?", "%"+kw+"%"),
				sq.Expr("description ILIKE ?", "%"+kw+"%"),
			})
		}
	}

	if opt.Source != nil && *opt.Source != "" {
		qb = qb.Where(sq.Eq{"publisher": *opt.Source})
	}

	if opt.From != nil {
		qb = qb.Where(sq.GtOrEq{"published_at": *opt.From})
	}

	if opt.To != nil {
		qb = qb.Where(sq.LtOrEq{"published_at": *opt.To})
	}

	if opt.Vector != nil && len(*opt.Vector) > 0 {
		qb = qb.OrderBy(fmt.Sprintf("vector <-> '%v'", pgvector.NewVector(*opt.Vector)))
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.NewsItem
	for rows.Next() {
		var item models.NewsItem
		var v pgvector.Vector
		if err := rows.Scan(
			&item.Guid,
			&item.Title,
			&item.Link,
			&item.Description,
			&item.PublishedAt,
			&item.Publisher,
			&v,
		); err != nil {
			return nil, err
		}
		item.Vector = v.Slice()
		items = append(items, item)
	}

	return items, nil
}

func (r *PgRepository) GetSourceLastParsed(ctx context.Context, source string) (time.Time, error) {
	tx := r.txManager.GetQueryEngine(ctx)

	query := "SELECT last_parsed FROM source_last_parsed WHERE source_id = $1"
	row := tx.QueryRow(ctx, query, source)

	var lastParsed time.Time
	if err := row.Scan(&lastParsed); err != nil {
		if err == pgx.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("failed to get last parsed time for source %s: %w", source, err)
	}

	return lastParsed, nil
}
