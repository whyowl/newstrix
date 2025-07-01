package sources

import (
	"context"
	"github.com/mmcdole/gofeed"
)

type Lenta struct {
	parser *gofeed.Parser
}

func NewLenta() *Lenta {
	return &Lenta{parser: gofeed.NewParser()}
}

func (l *Lenta) Name() string {
	return "Lenta.ru"
}

func (l *Lenta) Fetch(ctx context.Context) ([]NewsItem, error) {

	feed, err := l.parser.ParseURLWithContext("https://lenta.ru/rss/news", ctx)
	if err != nil {
		return nil, err
	}

	var items []NewsItem
	for _, entry := range feed.Items {
		items = append(items, NewsItem{
			Title:       entry.Title,
			Link:        entry.Link,
			Description: entry.Description,
			PublishedAt: entry.Published,
			Publisher:   "Lenta.ru",
		})
	}

	return items, nil
}
