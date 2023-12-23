package main

import (
	"fmt"
)

func handlePing(contents []string) (string, error) {
	if len(contents) == 1 {
		return "PONG", nil
	} else if len(contents) == 2 {
		return contents[1], nil
	} else {
		return "", fmt.Errorf("bad param")
	}
}
