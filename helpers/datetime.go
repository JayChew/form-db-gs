package helpers

import (
	"reflect"
	"time"
)

// Define format strings for MySQL date/time types
var DbDateTimeFormats = map[string]string{
	"DATE":      "2006-01-02",          // MySQL DATE format
	"DATETIME":  "2006-01-02 15:04:05", // MySQL DATETIME format
	"TIMESTAMP": "2006-01-02 15:04:05", // MySQL TIMESTAMP format
	"TIME":      "15:04:05",            // MySQL TIME format
}

// FormatTimeField formats a time.Time field based on its mysql_format tag
func FormatDateTimeBasedOnTag(field reflect.Value, dateTimeFormatTag string) interface{} {
	// Get the format string from MysqlDateTimeFormats map
	format, ok := DbDateTimeFormats[dateTimeFormatTag]
	if ok {
		return field.Interface().(time.Time).Format(format)
	}

	// Handle unknown format tags
	return field.Interface().(time.Time).Format("2006-01-02") // default format
}
