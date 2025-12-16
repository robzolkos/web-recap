package database

import (
	"net/url"
	"strings"
	"time"
)

// ConvertChromeTimestamp converts Chrome's timestamp format (microseconds since 1601-01-01) to Unix time
func ConvertChromeTimestamp(chromeTime int64) time.Time {
	// Chrome timestamp is in microseconds since 1601-01-01
	// Unix epoch is 1970-01-01
	// Difference: 11644473600 seconds = 11644473600000000 microseconds
	const chromeEpochDiff = 11644473600

	if chromeTime == 0 {
		return time.Time{}
	}

	unixSeconds := (chromeTime / 1000000) - chromeEpochDiff
	return time.Unix(unixSeconds, 0).UTC()
}

// ConvertFirefoxTimestamp converts Firefox's timestamp format (microseconds since epoch) to Unix time
func ConvertFirefoxTimestamp(firefoxTime int64) time.Time {
	if firefoxTime == 0 {
		return time.Time{}
	}

	// Firefox uses microseconds since Unix epoch
	unixSeconds := firefoxTime / 1000000
	unixNanos := (firefoxTime % 1000000) * 1000
	return time.Unix(unixSeconds, unixNanos).UTC()
}

// ConvertSafariTimestamp converts Safari's timestamp format (seconds since 2001-01-01) to Unix time
func ConvertSafariTimestamp(safariTime int64) time.Time {
	// Safari uses seconds since 2001-01-01
	// Unix epoch is 1970-01-01
	// Difference: 978307200 seconds
	const safariEpochDiff = 978307200

	unixSeconds := safariTime + safariEpochDiff
	return time.Unix(unixSeconds, 0).UTC()
}

// ExtractDomain extracts the domain from a URL string
func ExtractDomain(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	// Try to parse as URL
	u, err := url.Parse(urlStr)
	if err != nil {
		// If parsing fails, try to extract domain manually
		if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
			parts := strings.Split(urlStr, "/")
			if len(parts) > 2 {
				return parts[2]
			}
		}
		return urlStr
	}

	if u.Host != "" {
		return u.Host
	}

	return urlStr
}

// FilterByDateRange filters history entries by date range
func FilterByDateRange(entries []interface{}, startDate, endDate time.Time) []interface{} {
	if startDate.IsZero() && endDate.IsZero() {
		return entries
	}

	// Normalize dates to start of day
	if !startDate.IsZero() {
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	}

	if !endDate.IsZero() {
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
	}

	var filtered []interface{}
	for _, entry := range entries {
		filtered = append(filtered, entry)
	}

	return filtered
}
