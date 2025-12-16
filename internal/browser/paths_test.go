package browser

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestGetLinuxPath(t *testing.T) {
	tests := []struct {
		name      string
		browser   Type
		expectErr bool
		contains  string
	}{
		{
			name:     "Chrome",
			browser:  Chrome,
			contains: ".config/google-chrome",
		},
		{
			name:     "Chromium",
			browser:  Chromium,
			contains: ".config/chromium",
		},
		{
			name:     "Firefox",
			browser:  Firefox,
			contains: ".mozilla/firefox",
		},
		{
			name:      "Safari not available",
			browser:   Safari,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := "/home/testuser"
			path, err := getLinuxPath(home, tt.browser)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectErr && !strings.Contains(filepath.ToSlash(path), tt.contains) {
				t.Errorf("expected path to contain %q, got %q", tt.contains, path)
			}
		})
	}
}

func TestGetDarwinPath(t *testing.T) {
	tests := []struct {
		name      string
		browser   Type
		expectErr bool
		contains  string
	}{
		{
			name:     "Chrome",
			browser:  Chrome,
			contains: "Library/Application Support/Google/Chrome",
		},
		{
			name:     "Firefox",
			browser:  Firefox,
			contains: "Library/Application Support/Firefox",
		},
		{
			name:     "Safari",
			browser:  Safari,
			contains: "Library/Safari/History.db",
		},
		{
			name:      "Edge",
			browser:   Edge,
			contains:  "Library/Application Support/Microsoft Edge",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := "/Users/testuser"
			path, err := getDarwinPath(home, tt.browser)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectErr && !strings.Contains(filepath.ToSlash(path), tt.contains) {
				t.Errorf("expected path to contain %q, got %q", tt.contains, path)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	// Note: ExtractDomain is in the database package, so we'd need to import it there
	// For now, this is a placeholder for domain extraction tests
}
