package sources

import (
	"context"
	"github.com/mmcdole/gofeed"
	"newstrix/internal/models"
	"time"
)

type Tass struct {
	parser *gofeed.Parser
}

func NewTass() *Tass {
	return &Tass{parser: gofeed.NewParser()}
}

func (l *Tass) Name() string {
	return "Tass.ru"
}

func (l *Tass) Fetch(ctx context.Context, timeline time.Time) (*[]models.NewsItem, error) {

	feed, err := l.parser.ParseURLWithContext("https://tass.ru/rss/v2.xml", ctx)
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
				Publisher:   "Tass.ru",
			})
		}
	}

	return &items, nil
}
