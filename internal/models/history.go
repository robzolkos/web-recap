package models

import "time"

// HistoryEntry represents a single browser history entry
type HistoryEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitCount int       `json:"visit_count"`
	Domain     string    `json:"domain"`
	Browser    string    `json:"browser"`
}

// HistoryReport represents a collection of history entries for a specific time period
type HistoryReport struct {
	Browser      string          `json:"browser"`
	StartDate    time.Time       `json:"start_date"`
	EndDate      time.Time       `json:"end_date"`
	Timezone     string          `json:"timezone"`
	TotalEntries int             `json:"total_entries"`
	Entries      []HistoryEntry  `json:"entries"`
}

// BrowserType represents the type of browser
type BrowserType string

const (
	BrowserChrome    BrowserType = "chrome"
	BrowserChromium  BrowserType = "chromium"
	BrowserEdge      BrowserType = "edge"
	BrowserFirefox   BrowserType = "firefox"
	BrowserSafari    BrowserType = "safari"
	BrowserUnknown   BrowserType = "unknown"
)

func (b BrowserType) String() string {
	return string(b)
}
