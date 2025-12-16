---
name: web-recap
description: Extract and analyze browser history from Chrome, Chromium, Brave, Firefox, Safari, and Edge. Use when the user asks about their browsing history, web activity, visited websites, what they were reading, researching, or doing online. Outputs JSON for analysis of productivity, time tracking, research documentation, or daily summaries.
---

# web-recap SKILL

A Claude skill for extracting and analyzing browser history using the web-recap CLI tool.

## Overview

web-recap extracts browser history from Chrome, Chromium, Brave, Firefox, Safari, and Edge browsers, outputting JSON data that can be easily consumed by Claude and other language models for analysis, summarization, and insights.

## Installation

1. Install web-recap first (see [README.md](./README.md) for installation instructions)
2. Ensure web-recap binary is in your PATH

## Usage with Claude

### Basic Usage

```bash
# Get today's browser history
web-recap | claude --prompt "Summarize my web browsing today"

# Get history from a specific date
web-recap --date 2025-12-15 | claude --prompt "What were my main activities on this date?"

# Get history from a date range
web-recap --start-date 2025-12-01 --end-date 2025-12-15 | claude --prompt "Analyze my browsing patterns"

# Specify timezone for date interpretation
web-recap --date 2025-12-15 --tz America/New_York | claude --prompt "Summarize my activity on this date"

# Extract specific time period (12pm-1pm in your timezone)
web-recap --date 2025-12-15 --start-time 12:00 --end-time 13:00 | claude --prompt "What was I doing around noon?"

# Use UTC explicitly (useful for consistent results in automated workflows)
web-recap --date 2025-12-15 --utc | claude --prompt "Analyze my browsing"
```

### Common Analysis Prompts

```bash
# Productivity analysis
web-recap --date 2025-12-15 | claude --prompt "Analyze my browsing data and provide insights on my productivity today"

# Time tracking
web-recap --start-date 2025-12-09 --end-date 2025-12-15 | claude --prompt "Based on my browser history, how did I spend my time this week?"

# Research tracking
web-recap | claude --prompt "Identify research topics and domains I visited today"

# Distraction analysis
web-recap --date 2025-12-15 | claude --prompt "Which sites were likely distractions vs work-related?"

# Daily summary
web-recap --date 2025-12-15 | claude --prompt "Create a brief summary of my daily web activity"
```

## Output Format

web-recap outputs JSON in this format:

```json
{
  "browser": "chrome",
  "start_date": "2025-12-15T00:00:00Z",
  "end_date": "2025-12-15T23:59:59Z",
  "timezone": "America/New_York",
  "total_entries": 343,
  "entries": [
    {
      "timestamp": "2025-12-15T09:15:23Z",
      "url": "https://example.com/page",
      "title": "Example Page Title",
      "visit_count": 3,
      "domain": "example.com",
      "browser": "chrome"
    }
  ]
}
```

### JSON Field Reference

- **browser**: Source browser name
- **start_date**: Period start date (ISO 8601 UTC)
- **end_date**: Period end date (ISO 8601 UTC)
- **timezone**: Timezone used for date interpretation (shows which timezone was used when interpreting the requested dates)
- **total_entries**: Number of unique URLs visited
- **entries**: Array of history entries
  - **timestamp**: When the page was visited (always in UTC)
  - **url**: Full URL
  - **title**: Page title
  - **visit_count**: Total times this URL was visited
  - **domain**: Extracted domain
  - **browser**: Browser source

## Capabilities

### Data Extraction
- Extract history from any supported browser (Chrome, Chromium, Brave, Firefox, Safari, Edge)
- Filter by specific dates or date ranges
- Extract from multiple browsers simultaneously
- Custom database paths for non-standard installations

### LLM-Friendly Features
- Structured JSON output
- Clean URL and domain extraction
- Visit timestamps in ISO 8601 format
- Page titles for context
- Visit counts for frequency analysis

### Privacy-Safe
- No automatic transmission of data
- Full local processing
- Control over which browser history to extract
- Output only goes to stdout (you control where it goes)

## Advanced Usage

### Combine with other tools

```bash
# Export to file, then analyze
web-recap --all-browsers -o ~/web-history.json

# Pipe to jq for filtering
web-recap | jq '.entries[] | select(.domain | contains("github"))'

# Multiple days analysis
for i in {1..7}; do
  date -d "-$i days" +%Y-%m-%d | xargs -I {} web-recap --date {}
done | claude --prompt "Analyze my web activity trends"
```

### Combine multiple browsers

```bash
# Get history from all browsers
web-recap --all-browsers | claude --prompt "Which browser did I use most? What was I doing?"
```

### Time-based analysis

```bash
# Morning browsing (8am-12pm in your timezone)
web-recap --date "$(date +%Y-%m-%d)" --start-time 08:00 --end-time 12:00 | \
  claude --prompt "What did I do online this morning?"

# Afternoon work session (1pm-5pm)
web-recap --date "$(date +%Y-%m-%d)" --start-time 13:00 --end-time 17:00 | \
  claude --prompt "Summarize my afternoon work activities"

# Weekly review (use system timezone)
web-recap --start-date 2025-12-09 --end-date 2025-12-15 | \
  claude --prompt "Weekly summary: what were my main online activities?"

# Monthly report (forcing UTC for consistency)
web-recap --start-date 2025-12-01 --end-date 2025-12-31 --utc | \
  claude --prompt "Monthly internet usage report"

# Yesterday's activity in a specific timezone
web-recap --date "$(date -d yesterday +%Y-%m-%d)" --tz Europe/London | \
  claude --prompt "Summarize my activity from yesterday in London time"
```

## Use Cases

### 1. Daily Standup Summary
```bash
# Uses your local timezone by default
web-recap --date "$(date -d yesterday +%Y-%m-%d)" | \
  claude --prompt "Create a brief summary of my web-based work activities from yesterday for a standup"

# Or specify a timezone if working across time zones
web-recap --date "$(date -d yesterday +%Y-%m-%d)" --tz America/New_York | \
  claude --prompt "Create a brief summary of my web-based work activities from yesterday (EST) for a standup"
```

### 2. Research Documentation
```bash
web-recap --date 2025-12-15 | \
  claude --prompt "List the research sources I visited today and suggest how they connect"
```

### 3. Time Tracking
```bash
web-recap --start-date 2025-12-09 --end-date 2025-12-15 | \
  claude --prompt "Based on domains visited, estimate how much time I spent on work vs personal projects"
```

### 4. Distraction Audit
```bash
web-recap | \
  claude --prompt "Identify potential time-wasting sites and suggest productivity improvements"
```

### 5. Client/Project Analysis
```bash
web-recap --date 2025-12-15 | \
  claude --prompt "List all work-related domains I visited today for project ProjectX"
```

### 6. Learning Path Documentation
```bash
web-recap --start-date 2025-12-01 --end-date 2025-12-15 | \
  claude --prompt "Summarize the technical topics I researched this month"
```

## Timezone Considerations

### Understanding Timezone Behavior

By default, web-recap interprets dates in your **system's local timezone**. All output timestamps are in **UTC** for consistency:

- When you request `--date 2025-12-15`, it means December 15th in your local timezone
- The `timezone` field in the JSON shows which timezone was used
- Use `--tz` to explicitly specify a different timezone
- Use `--utc` to force UTC interpretation (useful for scripts and CI/CD)

### Common Patterns

```bash
# Query based on your local time (recommended for daily use)
web-recap --date 2025-12-15

# Query based on a specific timezone (multi-timezone teams)
web-recap --date 2025-12-15 --tz America/Los_Angeles

# Force UTC for reproducible results in automation
web-recap --date 2025-12-15 --utc

# Extract specific hours of your day
web-recap --date 2025-12-15 --start-time 09:00 --end-time 17:00
```

### Time Zone Names

Use IANA timezone identifiers (see [list of timezones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)):
- `America/New_York`, `America/Chicago`, `America/Los_Angeles`
- `Europe/London`, `Europe/Paris`, `Europe/Berlin`
- `Asia/Tokyo`, `Asia/Shanghai`, `Asia/Hong_Kong`
- `Australia/Sydney`
- `UTC` for Coordinated Universal Time

## Limitations

- **Browser must be closed** or database will be locked (tool handles this by copying, but may take time)
- **Safari only on macOS** - not available on Linux or Windows
- **Firefox multiple profiles** - tool uses the most recently modified profile
- **History retention** - depends on browser settings (typically 3-6 months)

## Tips for Better Results

1. **Be specific with dates** - provide exact dates for targeted analysis
2. **Include context in prompts** - tell Claude what you're analyzing for
3. **Use all-browsers flag** - get complete picture of your browsing
4. **Regular exports** - save history periodically before browser retention limits
5. **Combine with other data** - use alongside calendar, task managers for complete picture

## Integration with Workflows

### Daily Review
```bash
# Add to your daily standup script
web-recap --date "$(date +%Y-%m-%d)" > /tmp/web-recap.json
claude --context /tmp/web-recap.json --prompt "Daily summary"
```

### Weekly Analysis
```bash
#!/bin/bash
START=$(date -d "last monday" +%Y-%m-%d)
END=$(date +%Y-%m-%d)
web-recap --start-date $START --end-date $END | \
  claude --prompt "Weekly web activity analysis"
```

### Scheduled Exports
```bash
# Add to crontab for automatic daily exports
0 23 * * * /usr/local/bin/web-recap --all-browsers -o ~/history/$(date +\%Y-\%m-\%d).json
```

## Troubleshooting

### "command not found: web-recap"
- Ensure web-recap is installed and in your PATH
- Try the full path: `/usr/local/bin/web-recap`

### "No browsers detected"
- Check browser is installed at default location
- Use `web-recap list` to see detected browsers
- Use `--db-path` to specify custom location

### Empty results
- Browser databases may need to be locked-free (close browser)
- Check date range is correct
- Verify browser history retention settings

### Large outputs
- Use jq or other tools to filter JSON before passing to Claude
- Consider date ranges instead of all-time history

## Privacy Considerations

- web-recap processes history entirely locally
- No data is sent to external servers unless YOU pipe it to one
- You have full control over what data is extracted
- The tool does not modify or delete any browser history

## Contributing

Found issues or have suggestions? Please open an issue or PR on GitHub.

## See Also

- [README.md](./README.md) - Installation and basic usage
- [Supported Browsers](./README.md#supported-browsers) - Platform-specific information
