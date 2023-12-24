package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "6379"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Printf("Starting redis server on port %s ...\n", SERVER_PORT)
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening...", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Printf("Listening on %s ...\n", server.Addr())

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting...", err.Error())
			os.Exit(1)
		}
		go processClient(conn)
	}
}

func processClient(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		messageLen, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading...", err.Error())
		}
		var bufferStr string = string(buffer[:messageLen])
		fmt.Printf("Received: %s\n", strings.ReplaceAll(bufferStr, "\r\n", "\\r\\n"))

		messageContents, err := parseRESPMessage(buffer)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var output string = ""

		switch strings.ToUpper(messageContents[0]) {
		case "PING":
			res, err := handlePING(messageContents)
			if err != nil {
				output = ToSimpleError(err.Error())
			} else {
				output = ToBulkString(res)
			}

		case "ECHO":
			res, err := handleECHO(messageContents)
			if err != nil {
				output = ToSimpleError(err.Error())
			} else {
				output = ToBulkString(res)
			}

		case "GET":
			res, err := handleGET(messageContents)
			if err != nil {
				if err.Error() == "NULL" {
					output = ToNull()
				} else {
					output = ToSimpleError(err.Error())
				}
			} else {
				output = ToBulkString(res)
			}

		case "SET":
			res, err := handleSET(messageContents)
			if err != nil {
				output = ToSimpleError(err.Error())
			} else {
				output = ToSimpleString(res)
			}

		default:
			output = ToSimpleError(fmt.Sprintf("unknown command '%s'", messageContents[0]))
		}

		fmt.Printf("Sending: %s\n", strings.ReplaceAll(output, "\r\n", "\\r\\n"))
		conn.Write([]byte(output))
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
		fmt.Println("UNSUPPORTED")
	}

	if res != nil {
		return res, nil
	} else {
		return make([]string, 0), fmt.Errorf("Unsupported Message Type")
	}
}