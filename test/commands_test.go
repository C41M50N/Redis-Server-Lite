package test

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/C41M50N/redis-server-lite-go/internal/r"
	"github.com/C41M50N/redis-server-lite-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func createMockConnection() net.Conn {
	server, client := net.Pipe()
	go func() {
		utils.ProcessClient(server)
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

	args := []string{"PEEK"}

	client.Write(r.ToArray(args))
	response := readBuffer(client)

	assert.Equal(t, r.ToSimpleError("unknown command 'PEEK'"), response)
}

func TestPing1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"PING"}

	client.Write(r.ToArray(args))
	response := readBuffer(client)

	assert.Equal(t, r.ToBulkString("PONG"), response)
}

func TestPing2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	message := "anyone there?"
	args := []string{"PING", message}

	client.Write(r.ToArray(args))
	response := readBuffer(client)

	assert.Equal(t, r.ToBulkString(message), response)
}

func TestPing3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"PING", "echo", "chamber"}

	client.Write(r.ToArray(args))
	response := readBuffer(client)

	assert.Equal(t, r.ToSimpleError("wrong number of arguments for 'ping' command"), response)
}

func TestEcho1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	message := "Wubba lubba dub dub!"
	args := []string{"ECHO", message}

	client.Write(r.ToArray(args))
	response := readBuffer(client)

	assert.Equal(t, r.ToBulkString(message), response)
}

func TestEcho2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"ECHO", "Wubba lubba dub dub!", "Morty!!!"}

	client.Write(r.ToArray(args))
	response := readBuffer(client)

	assert.Equal(t, r.ToSimpleError("wrong number of arguments for 'echo' command"), response)
}

func TestStorage1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value := "salary", "123456"
	args := []string{"SET", key, value}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"GET", key}

	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString(value), response)
}

func TestStorage2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value := "description", "thing thing something thing"

	args := []string{"SET", key, value}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"GET", key}

	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString(value), response)
}

func TestStorage3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key := "key"
	args := []string{"SET", key}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("wrong number of arguments for 'set' command"), response)
}

func TestStorageEX1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", 3
	args := []string{"SET", key, value, "EX", strconv.Itoa(exp)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString(value), response)

	// GET after expiration
	time.Sleep(time.Duration(exp) * time.Second)
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToNull(), response)
}

func TestStorageEX2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "EX", exp}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("value is not an integer or out of range"), response)
}

func TestStorageEX3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "EX", strconv.Itoa(exp)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("invalid expire time in 'set' command"), response)
}

func TestStoragePX1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", 3000
	args := []string{"SET", key, value, "PX", strconv.Itoa(exp)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString(value), response)

	// GET after expiration
	time.Sleep(time.Duration(exp) * time.Millisecond)
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToNull(), response)
}

func TestStoragePX2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "PX", exp}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("value is not an integer or out of range"), response)
}

func TestStoragePX3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "PX", strconv.Itoa(exp)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("invalid expire time in 'set' command"), response)
}

func TestStorageEXAT1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", time.Now().Add(3*time.Second).Unix()
	args := []string{"SET", key, value, "EXAT", strconv.FormatInt(exp, 10)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString(value), response)

	// GET after expiration
	time.Sleep(time.Duration(exp-time.Now().Unix()) * time.Second)
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToNull(), response)
}

func TestStorageEXAT2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "EXAT", exp}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("value is not an integer or out of range"), response)
}

func TestStorageEXAT3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "EXAT", strconv.Itoa(exp)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("invalid expire time in 'set' command"), response)
}

func TestStoragePXAT1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", time.Now().Add(3*time.Second).UnixMilli()
	args := []string{"SET", key, value, "PXAT", strconv.FormatInt(exp, 10)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"GET", key}

	// GET immediately
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString(value), response)

	// GET after expiration
	time.Sleep(time.Duration(exp-time.Now().UnixMilli()) * time.Millisecond)
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToNull(), response)
}

func TestStoragePXAT2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", "BAD_EXP"
	args := []string{"SET", key, value, "PXAT", exp}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("value is not an integer or out of range"), response)
}

func TestStoragePXAT3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	key, value, exp := "user:paswd:loggedin", "true", -11111
	args := []string{"SET", key, value, "PXAT", strconv.Itoa(exp)}

	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("invalid expire time in 'set' command"), response)
}

func TestEXISTS1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	Set := func(key string, value string) {
		args := []string{"SET", key, value}
		client.Write(r.ToArray(args))
		response := readBuffer(client)
		assert.Equal(t, r.ToSimpleString("OK"), response)
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
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToInteger(3), response)
}

func TestEXISTS2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	Set := func(key string, value string) {
		args := []string{"SET", key, value}
		client.Write(r.ToArray(args))
		response := readBuffer(client)
		assert.Equal(t, r.ToSimpleString("OK"), response)
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
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToInteger(0), response)
}

func TestEXISTS3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"EXISTS"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("wrong number of arguments for 'EXISTS' command"), response)
}

func TestDEL1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	args := []string{"SET", "key1", "value1"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"SET", "key2", "value2"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	args = []string{"SET", "key3", "value3"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	// del
	args = []string{"DEL", "key0", "key1", "key2"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToInteger(2), response)
}

func TestDEL2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"DEL"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("wrong number of arguments for 'DEL' command"), response)
}

func TestINCR1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	args := []string{"SET", "salary", "123000"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	// incr
	args = []string{"INCR", "salary"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToInteger(123001), response)

	// get
	args = []string{"GET", "salary"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString("123001"), response)
}

func TestINCR2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	args := []string{"SET", "salary", "enough"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	// incr
	args = []string{"INCR", "salary"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToSimpleError("value is not an integer or out of range"), response)
}

func TestINCR3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"INCR", "key1"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToInteger(1), response)
}

func TestINCR4(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"INCR"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("wrong number of arguments for 'INCR' command"), response)
}

func TestDECR1(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	args := []string{"SET", "salary", "123000"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	// incr
	args = []string{"DECR", "salary"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToInteger(122999), response)

	// get
	args = []string{"GET", "salary"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToBulkString("122999"), response)
}

func TestDECR2(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	// set
	args := []string{"SET", "salary", "enough"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleString("OK"), response)

	// incr
	args = []string{"DECR", "salary"}
	client.Write(r.ToArray(args))
	response = readBuffer(client)
	assert.Equal(t, r.ToSimpleError("value is not an integer or out of range"), response)
}

func TestDECR3(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"DECR", "key"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToInteger(-1), response)
}

func TestDECR4(t *testing.T) {
	client := createMockConnection()
	defer client.Close()

	args := []string{"DECR"}
	client.Write(r.ToArray(args))
	response := readBuffer(client)
	assert.Equal(t, r.ToSimpleError("wrong number of arguments for 'DECR' command"), response)
}
