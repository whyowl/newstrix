package sources

import (
	"context"
	"github.com/mmcdole/gofeed"
	"newstrix/internal/models"
	"time"
)

type Kommersant struct {
	parser *gofeed.Parser
}

func NewKommersant() *Kommersant {
	return &Kommersant{parser: gofeed.NewParser()}
}

func (l *Kommersant) Name() string {
	return "Kommersant.ru"
}

func (l *Kommersant) Fetch(ctx context.Context, timeline time.Time) ([]models.NewsItem, error) {

	feed, err := l.parser.ParseURLWithContext("https://www.kommersant.ru/rss/corp.xml", ctx)
	if err != nil {
		return nil, err
	}

	var items []models.NewsItem
	for _, entry := range feed.Items {
		if entry.PublishedParsed.After(timeline) || entry.PublishedParsed.Equal(timeline) {
			items = append(items, models.NewsItem{
				Title:       entry.Title,
				Link:        entry.Link,
				Description: entry.Description,
				PublishedAt: *entry.PublishedParsed,
				Publisher:   "Kommersant.ru",
			})
		}
	}

	return items, nil
}
