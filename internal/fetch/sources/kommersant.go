package sources

import (
	"context"
	"github.com/mmcdole/gofeed"
	"newstrix/internal/models"
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

func (l *Kommersant) Fetch(ctx context.Context) ([]models.NewsItem, error) {

	feed, err := l.parser.ParseURLWithContext("https://www.kommersant.ru/rss/corp.xml", ctx)
	if err != nil {
		return nil, err
	}

	var items []models.NewsItem
	for _, entry := range feed.Items {
		items = append(items, models.NewsItem{
			Title:       entry.Title,
			Link:        entry.Link,
			Description: entry.Description,
			PublishedAt: entry.Published,
			Publisher:   "Kommersant.ru",
		})
	}

	return items, nil
}
