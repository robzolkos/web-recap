package database

import "errors"

var (
	ErrSafariNotAvailable = errors.New("Safari is only available on macOS")
	ErrUnsupportedBrowser = errors.New("unsupported browser type")
	ErrDatabaseError      = errors.New("database error")
)
