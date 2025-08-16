package fetch

import (
	"context"
	"fmt"
	"log"
	"newstrix/internal/embedding"
	"newstrix/internal/models"
	"newstrix/internal/storage"
	"sync"
	"time"
)

type Fetcher struct {
	sources    []models.Source
	embedder   *embedding.Embedder
	storage    storage.Facade
	maxWorkers int
	stats      *FetchStats
}

type FetchStats struct {
	mu              sync.RWMutex
	TotalSources    int
	SuccessfulFetch int
	FailedFetch     int
	TotalItems      int
	VectorizedItems int
	FailedItems     int
	LastRunTime     time.Time
}

type FetchResult struct {
	Source string
	Items  []models.NewsItem
	Error  error
}

type ProcessedItem struct {
	Item   models.NewsItem
	Source string
	Error  error
}

func NewFetcher(s []models.Source, e *embedding.Embedder, storage storage.Facade, maxWorkers int) *Fetcher {
	return &Fetcher{
		sources:    s,
		embedder:   e,
		storage:    storage,
		maxWorkers: maxWorkers,
		stats:      &FetchStats{TotalSources: len(s)},
	}
}

func (f *Fetcher) Run(ctx context.Context) error {
	startTime := time.Now()
	log.Printf("Starting concurrent fetch for %d sources with %d workers", len(f.sources), f.maxWorkers)

	// Reset stats for this run
	f.stats.mu.Lock()
	f.stats.LastRunTime = startTime
	f.stats.SuccessfulFetch = 0
	f.stats.FailedFetch = 0
	f.stats.TotalItems = 0
	f.stats.VectorizedItems = 0
	f.stats.FailedItems = 0
	f.stats.mu.Unlock()

	// Create channels for coordination
	sourceResults := make(chan FetchResult, len(f.sources))

	// Start source fetchers concurrently
	var wg sync.WaitGroup
	for _, source := range f.sources {
		wg.Add(1)
		go f.fetchSource(ctx, source, sourceResults, &wg)
	}

	// Close source results channel when all fetchers complete
	go func() {
		wg.Wait()
		close(sourceResults)
	}()

	// Process and store results
	err := f.processAndStore(ctx, sourceResults)

	// Log final stats
	duration := time.Since(startTime)
	f.stats.mu.RLock()
	log.Printf("Fetch completed in %v. Sources: %d/%d, Items: %d/%d, Vectorized: %d/%d",
		duration,
		f.stats.SuccessfulFetch, f.stats.TotalSources,
		f.stats.VectorizedItems, f.stats.TotalItems,
		f.stats.VectorizedItems, f.stats.TotalItems)
	f.stats.mu.RUnlock()

	return err
}

func (f *Fetcher) fetchSource(ctx context.Context, source models.Source, results chan<- FetchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Fetching from source: %s", source.Name())

	lastParsed, err := f.storage.GetSourceLastParsed(ctx, source.Name())
	if err != nil {
		log.Printf("Error getting last parsed time for source %s: %v", source.Name(), err)
		f.stats.mu.Lock()
		f.stats.FailedFetch++
		f.stats.mu.Unlock()
		results <- FetchResult{Source: source.Name(), Error: err}
		return
	}

	items, err := source.Fetch(ctx, lastParsed)
	if err != nil {
		log.Printf("Error fetching from source %s: %v", source.Name(), err)
		f.stats.mu.Lock()
		f.stats.FailedFetch++
		f.stats.mu.Unlock()
		results <- FetchResult{Source: source.Name(), Error: err}
		return
	}

	if len(*items) == 0 {
		log.Printf("No new items for source %s since %s", source.Name(), lastParsed)
		f.stats.mu.Lock()
		f.stats.SuccessfulFetch++
		f.stats.mu.Unlock()
		results <- FetchResult{Source: source.Name(), Items: []models.NewsItem{}}
		return
	}

	log.Printf("Fetched %d items from source %s", len(*items), source.Name())
	f.stats.mu.Lock()
	f.stats.SuccessfulFetch++
	f.stats.TotalItems += len(*items)
	f.stats.mu.Unlock()
	results <- FetchResult{Source: source.Name(), Items: *items}
}

func (f *Fetcher) processAndStore(ctx context.Context, sourceResults <-chan FetchResult) error {
	// Create a map to batch items by source
	sourceBatches := make(map[string][]models.NewsItem)
	var mu sync.Mutex

	// Process source results and collect items by source
	for result := range sourceResults {
		if result.Error != nil {
			log.Printf("Source %s failed: %v", result.Source, result.Error)
			continue
		}

		if len(result.Items) == 0 {
			continue
		}

		// Vectorize items for this source
		vectorizedItems, err := f.vectorizeItems(ctx, result.Items)
		if err != nil {
			log.Printf("Failed to vectorize items for source %s: %v", result.Source, err)
			f.stats.mu.Lock()
			f.stats.FailedItems += len(result.Items)
			f.stats.mu.Unlock()
			continue
		}

		// Add vectorized items to batch
		mu.Lock()
		sourceBatches[result.Source] = append(sourceBatches[result.Source], vectorizedItems...)
		mu.Unlock()

		f.stats.mu.Lock()
		f.stats.VectorizedItems += len(vectorizedItems)
		f.stats.mu.Unlock()
	}

	// Store all batches
	return f.storeBatches(ctx, sourceBatches)
}

func (f *Fetcher) vectorizeItems(ctx context.Context, items []models.NewsItem) ([]models.NewsItem, error) {
	var vectorizedItems []models.NewsItem
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Process items concurrently with limited workers
	semaphore := make(chan struct{}, f.maxWorkers)

	for i := range items {
		wg.Add(1)
		go func(item models.NewsItem) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Vectorize item
			if err := f.VectorizeWithRetry(ctx, &item, 3); err != nil {
				log.Printf("Failed to vectorize item %s: %v", item.Guid, err)
				f.stats.mu.Lock()
				f.stats.FailedItems++
				f.stats.mu.Unlock()
				return
			}

			// Add to results
			mu.Lock()
			vectorizedItems = append(vectorizedItems, item)
			mu.Unlock()
		}(items[i])
	}

	wg.Wait()
	return vectorizedItems, nil
}

func (f *Fetcher) storeBatches(ctx context.Context, sourceBatches map[string][]models.NewsItem) error {
	var wg sync.WaitGroup
	var errors []error
	var mu sync.Mutex

	// Store each source batch concurrently
	for source, items := range sourceBatches {
		if len(items) == 0 {
			continue
		}

		wg.Add(1)
		go func(sourceName string, sourceItems []models.NewsItem) {
			defer wg.Done()

			log.Printf("Storing %d items from source %s", len(sourceItems), sourceName)

			if err := f.storage.AddNews(ctx, &sourceItems, sourceName, time.Now()); err != nil {
				log.Printf("Failed to store %d items from source %s: %v", len(sourceItems), sourceName, err)
				mu.Lock()
				errors = append(errors, fmt.Errorf("source %s: %w", sourceName, err))
				mu.Unlock()
				return
			}

			log.Printf("Successfully stored %d items from source %s", len(sourceItems), sourceName)
		}(source, items)
	}

	wg.Wait()

	// Return first error if any occurred
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

func (f *Fetcher) Start(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := f.Run(ctx); err != nil {
				log.Printf("Fetch run failed: %v", err)
			}
		case <-ctx.Done():
			log.Println("Fetcher shutting down...")
			return nil
		}
	}
}

func (f *Fetcher) Vectorize(ctx context.Context, item *models.NewsItem) error {
	vector, err := f.embedder.Vectorize(ctx, item.Title+" "+item.Description)
	if err != nil {
		return fmt.Errorf("failed to vectorize item %s: %w", item.Guid, err)
	}
	item.Vector = vector
	return nil
}

func (f *Fetcher) VectorizeWithRetry(ctx context.Context, item *models.NewsItem, maxRetries int) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := f.Vectorize(ctx, item); err != nil {
			lastErr = err
			if isRetryableError(err) && attempt < maxRetries {
				backoff := time.Duration(attempt*attempt) * time.Second
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
					continue
				}
			}
			return err
		}
		return nil
	}
	return fmt.Errorf("vectorization failed after %d attempts: %w", maxRetries, lastErr)
}

func (f *Fetcher) AddNews(ctx context.Context, items *[]models.NewsItem, source string, lastParsed time.Time) error {
	if err := f.storage.AddNews(ctx, items, source, lastParsed); err != nil {
		return fmt.Errorf("failed to add news to storage: %w", err)
	}
	log.Printf("Added %d news items to storage from source %s", len(*items), source)
	return nil
}

func isRetryableError(err error) bool {
	// todo Add logic to determine if error is transient
	// For now, assume network/timeout errors are retryable
	return err != nil
}

// GetStats returns a copy of the current fetch statistics
func (f *Fetcher) GetStats() FetchStats {
	f.stats.mu.RLock()
	defer f.stats.mu.RUnlock()
	return *f.stats
}
