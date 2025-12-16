package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/rzolkos/web-recap/internal/browser"
	"github.com/rzolkos/web-recap/internal/database"
	"github.com/rzolkos/web-recap/internal/output"
)

var (
	browserType    string
	date           string
	startDate      string
	endDate        string
	startTime      string
	endTime        string
	timeHour       string
	timezone       string
	utcMode        bool
	outputFile     string
	dbPath         string
	allBrowsers    bool
	version        = "0.1.0-alpha"
)

var rootCmd = &cobra.Command{
	Use:   "web-recap",
	Short: "Extract browser history in LLM-friendly JSON format",
	Long: `web-recap extracts browser history from Chrome, Chromium, Firefox, Safari, and Edge
and outputs it in JSON format suitable for analysis by LLMs and other tools.

Date and time inputs are interpreted in your local timezone by default.

Examples:
  web-recap                          # Extract today's history from default browser
  web-recap --browser chrome         # Extract from Chrome specifically
  web-recap --date 2025-12-15        # Extract history from specific date (local time)
  web-recap --date 2025-12-15 --time 12           # Extract 12pm hour (12:00-12:59)
  web-recap --date 2025-12-15 --start-time 12:00 --end-time 13:00  # Time range
  web-recap --tz America/New_York --date 2025-12-15  # Explicit timezone
  web-recap --start-date 2025-12-01 --end-date 2025-12-15  # Date range
  web-recap --all-browsers -o history.json  # All browsers to file
`,
	RunE: runWeb,
}

func init() {
	rootCmd.Flags().StringVarP(&browserType, "browser", "b", "auto", "Browser type: auto, chrome, chromium, edge, brave, firefox, or safari")
	rootCmd.Flags().StringVar(&date, "date", "", "Specific date (YYYY-MM-DD, interpreted in local timezone)")
	rootCmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD, interpreted in local timezone)")
	rootCmd.Flags().StringVar(&endDate, "end-date", "", "End date (YYYY-MM-DD, interpreted in local timezone)")
	rootCmd.Flags().StringVar(&startTime, "start-time", "", "Start time (HH:MM format)")
	rootCmd.Flags().StringVar(&endTime, "end-time", "", "End time (HH:MM format)")
	rootCmd.Flags().StringVar(&timeHour, "time", "", "Time hour shorthand (e.g., '12' for 12:00-12:59)")
	rootCmd.Flags().StringVar(&timezone, "tz", "", "Timezone (e.g., America/New_York, UTC, local for system timezone)")
	rootCmd.Flags().BoolVar(&utcMode, "utc", false, "Treat all dates/times as UTC instead of local timezone")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().StringVar(&dbPath, "db-path", "", "Custom database path")
	rootCmd.Flags().BoolVar(&allBrowsers, "all-browsers", false, "Extract from all detected browsers")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// getTimezone returns the appropriate timezone based on flags
func getTimezone(tzFlag string, utcFlag bool) (*time.Location, error) {
	if utcFlag {
		return time.UTC, nil
	}

	if tzFlag != "" {
		if tzFlag == "local" {
			return time.Local, nil
		}
		loc, err := time.LoadLocation(tzFlag)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %q: %v", tzFlag, err)
		}
		return loc, nil
	}

	// Default to system local timezone
	return time.Local, nil
}

// parseDateTimeInLocation parses a date and optional time in a specific timezone
func parseDateTimeInLocation(dateStr, timeStr string, loc *time.Location) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	// Parse date
	dateTime, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}

	if timeStr == "" {
		// No time specified, use start of day
		return dateTime, nil
	}

	// Parse time
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format (use HH:MM): %v", err)
	}

	// Combine date + time in the specified timezone
	return time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(),
		t.Hour(), t.Minute(), 0, 0, loc), nil
}

// parseHour parses a single hour value (0-23)
func parseHour(hourStr string) (int, error) {
	var hour int
	_, err := fmt.Sscanf(hourStr, "%d", &hour)
	if err != nil || hour < 0 || hour > 23 {
		return 0, fmt.Errorf("invalid hour (must be 0-23): %s", hourStr)
	}
	return hour, nil
}

func runWeb(cmd *cobra.Command, args []string) error {
	// Get timezone
	loc, err := getTimezone(timezone, utcMode)
	if err != nil {
		return err
	}

	// Parse dates with timezone
	var startTimeValue, endTimeValue time.Time
	var err2 error

	if date != "" {
		// Single date mode
		start, err := parseDateTimeInLocation(date, "", loc)
		if err != nil {
			return err
		}

		if timeHour != "" {
			// --time 12 means 12:00-12:59
			hour, err := parseHour(timeHour)
			if err != nil {
				return err
			}
			startTimeValue = time.Date(start.Year(), start.Month(), start.Day(),
				hour, 0, 0, 0, loc)
			endTimeValue = startTimeValue.Add(1 * time.Hour)
		} else if startTime != "" || endTime != "" {
			// Explicit time range
			var st, et string
			if startTime != "" {
				st = startTime
			} else {
				st = "00:00"
			}
			if endTime != "" {
				et = endTime
			} else {
				et = "23:59"
			}

			startTimeValue, err = parseDateTimeInLocation(date, st, loc)
			if err != nil {
				return err
			}
			endTimeValue, err = parseDateTimeInLocation(date, et, loc)
			if err != nil {
				return err
			}
		} else {
			// Full day
			startTimeValue = start
			endTimeValue = start.Add(24 * time.Hour)
		}
	} else if startDate != "" || endDate != "" {
		// Date range mode (existing logic, updated to use timezone)
		if startDate != "" {
			startTimeValue, err2 = parseDateTimeInLocation(startDate, "", loc)
			if err2 != nil {
				return err2
			}
		}

		if endDate != "" {
			endTimeValue, err2 = parseDateTimeInLocation(endDate, "", loc)
			if err2 != nil {
				return err2
			}
			endTimeValue = endTimeValue.Add(24 * time.Hour)
		}
	} else {
		// No date specified - default to today
		now := time.Now().In(loc)
		startTimeValue = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endTimeValue = startTimeValue.Add(24 * time.Hour)
	}

	// Convert to UTC for database query (important!)
	startTimeValue = startTimeValue.UTC()
	endTimeValue = endTimeValue.UTC()

	// Get browser
	detector := browser.NewDetector()
	var b *browser.Browser

	// Default to all browsers if no specific browser and no --all-browsers flag
	useAllBrowsers := allBrowsers || browserType == "auto"

	if useAllBrowsers {
		// Handle multiple browsers
		entries, err := database.QueryMultipleBrowsers(detector, startTimeValue, endTimeValue)
		if err != nil {
			return fmt.Errorf("failed to query browsers: %v", err)
		}

		// Write output
		out := os.Stdout
		if outputFile != "" {
			f, err := os.Create(outputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file: %v", err)
			}
			defer f.Close()
			out = f
		}

		return output.FormatJSON(out, entries, "all", startTimeValue, endTimeValue, timezone)
	}

	// Get specific browser
	bType := browser.Type(browserType)
	if dbPath != "" {
		// Validate custom path
		info, err := os.Stat(dbPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("database file not found: %s", dbPath)
			}
			return fmt.Errorf("cannot access database file: %v", err)
		}
		if info.IsDir() {
			return fmt.Errorf("path is a directory, not a file: %s", dbPath)
		}

		// Use custom path
		b = &browser.Browser{
			Type: bType,
			Name: string(bType),
			Path: dbPath,
		}
	} else {
		var err error
		b, err = detector.GetBrowser(bType)
		if err != nil {
			return fmt.Errorf("failed to get browser: %v", err)
		}
	}

	// Query history
	entries, err := database.Query(b, startTimeValue, endTimeValue)
	if err != nil {
		return fmt.Errorf("failed to query history: %v", err)
	}

	// Write output
	out := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer f.Close()
		out = f
	}

	return output.FormatJSON(out, entries, b.Name, startTimeValue, endTimeValue, timezone)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("web-recap version %s\n", version)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List detected browsers",
	RunE: func(cmd *cobra.Command, args []string) error {
		detector := browser.NewDetector()
		browsers := detector.Detect()

		if len(browsers) == 0 {
			fmt.Println("No browsers detected")
			return nil
		}

		fmt.Println("Detected browsers:")
		for _, b := range browsers {
			fmt.Printf("  - %s (%s): %s\n", b.Name, b.Type, b.Path)
		}

		return nil
	},
}
