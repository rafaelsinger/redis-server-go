package app

import (
	"fmt"
	"net"
)

func StartServer() {
	ln, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error starting TCP connection")
	}
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go func(c net.Conn) {
			conn.Write([]byte("+PONG\r\n"))
			c.Close()
		}(conn)
	}
}
