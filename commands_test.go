package main

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createMockConnection() net.Conn {
	server, client := net.Pipe()
	go func() {
		processClient(server)
		server.Close()
	}()

	return client
}

func readBuffer(client net.Conn) []byte {
	buffer := make([]byte, 1024)
	messageLen, err := client.Read(buffer)
	if err != nil {
		fmt.Printf("Issue Reading: %s\n", err.Error())
	}
	return buffer[:messageLen]
}

func TestUnknownCommand(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	client.Write([]byte(ToArray([]string{"PEEK"})))
	response := readBuffer(client)

	assert.Equal(t, []byte("-unknown command 'PEEK'\r\n"), response)
}

func TestPing1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	client.Write([]byte(ToArray([]string{"PING"})))
	response := readBuffer(client)

	assert.Equal(t, []byte(ToBulkString("PONG")), response)
}

func TestPing2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	arg := "rainbow"

	client.Write([]byte(ToArray([]string{"PING", arg})))
	response := readBuffer(client)

	assert.Equal(t, []byte(ToBulkString(arg)), response)
}

func TestPing3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"echo", "chamber"}

	client.Write([]byte(ToArray(append([]string{"PING"}, args...))))
	response := readBuffer(client)

	assert.Equal(t, []byte(ToSimpleError("wrong number of arguments for 'ping' command")), response)
}

func TestEcho1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	arg := "Wubba lubba dub dub!"

	client.Write([]byte(ToArray([]string{"ECHO", arg})))
	response := readBuffer(client)

	assert.Equal(t, []byte(ToBulkString(arg)), response)
}

func TestEcho2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"Wubba lubba dub dub!", "Morty!!!"}

	client.Write([]byte(ToArray(append([]string{"ECHO"}, args...))))
	response := readBuffer(client)

	assert.Equal(t, []byte(ToSimpleError("wrong number of arguments for 'echo' command")), response)
}

func TestStorage1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value := "salary", "123456"

	client.Write([]byte(ToArray(append([]string{"SET"}, key, value))))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleString("OK")), response)

	client.Write([]byte(ToArray(append([]string{"GET"}, key))))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToBulkString(value)), response)
}

func TestStorage2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value := "description", "thing thing something thing"

	client.Write([]byte(ToArray(append([]string{"SET"}, key, value))))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleString("OK")), response)

	client.Write([]byte(ToArray(append([]string{"GET"}, key))))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToBulkString(value)), response)
}

func TestStorage3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key := "key"

	client.Write([]byte(ToArray(append([]string{"SET"}, key))))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("wrong number of arguments for 'set' command")), response)
}

func TestStorageEX1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", 3
	args := []string{"SET", key, value, "EX", strconv.Itoa(exp)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleString("OK")), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToBulkString(value)), response)

	// GET after expiration
	time.Sleep(time.Duration(exp) * time.Second)
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToNull()), response)
}

func TestStorageEX2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "EX", exp}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("value is not an integer or out of range")), response)
}

func TestStorageEX3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "EX", strconv.Itoa(exp)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("invalid expire time in 'set' command")), response)
}

func TestStoragePX1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", 3000
	args := []string{"SET", key, value, "PX", strconv.Itoa(exp)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleString("OK")), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToBulkString(value)), response)

	// GET after expiration
	time.Sleep(time.Duration(exp) * time.Millisecond)
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToNull()), response)
}

func TestStoragePX2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "PX", exp}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("value is not an integer or out of range")), response)
}

func TestStoragePX3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "PX", strconv.Itoa(exp)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("invalid expire time in 'set' command")), response)
}

func TestStorageEXAT1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", time.Now().Add(3*time.Second).Unix()
	args := []string{"SET", key, value, "EXAT", strconv.FormatInt(exp, 10)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleString("OK")), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToBulkString(value)), response)

	// GET after expiration
	time.Sleep(time.Duration(exp-time.Now().Unix()) * time.Second)
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToNull()), response)
}

func TestStorageEXAT2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "EXAT", exp}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("value is not an integer or out of range")), response)
}

func TestStorageEXAT3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "EXAT", strconv.Itoa(exp)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("invalid expire time in 'set' command")), response)
}

func TestStoragePXAT1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", time.Now().Add(3*time.Second).UnixMilli()
	args := []string{"SET", key, value, "PXAT", strconv.FormatInt(exp, 10)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleString("OK")), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToBulkString(value)), response)

	// GET after expiration
	time.Sleep(time.Duration(exp-time.Now().UnixMilli()) * time.Millisecond)
	client.Write([]byte(ToArray(args)))
	response = readBuffer(client)
	assert.Equal(t, []byte(ToNull()), response)
}

func TestStoragePXAT2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "PXAT", exp}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("value is not an integer or out of range")), response)
}

func TestStoragePXAT3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "PXAT", strconv.Itoa(exp)}

	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("invalid expire time in 'set' command")), response)
}

func TestEXISTS1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	Set := func(key string, value string) {
		args := []string{"SET", key, value}
		client.Write([]byte(ToArray(args)))
		response := readBuffer(client)
		assert.Equal(t, []byte(ToSimpleString("OK")), response)
	}

	db := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	keys := make([]string, 0, len(db))
	for k, v := range db {
		Set(k, v)
		keys = append(keys, k)
	}

	// exists
	args := []string{"EXISTS"}
	args = append(args, keys...)
	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToInteger(3)), response)
}

func TestEXISTS2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	Set := func(key string, value string) {
		args := []string{"SET", key, value}
		client.Write([]byte(ToArray(args)))
		response := readBuffer(client)
		assert.Equal(t, []byte(ToSimpleString("OK")), response)
	}

	db := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range db {
		Set(k, v)
	}

	// exists
	args := []string{"EXISTS", "randomkey1", "randomkey2"}
	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToInteger(0)), response)
}

func TestEXISTS3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"EXISTS"}
	client.Write([]byte(ToArray(args)))
	response := readBuffer(client)
	assert.Equal(t, []byte(ToSimpleError("wrong number of arguments for 'EXISTS' command")), response)
}
