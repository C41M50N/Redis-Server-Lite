package main

import (
	"fmt"
	"net"
	"os"

	"github.com/C41M50N/Redis-Server-Lite/internal/utils"
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
		go utils.ProcessClient(conn)
	}
}
