package fetch

import (
	"context"
	"newstrix/internal/embedding"
	"newstrix/internal/models"
	"newstrix/internal/storage"
	"testing"
	"time"
)

func TestFetcherBatching(t *testing.T) {
	// This is a basic test to ensure the fetcher structure is correct
	// In a real implementation, you'd want to mock the dependencies

	t.Run("NewFetcher creates fetcher with correct configuration", func(t *testing.T) {
		sources := []models.Source{} // Empty for test
		var embedder *embedding.Embedder
		var facade storage.Facade

		fetcher := NewFetcher(sources, embedder, facade, 5)

		if fetcher.maxWorkers != 5 {
			t.Errorf("Expected maxWorkers to be 5, got %d", fetcher.maxWorkers)
		}

		if fetcher.stats.TotalSources != 0 {
			t.Errorf("Expected TotalSources to be 0, got %d", fetcher.stats.TotalSources)
		}
	})

	t.Run("Fetcher stats are properly initialized", func(t *testing.T) {
		sources := []models.Source{} // Empty for test
		var embedder *embedding.Embedder
		var facade storage.Facade

		fetcher := NewFetcher(sources, embedder, facade, 10)
		stats := fetcher.GetStats()

		if stats.TotalSources != 0 {
			t.Errorf("Expected TotalSources to be 0, got %d", stats.TotalSources)
		}

		if stats.SuccessfulFetch != 0 {
			t.Errorf("Expected SuccessfulFetch to be 0, got %d", stats.SuccessfulFetch)
		}

		if stats.FailedFetch != 0 {
			t.Errorf("Expected FailedFetch to be 0, got %d", stats.FailedFetch)
		}
	})
}

// Mock implementations for testing (you'd typically use a mocking library)
type MockSource struct {
	name string
}

func (m *MockSource) Name() string {
	return m.name
}

func (m *MockSource) Fetch(ctx context.Context, timeline time.Time) (*[]models.NewsItem, error) {
	// Mock implementation
	return &[]models.NewsItem{}, nil
}
