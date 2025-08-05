package models

import "context"

type Source interface {
	Name() string
	Fetch(ctx context.Context) ([]NewsItem, error)
}

type NewsItem struct {
	Guid        string    `db:"id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Description string    `db:"description"`
	PublishedAt string    `db:"published_at"`
	Publisher   string    `db:"publisher"`
	Vector      []float32 `db:"vector"`
}
