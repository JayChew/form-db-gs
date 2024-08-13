package helpers

import (
	"strings"
	"time"
)

func FormatDateTime(input string) string {
	// Check if the string is in ISO 8601 format
	if strings.Contains(input, "T") {
		// Attempt to parse as date-time
		t, err := time.Parse(time.RFC3339, input)
		if err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	// Check if the input is in time-only format (HH:MM:SS)
	// _, err := time.Parse("15:04:05", input)
	// if err == nil {
	// 	fmt.Println(input)
	// 	return input
	// }

	// If parsing fails or not a recognized format, return as is
	return input
}
