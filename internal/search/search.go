package search

import (
	"context"
	"fmt"
	"newstrix/internal/models"
	"time"
)

const (
	DefaultLimit     = 20
	MaxLimit         = 100
	MaxQueryLength   = 300
	MaxSourceLength  = 100
	DefaultDateRange = 6 * time.Hour
)

type QueryOption struct {
	Query *string
	models.SearchParams
}

type SearchRepository interface {
	GetByID(ctx context.Context, id string) (*models.NewsItem, error)
	SearchByFilters(ctx context.Context, opt models.SearchParams) ([]models.NewsItem, error)
}

type Vectorizer interface {
	Vectorize(ctx context.Context, text string) ([]float32, error)
}

type SearchEngine struct {
	ctx      context.Context
	embedder Vectorizer
	storage  SearchRepository
}

func NewSearchEngine(ctx context.Context, embedder Vectorizer, storage SearchRepository) *SearchEngine {
	return &SearchEngine{
		ctx:      ctx,
		embedder: embedder,
		storage:  storage,
	}
}

func (s *SearchEngine) GetByID(ctx context.Context, id *string) (*models.NewsItem, error) {
	if id == nil || *id == "" {
		return nil, nil
	}
	item, err := s.storage.GetByID(ctx, *id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *SearchEngine) SearchByKeywords(ctx context.Context, keywords []string) ([]models.NewsItem, error) {
	if len(keywords) == 0 {
		return nil, fmt.Errorf("keywords cannot be empty")
	}
	items, err := s.storage.SearchByFilters(ctx, models.SearchParams{
		Keywords: &keywords,
		Limit:    DefaultLimit,
	})

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *SearchEngine) SearchBySemanticQuery(ctx context.Context, query string, limit int) ([]models.NewsItem, error) {
	if query == "" {
		return nil, fmt.Errorf("invalid query: query='%s'", query)
	}
	if len(query) > MaxQueryLength {
		query = query[:MaxQueryLength]
	}
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	vec, err := s.embedder.Vectorize(ctx, query)
	if err != nil {
		return nil, err
	}

	items, err := s.storage.SearchByFilters(ctx, models.SearchParams{
		Vector: &vec,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *SearchEngine) SearchAdvanced(ctx context.Context, params QueryOption) ([]models.NewsItem, error) {

	if params.Query == nil && params.Source == nil && params.From == nil && params.To == nil && params.Keywords == nil {
		return nil, fmt.Errorf("at least one search parameter must be provided")
	}

	if params.Source != nil {
		if len(*params.Source) > MaxSourceLength {
			return nil, fmt.Errorf("source name too long")
		}
	}

	if params.From != nil && params.To != nil && params.From.After(*params.To) {
		return nil, fmt.Errorf("invalid date range: from='%s', to='%s'", params.From, params.To)
	}

	if params.Limit <= 0 {
		params.Limit = DefaultLimit
	}
	if params.Limit > MaxLimit {
		params.Limit = MaxLimit
	}

	if params.From == nil && params.To != nil {
		from := params.To.Add(-DefaultDateRange)
		params.From = &from
	}
	if params.From != nil && params.To == nil {
		to := time.Now()
		params.To = &to
	}

	if params.Query != nil {
		vec, err := s.embedder.Vectorize(ctx, *params.Query)
		if err != nil {
			return nil, fmt.Errorf("error vectorizing query: %w", err)
		}
		params.Vector = &vec
	}

	request := models.SearchParams{
		Keywords: params.Keywords,
		Vector:   params.Vector,
		Source:   params.Source,
		From:     params.From,
		To:       params.To,
		Limit:    params.Limit,
	}

	return s.storage.SearchByFilters(ctx, request)
}
