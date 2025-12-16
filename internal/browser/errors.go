package browser

import "errors"

var (
	ErrUnsupportedPlatform    = errors.New("unsupported platform")
	ErrBrowserNotAvailable    = errors.New("browser not available on this platform")
	ErrUnknownBrowser         = errors.New("unknown browser type")
	ErrFirefoxProfileNotFound = errors.New("no Firefox profile found")
	ErrDatabaseNotFound       = errors.New("database file not found")
	ErrDatabaseLocked         = errors.New("database is locked")
)
