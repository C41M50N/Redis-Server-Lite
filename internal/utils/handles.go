package utils

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"
)

// goroutine (concurrency) safe key-value store
var db = sync.Map{}

// https://redis.io/commands/ping/
func HandlePING(contents []string) (string, error) {
	if len(contents) == 1 {
		return "PONG", nil
	} else if len(contents) == 2 {
		return contents[1], nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'ping' command")
	}
}

// https://redis.io/commands/echo/
func HandleECHO(contents []string) (string, error) {
	if len(contents) == 2 {
		return contents[1], nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'echo' command")
	}
}

// https://redis.io/commands/set/
func HandleSET(contents []string) (string, error) {
	if len(contents) == 3 {
		key := contents[1]
		value := contents[2]
		db.Store(key, value)
		return "OK", nil
	} else if len(contents) == 5 {
		key := contents[1]
		value := contents[2]
		switch contents[3] {
		case "EX":
			delta, err := strconv.Atoi(contents[4])
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)
			time.AfterFunc(time.Duration(delta)*time.Second, func() { db.Delete(key) })
			return "OK", nil

		case "PX":
			delta, err := strconv.Atoi(contents[4])
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)
			time.AfterFunc(time.Duration(delta)*time.Millisecond, func() { db.Delete(key) })
			return "OK", nil

		case "EXAT":
			timestamp, err := strconv.ParseInt(contents[4], 10, 64)
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if timestamp <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)

			delta := timestamp - time.Now().Unix()
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			time.AfterFunc(time.Duration(delta)*time.Second, func() { db.Delete(key) })
			return "OK", nil

		case "PXAT":
			timestamp, err := strconv.ParseInt(contents[4], 10, 64)
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if timestamp <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)

			delta := timestamp - time.Now().UnixMilli()
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			time.AfterFunc(time.Duration(delta)*time.Millisecond, func() { db.Delete(key) })
			return "OK", nil

		default:
			return "", fmt.Errorf("syntax error")
		}
	}
	return "", fmt.Errorf("wrong number of arguments for 'set' command")
}

// https://redis.io/commands/get/
func HandleGET(contents []string) (string, error) {
	if len(contents) == 2 {
		key := contents[1]
		value, ok := db.Load(key)
		if !ok {
			return "", fmt.Errorf("NULL")
		}
		return value.(string), nil
	}
	return "", fmt.Errorf("wrong number of arguments for 'get' command")
}

// https://redis.io/commands/exists/
func HandleEXISTS(contents []string) (int, error) {
	if len(contents) >= 2 {
		count := 0
		keys := contents[1:]
		for _, key := range keys {
			_, ok := db.Load(key)
			if ok {
				count++
			}
		}
		return count, nil
	}
	return -1, fmt.Errorf("wrong number of arguments for 'EXISTS' command")
}

// https://redis.io/commands/del/
func HandleDEL(contents []string) (int, error) {
	if len(contents) >= 2 {
		count := 0
		keys := contents[1:]
		for _, key := range keys {
			_, loaded := db.LoadAndDelete(key)
			if loaded {
				count++
			}
		}
		return count, nil
	}
	return -1, fmt.Errorf("wrong number of arguments for 'DEL' command")
}

// https://redis.io/commands/incr/
func HandleINCR(contents []string) (int, error) {
	if len(contents) == 2 {
		key := contents[1]
		value, ok := db.Load(key)
		if !ok {
			value = "0"
		}
		intValue, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return -1, fmt.Errorf("value is not an integer or out of range")
		}
		intValue++
		db.Store(key, fmt.Sprint(intValue))
		return int(intValue), nil
	}
	return -1, fmt.Errorf("wrong number of arguments for 'INCR' command")
}

// https://redis.io/commands/decr/
func HandleDECR(contents []string) (int, error) {
	if len(contents) == 2 {
		key := contents[1]
		value, ok := db.Load(key)
		if !ok {
			value = "0"
		}
		intValue, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return -1, fmt.Errorf("value is not an integer or out of range")
		}
		intValue--
		db.Store(key, fmt.Sprint(intValue))
		return int(intValue), nil
	}
	return -1, fmt.Errorf("wrong number of arguments for 'DECR' command")
}

// https://redis.io/commands/lpush/
func HandleLPUSH(contents []string) (int, error) {
	if len(contents) >= 3 {
		key := contents[1]
		elements := contents[2:]
		slices.Reverse(elements)

		var listValue []string

		value, ok := db.Load(key)
		if !ok {
			listValue = make([]string, 0)
		} else {
			var ok bool
			listValue, ok = value.([]string)
			if !ok {
				return -1, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
			}
		}

		listValue = append(elements, listValue...)
		db.Store(key, listValue)
		return len(listValue), nil
	}
	return -1, fmt.Errorf("wrong number of arguments for 'LPUSH' command")
}

func HandleRPUSH(contents []string) (int, error) {
	if len(contents) >= 3 {
		key := contents[1]
		elements := contents[2:]

		var listValue []string

		value, ok := db.Load(key)
		if !ok {
			listValue = make([]string, 0)
		} else {
			var ok bool
			listValue, ok = value.([]string)
			if !ok {
				return -1, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
			}
		}

		listValue = append(listValue, elements...)
		db.Store(key, listValue)
		return len(listValue), nil
	}
	return -1, fmt.Errorf("wrong number of arguments for 'RPUSH' command")
}

func HandleLRANGE(contents []string) ([]string, error) {
	if len(contents) == 4 {
		key := contents[1]

		var list []string

		value, ok := db.Load(key)
		if !ok {
			return make([]string, 0), nil
		} else {
			var ok bool
			list, ok = value.([]string)
			if !ok {
				return []string{}, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
			}
		}

		start, err := strconv.ParseInt(contents[2], 10, 64)
		if err != nil {
			return []string{}, fmt.Errorf("value is not an integer or out of range")
		}

		stop, err := strconv.ParseInt(contents[3], 10, 64)
		if err != nil {
			return []string{}, fmt.Errorf("value is not an integer or out of range")
		}

		fmt.Println(start, stop, list)

		if int(start) >= len(list) {
			return make([]string, 0), nil
		} else if start == 0 && stop == -1 {
			return list, nil
		}

		if start < 0 {
			start = int64(len(list) + int(start))
			if start < 0 {
				start = 0
			}
		}

		if stop < 0 {
			stop = int64(len(list) + int(stop) + 1)
			if stop < 0 {
				stop = int64(len(list))
			}
		} else if int(stop) > len(list) {
			stop = int64(len(list))
		} else {
			stop++
		}

		if start > stop {
			return make([]string, 0), nil
		}

		return list[start:stop], nil
	}
	return []string{}, fmt.Errorf("wrong number of arguments for 'LRANGE' command")
}
