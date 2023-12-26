# Redis Server Lite
A simplified Redis server implementation written in Golang. Solution for [Write Your Own Redis Server | Coding Challenges](https://codingchallenges.fyi/challenges/challenge-redis)

## Supported Commands
- PING
- ECHO
- SET
- GET

## How It Works
The server listens for clients. Once a client connects, a go routine is created to handle the client's session. During a session the client can send commands with RESP (Redis serialization protocol) messages. The server parses a message and then sends a response in the same RESP message format.
