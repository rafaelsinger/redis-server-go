package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error starting TCP connection")
	}
	conn, _ := ln.Accept()
	for {
		conn.Write([]byte("$4\r\nPONG\r\n"))
	}
}
