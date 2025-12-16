package database

import (
	"sort"
	"time"

	"github.com/rzolkos/web-recap/internal/browser"
	"github.com/rzolkos/web-recap/internal/models"
)

// HistoryQuerier defines the interface for querying browser history
type HistoryQuerier interface {
	GetHistory(startDate, endDate time.Time) ([]models.HistoryEntry, error)
}

// NewQuerier creates a new history querier for the given browser
func NewQuerier(b *browser.Browser) (HistoryQuerier, error) {
	switch b.Type {
	case browser.Chrome, browser.Chromium, browser.Edge, browser.Brave:
		return NewChromeHandler(b.Path), nil
	case browser.Firefox:
		return NewFirefoxHandler(b.Path), nil
	case browser.Safari:
		return NewSafariHandler(b.Path), nil
	default:
		return nil, ErrUnsupportedBrowser
	}
}

// Query retrieves history entries from a specific browser
func Query(b *browser.Browser, startDate, endDate time.Time) ([]models.HistoryEntry, error) {
	querier, err := NewQuerier(b)
	if err != nil {
		return nil, err
	}

	entries, err := querier.GetHistory(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Sort by timestamp descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	return entries, nil
}

// QueryMultipleBrowsers retrieves history from all detected browsers
func QueryMultipleBrowsers(detector *browser.Detector, startDate, endDate time.Time) ([]models.HistoryEntry, error) {
	var allEntries []models.HistoryEntry

	detectedBrowsers := detector.Detect()
	for _, b := range detectedBrowsers {
		browser := b // Copy to avoid pointer issues
		entries, err := Query(&browser, startDate, endDate)
		if err != nil {
			// Log error but continue with other browsers
			continue
		}
		allEntries = append(allEntries, entries...)
	}

	// Sort all entries by timestamp descending
	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].Timestamp.After(allEntries[j].Timestamp)
	})

	return allEntries, nil
}
