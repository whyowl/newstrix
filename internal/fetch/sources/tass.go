package sources

import (
	"context"
	"github.com/mmcdole/gofeed"
	"newstrix/internal/models"
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

func (l *Tass) Fetch(ctx context.Context) ([]models.NewsItem, error) {

	feed, err := l.parser.ParseURLWithContext("https://tass.ru/rss/v2.xml", ctx)
	if err != nil {
		return nil, err
	}

	var items []models.NewsItem
	for _, entry := range feed.Items {
		items = append(items, models.NewsItem{
			Guid:        entry.GUID,
			Title:       entry.Title,
			Link:        entry.Link,
			Description: entry.Description,
			PublishedAt: parseTime(entry.Published),
			Publisher:   "Tass.ru",
		})
	}

	return items, nil
}
