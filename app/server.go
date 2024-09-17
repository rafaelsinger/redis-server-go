package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
		parseRESP(buffer, conn)
	}
}

func parseRESP(request []byte, conn net.Conn) {
	// RESP arrays format: *<number-of-elements>\r\n<element-1>...<element-n>

	// initialize using conn?
	reader := bufio.NewReader(bytes.NewReader(request))
	for {
		// read until the first delimiter to get the number of elements
		_, err := reader.ReadString(byte('\n'))
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Error parsing request:", err)
			return
		}
		for {
			data, err := reader.ReadString(byte('\n'))
			// remove delimiter
			data = strings.TrimSuffix(data, "\r\n")
			if err != nil {
				break
			}
			if data == "" {
				break
			}
			dataLength, err := strconv.Atoi(data[1:])
			if err != nil {
				fmt.Println("Error processing data length", err)
			}
			switch data[0] {
			case '+':
				// TODO
			case '$':
				parseBulkString(reader, dataLength, conn)
			case ':':
				//TODO
			default:
				fmt.Println("Unknown request type", data)
			}
		}
	}
}

func parseBulkString(reader io.Reader, dataLength int, conn net.Conn) {
	dataBuffer := make([]byte, dataLength)
	n, err := io.ReadFull(reader, dataBuffer)
	if err != nil {
		fmt.Println("Error reading data", err)
	}
	if n < dataLength {
		fmt.Printf("Unable to parse full bulk string. Expected %d, read %d", dataLength, n)
	}
	data := string(dataBuffer)
	switch data {
	case "PING":
		conn.Write([]byte("+PONG\r\n"))
	default:
		fmt.Println("Unknown simple string command")
	}
}
