package models

import (
	"context"
	"time"
)

type Source interface {
	Name() string
	Fetch(ctx context.Context) ([]NewsItem, error)
}

type NewsItem struct {
	Guid        string    `db:"id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Description string    `db:"description"`
	PublishedAt time.Time `db:"published_at"`
	Publisher   string    `db:"publisher"`
	Vector      []float32 `db:"vector"`
}

type SearchParams struct {
	Keywords *[]string
	Vector   *[]float32
	Source   *string
	From     *time.Time
	To       *time.Time
	Limit    int
}
