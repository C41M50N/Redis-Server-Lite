# Redis Server Lite
A simplified Redis server implementation written in Golang. This is my solution for [John Crickett's Write Your Own Redis Server Coding Challenge](https://codingchallenges.fyi/challenges/challenge-redis).

## How It Works
The server listens for clients. Once a client connects, a go routine is created to handle the client's session. During a session the client can send commands with RESP (Redis serialization protocol) messages. The server parses a message and then sends a response in the same RESP message format.

## How to Run
```bash
git clone https://github.com/C41M50N/Redis-Server-Lite
cd Redis-Server-Lite
go run cmd/redis-server/server.go
```

## Supported Commands
A list of the server's supported commands and their usage syntax.

### PING
Returns `PONG` if no argument is provided, otherwise return a copy of the argument as a bulk.
```
PING [message]
```

### ECHO
Returns `message`.
```
ECHO message
```

### SET
Set `key` to hold the string `value`. If `key` already holds a value, it is overwritten, regardless of its type.
```
SET key value [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds]
```

### GET
Returns the value of `key`. If the key does not exist the special value `nil` is returned.
```
GET key
```

### EXISTS
Returns the number of keys that exist from those specified as arguments.
```
EXISTS key [key ...]
```

### DEL
Removes the specified keys. A key is ignored if it does not exist. Returns the number of keys that were removed.
```
DEL key [key ...]
```

### INCR
Increments the number stored at `key` by one. If the key does not exist, it is set to `0` before performing the operation. An error is returned if the key contains a value of the wrong type or contains a string that can not be represented as integer. This operation is limited to 64 bit signed integers. Returns the value of the key (as an integer) after the increment.
```
INCR key
```

### DECR
Decrements the number stored at `key` by one. If the key does not exist, it is set to `0` before performing the operation. An error is returned if the key contains a value of the wrong type or contains a string that can not be represented as integer. This operation is limited to 64 bit signed integers. Returns the value of the key (as an integer) after the decrement.
```
DECR key
```

### LPUSH
Insert all the specified values at the head of the list stored at `key`. If `key` does not exist, it is created as empty list before performing the push operations. When `key` holds a value that is not a list, an error is returned. Returns the length of the list after the push operation.
```
LPUSH key element [element ...]
```

### RPUSH
Insert all the specified values at the tail of the list stored at `key`. If `key` does not exist, it is created as empty list before performing the push operation. When `key` holds a value that is not a list, an error is returned. Returns the length of the list after the push operation.
```
RPUSH key element [element ...]
```

### LRANGE
Returns the specified elements of the list stored at `key`. If the data type of the value at the given `key` is not a list, then an error is returned. The offsets `start` and `stop` are zero-based indexes, with `0` being the first element of the list (the head of the list), `1` being the next element and so on. The offsets can also be negative numbers indicating offsets starting at the end of the list.
```
LRANGE key start stop
```

## Benchmarks
The following benchmarks were performed on my M2 MacBook Pro.

### Custom Redis
```bash
> redis-benchmark -t set,get -n 100000 -q
SET: 42735.04 requests per second, p50=0.687 msec
GET: 46019.32 requests per second, p50=0.615 msec

> redis-benchmark -t set,get -n 100000 -q
SET: 44014.08 requests per second, p50=0.695 msec
GET: 46210.72 requests per second, p50=0.607 msec

> redis-benchmark -t set,get -n 100000 -q
SET: 44503.79 requests per second, p50=0.663 msec
GET: 47169.81 requests per second, p50=0.591 msec
```

### Production Redis
```bash
> redis-benchmark -t set,get -n 100000 -q
SET: 165289.25 requests per second, p50=0.143 msec
GET: 159489.64 requests per second, p50=0.143 msec

> redis-benchmark -t set,get -n 100000 -q
SET: 163934.42 requests per second, p50=0.143 msec
GET: 165016.50 requests per second, p50=0.143 msec

> redis-benchmark -t set,get -n 100000 -q
SET: 171821.30 requests per second, p50=0.143 msec
GET: 168350.17 requests per second, p50=0.143 msec
```
