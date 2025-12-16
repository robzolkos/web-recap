package browser

type Type string

const (
	Chrome   Type = "chrome"
	Chromium Type = "chromium"
	Edge     Type = "edge"
	Firefox  Type = "firefox"
	Safari   Type = "safari"
	Brave    Type = "brave"
	Auto     Type = "auto"
)

// Browser represents a detected browser with its database path
type Browser struct {
	Type Type
	Name string
	Path string
}
