package models

import (
	"context"
	"time"
)

type Source interface {
	Name() string
	Fetch(ctx context.Context, timeline time.Time) (*[]NewsItem, error)
}

type NewsItem struct {
	Guid        string    `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Publisher   string    `json:"publisher"`
	Vector      []float32 `json:"-"`
}

type SearchParams struct {
	Keywords *[]string
	Vector   *[]float32
	Source   *string
	From     *time.Time
	To       *time.Time
	Limit    int
}
