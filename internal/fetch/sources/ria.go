package sources

import (
	"context"
	"github.com/mmcdole/gofeed"
	"newstrix/internal/models"
	"time"
)

type Ria struct {
	parser *gofeed.Parser
}

func NewRia() *Ria {
	return &Ria{parser: gofeed.NewParser()}
}

func (l *Ria) Name() string {
	return "Ria.ru"
}

func (l *Ria) Fetch(ctx context.Context, timeline time.Time) (*[]models.NewsItem, error) {

	feed, err := l.parser.ParseURLWithContext("https://ria.ru/export/rss2/archive/index.xml", ctx)
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
				Publisher:   "Ria.ru",
			})
		}
	}

	return &items, nil
}
