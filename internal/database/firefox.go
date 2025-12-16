package database

import (
	"database/sql"
	"io"
	"os"
	"time"

	"github.com/rzolkos/web-recap/internal/models"
	_ "modernc.org/sqlite"
)

// FirefoxHandler handles Firefox browser history
type FirefoxHandler struct {
	dbPath string
}

// NewFirefoxHandler creates a new Firefox history handler
func NewFirefoxHandler(dbPath string) *FirefoxHandler {
	return &FirefoxHandler{
		dbPath: dbPath,
	}
}

// GetHistory retrieves history entries from Firefox
func (h *FirefoxHandler) GetHistory(startDate, endDate time.Time) ([]models.HistoryEntry, error) {
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
	var query string
	var args []interface{}

	if !startDate.IsZero() || !endDate.IsZero() {
		query = `
		SELECT
			h.visit_date,
			p.url,
			p.title,
			p.visit_count
		FROM moz_historyvisits h
		JOIN moz_places p ON h.place_id = p.id
		WHERE h.visit_date > 0
		`

		if !startDate.IsZero() {
			// Firefox uses microseconds since epoch
			firefoxStart := startDate.Unix() * 1000000
			query += ` AND h.visit_date >= ?`
			args = append(args, firefoxStart)
		}

		if !endDate.IsZero() {
			// Only add 24 hours if the end time is at midnight (user specified just a date)
			endTimestamp := endDate.Unix()
			if endDate.Hour() == 0 && endDate.Minute() == 0 && endDate.Second() == 0 {
				endTimestamp += 86400
			}
			// Firefox uses microseconds since epoch
			firefoxEnd := endTimestamp * 1000000
			query += ` AND h.visit_date < ?`
			args = append(args, firefoxEnd)
		}

		query += ` ORDER BY h.visit_date DESC`
	} else {
		query = `
		SELECT
			h.visit_date,
			p.url,
			p.title,
			p.visit_count
		FROM moz_historyvisits h
		JOIN moz_places p ON h.place_id = p.id
		WHERE h.visit_date > 0
		ORDER BY h.visit_date DESC
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
		var firefoxTime int64
		var url, title string
		var visitCount int

		if err := rows.Scan(&firefoxTime, &url, &title, &visitCount); err != nil {
			continue
		}

		timestamp := ConvertFirefoxTimestamp(firefoxTime)
		if timestamp.IsZero() {
			continue
		}

		entries = append(entries, models.HistoryEntry{
			Timestamp:  timestamp,
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			Domain:     ExtractDomain(url),
			Browser:    "firefox",
		})
	}

	return entries, rows.Err()
}

// copyDatabase copies the Firefox database to a temporary file
func (h *FirefoxHandler) copyDatabase() (string, error) {
	src, err := os.Open(h.dbPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.CreateTemp("", "web-recap-firefox-*.db")
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
