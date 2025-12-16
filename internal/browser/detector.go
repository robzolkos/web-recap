package browser

// Detector detects available browsers on the system
type Detector struct{}

// NewDetector creates a new browser detector
func NewDetector() *Detector {
	return &Detector{}
}

// Detect returns a list of available browsers
func (d *Detector) Detect() []Browser {
	var browsers []Browser

	// Check each browser type
	for _, bType := range []Type{Chrome, Chromium, Edge, Brave, Firefox, Safari} {
		path, err := GetDatabasePath(bType)
		if err != nil {
			continue
		}

		// For Firefox, handle profile detection
		if bType == Firefox {
			profilePath, err := GetFirefoxProfilePath(path)
			if err == nil {
				browsers = append(browsers, Browser{
					Type: Firefox,
					Name: "Firefox",
					Path: profilePath,
				})
			}
			continue
		}

		// For other browsers, check if the database file exists
		if fileExists(path) {
			name := string(bType)
			if bType == Chrome {
				name = "Google Chrome"
			} else if bType == Chromium {
				name = "Chromium"
			} else if bType == Edge {
				name = "Microsoft Edge"
			} else if bType == Brave {
				name = "Brave"
			} else if bType == Safari {
				name = "Safari"
			}

			browsers = append(browsers, Browser{
				Type: bType,
				Name: name,
				Path: path,
			})
		}
	}

	return browsers
}

// GetBrowser returns a specific browser, detecting if necessary
func (d *Detector) GetBrowser(browserType Type) (*Browser, error) {
	if browserType == Auto {
		browsers := d.Detect()
		if len(browsers) == 0 {
			return nil, ErrDatabaseNotFound
		}
		// Return the first detected browser
		return &browsers[0], nil
	}

	path, err := GetDatabasePath(browserType)
	if err != nil {
		return nil, err
	}

	// For Firefox, handle profile detection
	if browserType == Firefox {
		profilePath, err := GetFirefoxProfilePath(path)
		if err != nil {
			return nil, err
		}
		return &Browser{
			Type: Firefox,
			Name: "Firefox",
			Path: profilePath,
		}, nil
	}

	// For other browsers, check if the database file exists
	if !fileExists(path) {
		return nil, ErrDatabaseNotFound
	}

	name := string(browserType)
	if browserType == Chrome {
		name = "Google Chrome"
	} else if browserType == Chromium {
		name = "Chromium"
	} else if browserType == Edge {
		name = "Microsoft Edge"
	} else if browserType == Brave {
		name = "Brave"
	} else if browserType == Safari {
		name = "Safari"
	}

	return &Browser{
		Type: browserType,
		Name: name,
		Path: path,
	}, nil
}
