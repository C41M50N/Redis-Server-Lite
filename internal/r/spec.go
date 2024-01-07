package r

import (
	"fmt"
)

type Bytes = []byte

// https://redis.io/docs/reference/protocol-spec/#simple-strings
func ToSimpleString(value string) Bytes {
	return Bytes(fmt.Sprintf("+%s\r\n", value))
}

// https://redis.io/docs/reference/protocol-spec/#simple-errors
func ToSimpleError(value string) Bytes {
	return Bytes(fmt.Sprintf("-%s\r\n", value))
}

// https://redis.io/docs/reference/protocol-spec/#integers
func ToInteger(value int) Bytes {
	return Bytes(fmt.Sprintf(":%d\r\n", value))
}

// https://redis.io/docs/reference/protocol-spec/#bulk-strings
func ToBulkString(value string) Bytes {
	return Bytes(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
}

// https://redis.io/docs/reference/protocol-spec/#arrays
func ToArray(value []string) Bytes {
	output := fmt.Sprintf("*%d\r\n", len(value))
	for _, element := range value {
		output += fmt.Sprintf("$%d\r\n%s\r\n", len(element), element)
	}
	return Bytes(output)
}

// https://redis.io/docs/reference/protocol-spec/#nulls
func ToNull() Bytes {
	return Bytes("_\r\n")
}

func ToNullArray() Bytes {
	return Bytes("*-1\r\n")
}
