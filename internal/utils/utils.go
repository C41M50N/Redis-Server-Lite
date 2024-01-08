package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Chuckinator2020/redis-server-lite-go/internal/r"
)

func ProcessClient(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		messageLen, err := conn.Read(buffer)
		if err != nil && messageLen == 0 {
			conn.Close()
			break
		}
		var bufferStr string = string(buffer[:messageLen])
		fmt.Printf("Received (%d): %s\n", messageLen, strings.ReplaceAll(bufferStr, "\r\n", "\\r\\n"))

		// handle redis-benchmark CONFIG request
		if bufferStr == "*3\r\n$6\r\nCONFIG\r\n$3\r\nGET\r\n$4\r\nsave\r\n*3\r\n$6\r\nCONFIG\r\n$3\r\nGET\r\n$10\r\nappendonly\r\n" {
			fmt.Println("CONFIG BS 4 redis-benchmark...")
			break
		}

		messageContents, err := parseRESPMessage(buffer)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var output []byte

		switch strings.ToUpper(messageContents[0]) {
		case "PING":
			res, err := HandlePING(messageContents)
			if err != nil {
				output = r.ToSimpleError(err.Error())
			} else {
				output = r.ToBulkString(res)
			}

		case "ECHO":
			res, err := HandleECHO(messageContents)
			if err != nil {
				output = r.ToSimpleError(err.Error())
			} else {
				output = r.ToBulkString(res)
			}

		case "GET":
			res, err := HandleGET(messageContents)
			if err != nil {
				if err.Error() == "NULL" {
					output = r.ToNull()
				} else {
					output = r.ToSimpleError(err.Error())
				}
			} else {
				output = r.ToBulkString(res)
			}

		case "SET":
			res, err := HandleSET(messageContents)
			if err != nil {
				output = r.ToSimpleError(err.Error())
			} else {
				output = r.ToSimpleString(res)
			}

		case "EXISTS":
			res, err := HandleEXISTS(messageContents)
			if err != nil {
				output = r.ToSimpleError(err.Error())
			} else {
				output = r.ToInteger(res)
			}

		case "DEL":
			res, err := HandleDEL(messageContents)
			if err != nil {
				output = r.ToSimpleError(err.Error())
			} else {
				output = r.ToInteger(res)
			}

		case "INCR":
			res, err := HandleINCR(messageContents)
			if err != nil {
				output = r.ToSimpleError(err.Error())
			} else {
				output = r.ToInteger(res)
			}

		case "DECR":
			res, err := HandleDECR(messageContents)
			if err != nil {
				output = r.ToSimpleError(err.Error())
			} else {
				output = r.ToInteger(res)
			}

		default:
			output = r.ToSimpleError(fmt.Sprintf("unknown command '%s'", messageContents[0]))
		}

		fmt.Printf("Sending: %s\n", strings.ReplaceAll(string(output), "\r\n", "\\r\\n"))
		conn.Write(output)
	}
}

func parseRESPMessage(buffer []byte) ([]string, error) {
	var res []string
	switch buffer[0] {
	case byte('*'): // array
		var arrayLen uint64

		var contents = bytes.SplitN(buffer, []byte("\r\n"), 2)

		var arrayLenBytes []byte = (contents[0])[1:]
		arrayLen, _ = strconv.ParseUint(string(arrayLenBytes), 10, 64)

		var rawArrayBytes = bytes.Split(contents[1], []byte("\r\n"))
		var arrayBytes = make([][]byte, 0)
		for i := 0; i < int(arrayLen); i++ {
			arrayBytes = append(arrayBytes, rawArrayBytes[(2*i)+1])
		}

		var arrayStrings []string
		for _, byteArray := range arrayBytes {
			arrayStrings = append(arrayStrings, string(byteArray))
		}
		var fancyArrayString, _ = json.Marshal(arrayStrings)

		fmt.Printf("Type: Array; Size: %d; Contents: %v\n", arrayLen, string(fancyArrayString))

		res = arrayStrings

	default:
		fmt.Println("recieved unsupported message type")
	}

	if res != nil {
		return res, nil
	} else {
		return make([]string, 0), fmt.Errorf("unsupported message type")
	}
}
