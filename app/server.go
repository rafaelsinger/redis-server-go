package app

import (
	"fmt"
	"io"
	"net"
	"os"
)

func StartServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error starting TCP connection:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 6379...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed connection")
				return
			}
			fmt.Println("Error reading connection:", err)
			return
		}
		buffer = buffer[:n]
		fmt.Printf("Received %s", string(buffer))

		// r := bytes.NewReader(buffer)
		// TODO: read bytes to handle multiple requests from the same connection

		conn.Write([]byte("+PONG\r\n"))
	}
}
