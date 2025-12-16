package database

import (
	"database/sql"
	"io"
	"os"
	"time"

	"github.com/rzolkos/web-recap/internal/models"
	_ "modernc.org/sqlite"
)

// ChromeHandler handles Chrome/Chromium/Edge browser history
type ChromeHandler struct {
	dbPath string
}

// NewChromeHandler creates a new Chrome history handler
func NewChromeHandler(dbPath string) *ChromeHandler {
	return &ChromeHandler{
		dbPath: dbPath,
	}
}

// GetHistory retrieves history entries from Chrome
func (h *ChromeHandler) GetHistory(startDate, endDate time.Time) ([]models.HistoryEntry, error) {
	// Copy database to temp location to avoid locking issues
	tempDB, err := h.copyDatabase()
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempDB)

	db, err := sql.Open("sqlite", tempDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare date filters
	// Query the visits table joined with urls to get individual visit records
	// (not just last_visit_time per URL)
	var query string
	var args []interface{}

	if !startDate.IsZero() || !endDate.IsZero() {
		query = `
		SELECT
			v.visit_time,
			u.url,
			u.title,
			u.visit_count
		FROM visits v
		JOIN urls u ON v.url = u.id
		WHERE v.visit_time > 0
		`

		if !startDate.IsZero() {
			chromeStart := (startDate.Unix() + 11644473600) * 1000000
			query += ` AND v.visit_time >= ?`
			args = append(args, chromeStart)
		}

		if !endDate.IsZero() {
			// Only add 24 hours if the end time is at midnight (user specified just a date)
			endTimestamp := endDate.Unix()
			if endDate.Hour() == 0 && endDate.Minute() == 0 && endDate.Second() == 0 {
				endTimestamp += 86400
			}
			chromeEnd := (endTimestamp + 11644473600) * 1000000
			query += ` AND v.visit_time < ?`
			args = append(args, chromeEnd)
		}

		query += ` ORDER BY v.visit_time DESC`
	} else {
		query = `
		SELECT
			v.visit_time,
			u.url,
			u.title,
			u.visit_count
		FROM visits v
		JOIN urls u ON v.url = u.id
		WHERE v.visit_time > 0
		ORDER BY v.visit_time DESC
		LIMIT 10000
		`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.HistoryEntry

	for rows.Next() {
		var chromeTime int64
		var url, title string
		var visitCount int

		if err := rows.Scan(&chromeTime, &url, &title, &visitCount); err != nil {
			continue
		}

		timestamp := ConvertChromeTimestamp(chromeTime)
		if timestamp.IsZero() {
			continue
		}

		entries = append(entries, models.HistoryEntry{
			Timestamp:  timestamp,
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			Domain:     ExtractDomain(url),
			Browser:    "chrome",
		})
	}

	return entries, rows.Err()
}

// copyDatabase copies the Chrome database to a temporary file
func (h *ChromeHandler) copyDatabase() (string, error) {
	src, err := os.Open(h.dbPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.CreateTemp("", "web-recap-chrome-*.db")
	if err != nil {
		return "", err
	}
	tmpFile := dst.Name()
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(tmpFile)
		return "", err
	}

	return tmpFile, nil
}
