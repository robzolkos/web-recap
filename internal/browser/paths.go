package browser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetDatabasePath returns the database path for a given browser type on the current platform
func GetDatabasePath(browserType Type) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "linux":
		return getLinuxPath(home, browserType)
	case "darwin":
		return getDarwinPath(home, browserType)
	case "windows":
		return getWindowsPath(browserType)
	default:
		return "", ErrUnsupportedPlatform
	}
}

func getLinuxPath(home string, browserType Type) (string, error) {
	switch browserType {
	case Chrome:
		return filepath.Join(home, ".config/google-chrome/Default/History"), nil
	case Chromium:
		return filepath.Join(home, ".config/chromium/Default/History"), nil
	case Edge:
		return filepath.Join(home, ".config/microsoft-edge/Default/History"), nil
	case Brave:
		return filepath.Join(home, ".config/BraveSoftware/Brave-Browser/Default/History"), nil
	case Firefox:
		// Firefox uses profile directory, we'll handle this in detector
		return filepath.Join(home, ".mozilla/firefox"), nil
	case Safari:
		// Safari not available on Linux
		return "", ErrBrowserNotAvailable
	case Auto:
		return "", nil
	default:
		return "", ErrUnknownBrowser
	}
}

func getDarwinPath(home string, browserType Type) (string, error) {
	switch browserType {
	case Chrome:
		return filepath.Join(home, "Library/Application Support/Google/Chrome/Default/History"), nil
	case Chromium:
		return filepath.Join(home, "Library/Application Support/Chromium/Default/History"), nil
	case Edge:
		return filepath.Join(home, "Library/Application Support/Microsoft Edge/Default/History"), nil
	case Brave:
		return filepath.Join(home, "Library/Application Support/BraveSoftware/Brave-Browser/Default/History"), nil
	case Firefox:
		return filepath.Join(home, "Library/Application Support/Firefox"), nil
	case Safari:
		return filepath.Join(home, "Library/Safari/History.db"), nil
	case Auto:
		return "", nil
	default:
		return "", ErrUnknownBrowser
	}
}

func getWindowsPath(browserType Type) (string, error) {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		appData = filepath.Join(home, "AppData/Local")
	}

	switch browserType {
	case Chrome:
		return filepath.Join(appData, `Google\Chrome\User Data\Default\History`), nil
	case Chromium:
		return filepath.Join(appData, `Chromium\User Data\Default\History`), nil
	case Edge:
		return filepath.Join(appData, `Microsoft\Edge\User Data\Default\History`), nil
	case Brave:
		return filepath.Join(appData, `BraveSoftware\Brave-Browser\User Data\Default\History`), nil
	case Firefox:
		return filepath.Join(appData, "Mozilla/Firefox"), nil
	case Safari:
		// Safari not available on Windows
		return "", ErrBrowserNotAvailable
	case Auto:
		return "", nil
	default:
		return "", ErrUnknownBrowser
	}
}

// GetFirefoxProfilePath returns the active Firefox profile path
func GetFirefoxProfilePath(profileBaseDir string) (string, error) {
	if !fileExists(profileBaseDir) {
		return "", ErrFirefoxProfileNotFound
	}

	// Try to find the default profile or most recently modified profile
	entries, err := os.ReadDir(profileBaseDir)
	if err != nil {
		return "", err
	}

	var mostRecentPath string
	var mostRecentTime int64

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Look for .default-release or .default profiles first
		if strings.HasSuffix(name, ".default-release") || strings.HasSuffix(name, ".default") {
			placesPath := filepath.Join(profileBaseDir, name, "places.sqlite")
			if fileExists(placesPath) {
				return placesPath, nil
			}
		}

		// Otherwise, keep track of the most recently modified profile
		info, err := entry.Info()
		if err != nil {
			continue
		}

		modTime := info.ModTime().Unix()
		if modTime > mostRecentTime {
			mostRecentTime = modTime
			placesPath := filepath.Join(profileBaseDir, name, "places.sqlite")
			if fileExists(placesPath) {
				mostRecentPath = placesPath
			}
		}
	}

	if mostRecentPath != "" {
		return mostRecentPath, nil
	}

	return "", ErrFirefoxProfileNotFound
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
