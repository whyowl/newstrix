package sources

import "context"

type Source interface {
	Name() string
	Fetch(ctx context.Context) ([]NewsItem, error)
}

type NewsItem struct {
	Guid        string
	Title       string
	Link        string
	Description string
	PublishedAt string
	Publisher   string
}
