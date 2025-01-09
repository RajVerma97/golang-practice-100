package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ScrapedData struct {
	URL     string
	Content string
	Error   error
}

type Scraper struct {
	client      *http.Client
	rateLimiter chan struct{}
	maxRetries  int
	timeout     time.Duration
}

func NewScraper(maxConcurrent int, maxRetries int, timeout time.Duration) *Scraper {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{

			MaxIdleConns:        100,
			MaxConnsPerHost:     100,
			DisableCompression:  false,
			DisableKeepAlives:   false,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 10,
		},
	}

	return &Scraper{
		client:      client,
		rateLimiter: make(chan struct{}, maxConcurrent),
		maxRetries:  maxRetries,
		timeout:     timeout,
	}
}

func (s *Scraper) fetchWithRetry(ctx context.Context, url string) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Apply rate limiting
		select {
		case s.rateLimiter <- struct{}{}:
			defer func() { <-s.rateLimiter }()
		case <-ctx.Done():
			return "", ctx.Err()
		}

		// Create request with context
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		// Add headers to make the request more reliable
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("status code: %d", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		return string(body), nil
	}

	return "", fmt.Errorf("failed after %d retries: %w", s.maxRetries, lastErr)
}
func (s *Scraper) ScrapeURLs(urls []string) (<-chan ScrapedData, error) {
	if len(urls) == 0 {
		return nil, errors.New("no URLs provided")
	}

	resultChan := make(chan ScrapedData, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Create a context with timeout for each URL
			ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
			defer cancel()

			content, err := s.fetchWithRetry(ctx, url)
			result := ScrapedData{
				URL:     url,
				Content: content,
				Error:   err,
			}

			select {
			case resultChan <- result:
			case <-ctx.Done():
				return
			}
		}(url)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan, nil
}

func main() {
	start := time.Now()

	urls := []string{
		"https://www.bbc.com/news",
		"https://twitter.com/",
		"https://reddit.com/",
	}

	scraper := NewScraper(
		3,              // Max concurrent requests
		3,              // Max retries
		10*time.Second, // Timeout per request
	)

	results, err := scraper.ScrapeURLs(urls)
	if err != nil {
		fmt.Printf("Failed to initialize scraping: %v\n", err)
		return
	}

	successCount := 0
	failCount := 0

	for result := range results {
		if result.Error != nil {
			failCount++
			fmt.Printf("Failed to scrape %s: %v\n", result.URL, result.Error)
			continue
		}

		successCount++
		fmt.Printf("\nSuccessfully scraped %s (content length: %d)\n",
			result.URL, len(result.Content))
		fmt.Println(result.Content[:100])

	}

	duration := time.Since(start)
	fmt.Printf("\nScraping Summary:\n")
	fmt.Printf("Total URLs: %d\n", len(urls))
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Printf("Total time: %v\n", duration)
}
