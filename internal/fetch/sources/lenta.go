package sources

import (
	"context"
	"github.com/mmcdole/gofeed"
	"newstrix/internal/models"
	"time"
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

func (l *Lenta) Fetch(ctx context.Context, timeline time.Time) (*[]models.NewsItem, error) {

	feed, err := l.parser.ParseURLWithContext("https://lenta.ru/rss/news", ctx)
	if err != nil {
		return nil, err
	}

	var items []models.NewsItem
	for _, entry := range feed.Items {
		if entry.PublishedParsed.After(timeline) || entry.PublishedParsed.Equal(timeline) {
			items = append(items, models.NewsItem{
				Guid:        entry.GUID,
				Title:       entry.Title,
				Link:        entry.Link,
				Description: entry.Description,
				PublishedAt: *entry.PublishedParsed,
				Publisher:   "Lenta.ru",
			})
		}
	}

	return &items, nil
}
