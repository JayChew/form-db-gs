package helpers

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateInsertIntoQuery(model interface{}, tableName string) (string, []interface{}, error) {
	v := reflect.ValueOf(model)
  t := reflect.TypeOf(model)

	if v.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag == "" || dbTag == "-" {
			continue
		}

		columns = append(columns, dbTag)
		placeholders = append(placeholders, "?")
		values = append(values, v.Field(i).Interface())
	}

	// Construct the SQL query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, values, nil
}