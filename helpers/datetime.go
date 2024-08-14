package helpers

import (
	"regexp"
	"time"
)

// Regular expressions for date-time and time-only formats
var (
	// Strict ISO 8601 regex pattern
	iso8601Regex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z$`)
	// Time-only format (HH:MM:SS)
	timeOnlyRegex = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)
)

func FormatDateTime(input string) string {
	// Check if the string is in ISO 8601 format
	if iso8601Regex.MatchString(input) {
		// Attempt to parse as date-time
		t, err := time.Parse(time.RFC3339, input)
		if err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	// Check if the input is in time-only format (HH:MM:SS)
	if timeOnlyRegex.MatchString(input) {
		return input
	}

	// If parsing fails or not a recognized format, return as is
	return input
}
