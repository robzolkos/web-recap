package database

import (
	"database/sql"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/rzolkos/web-recap/internal/models"
	_ "modernc.org/sqlite"
)

// SafariHandler handles Safari browser history (macOS only)
type SafariHandler struct {
	dbPath string
}

// NewSafariHandler creates a new Safari history handler
func NewSafariHandler(dbPath string) *SafariHandler {
	return &SafariHandler{
		dbPath: dbPath,
	}
}

// GetHistory retrieves history entries from Safari
func (h *SafariHandler) GetHistory(startDate, endDate time.Time) ([]models.HistoryEntry, error) {
	// Safari is only available on macOS
	if runtime.GOOS != "darwin" {
		return nil, ErrSafariNotAvailable
	}

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
	// Query history_visits joined with history_items to get individual visit records
	// (not just the last visit per URL)
	var query string
	var args []interface{}

	if !startDate.IsZero() || !endDate.IsZero() {
		query = `
		SELECT
			hv.visit_time,
			hi.url,
			COALESCE(hv.title, hi.url) as title,
			hi.visit_count
		FROM history_visits hv
		JOIN history_items hi ON hv.history_item = hi.id
		WHERE hv.visit_time > 0
		`

		if !startDate.IsZero() {
			// Safari uses seconds since 2001-01-01
			const safariEpochDiff = 978307200
			safariStart := startDate.Unix() - safariEpochDiff
			query += ` AND hv.visit_time >= ?`
			args = append(args, safariStart)
		}

		if !endDate.IsZero() {
			// Only add 24 hours if the end time is at midnight (user specified just a date)
			endTimestamp := endDate.Unix()
			if endDate.Hour() == 0 && endDate.Minute() == 0 && endDate.Second() == 0 {
				endTimestamp += 86400
			}
			// Safari uses seconds since 2001-01-01
			const safariEpochDiff = 978307200
			safariEnd := endTimestamp - safariEpochDiff
			query += ` AND hv.visit_time < ?`
			args = append(args, safariEnd)
		}

		query += ` ORDER BY hv.visit_time DESC`
	} else {
		query = `
		SELECT
			hv.visit_time,
			hi.url,
			COALESCE(hv.title, hi.url) as title,
			hi.visit_count
		FROM history_visits hv
		JOIN history_items hi ON hv.history_item = hi.id
		WHERE hv.visit_time > 0
		ORDER BY hv.visit_time DESC
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
		var safariTime int64
		var url, title string
		var visitCount int

		if err := rows.Scan(&safariTime, &url, &title, &visitCount); err != nil {
			continue
		}

		timestamp := ConvertSafariTimestamp(safariTime)
		if timestamp.IsZero() {
			continue
		}

		entries = append(entries, models.HistoryEntry{
			Timestamp:  timestamp,
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			Domain:     ExtractDomain(url),
			Browser:    "safari",
		})
	}

	return entries, rows.Err()
}

// copyDatabase copies the Safari database to a temporary file
func (h *SafariHandler) copyDatabase() (string, error) {
	src, err := os.Open(h.dbPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.CreateTemp("", "web-recap-safari-*.db")
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
