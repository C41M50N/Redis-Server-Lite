package main

import (
	"fmt"
)

var db = make(map[string]string)

// https://redis.io/commands/ping/
func handlePING(contents []string) (string, error) {
	if len(contents) == 1 {
		return "PONG", nil
	} else if len(contents) == 2 {
		return contents[1], nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'ping' command")
	}
}

// https://redis.io/commands/echo/
func handleECHO(contents []string) (string, error) {
	if len(contents) == 2 {
		return contents[1], nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'echo' command")
	}
}

// https://redis.io/commands/set/
func handleSET(contents []string) (string, error) {
	if len(contents) == 3 {
		key := contents[1]
		value := contents[2]
		db[key] = value
		return "OK", nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'set' command")
	}
}

// https://redis.io/commands/get/
func handleGET(contents []string) (string, error) {
	if len(contents) == 2 {
		key := contents[1]
		value, ok := db[key]
		if ok != true {
			return "", fmt.Errorf("NULL")
		}
		return value, nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'get' command")
	}
}
