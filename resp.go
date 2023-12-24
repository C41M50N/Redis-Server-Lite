package main

import (
	"fmt"
)

// https://redis.io/docs/reference/protocol-spec/#simple-strings
func ToSimpleString(value string) string {
	return fmt.Sprintf("+%s\r\n", value)
}

// https://redis.io/docs/reference/protocol-spec/#simple-errors
func ToSimpleError(value string) string {
	return fmt.Sprintf("-%s\r\n", value)
}

// https://redis.io/docs/reference/protocol-spec/#integers
func ToInteger(value int) string {
	return fmt.Sprintf(":%d\r\n", value)
}

// https://redis.io/docs/reference/protocol-spec/#bulk-strings
func ToBulkString(value string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
}

// https://redis.io/docs/reference/protocol-spec/#arrays
func ToArray(value []string) string {
	output := fmt.Sprintf("*%d\r\n", len(value))
	for _, element := range value {
		output += fmt.Sprintf("$%d\r\n%s\r\n", len(element), element)
	}
	return output
}

// https://redis.io/docs/reference/protocol-spec/#nulls
func ToNull() string {
	return "_\r\n"
}
