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
	reader := bufio.NewReader(bytes.NewReader(request))
	for {
		// get number of elements in array
		numElementString, err := reader.ReadString(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error parsing request:", err)
			}
			return
		}
		if numElementString[0] != '*' {
			fmt.Println("Error: expected array")
			return
		}
		numElementString = strings.TrimSuffix(numElementString, "\r\n")
		numElements, err := strconv.Atoi(numElementString[1:])
		if err != nil {
			fmt.Println("Error parsing number of elmements:", err)
			return
		}
		// read over each element
		for i := 0; i < numElements; i++ {
			dataInfoString, err := reader.ReadString(byte('\n'))
			if err != nil {
				fmt.Println("Error reading data", err)
				break // don't return since we can still read other valid elements
			}
			// remove delimiter
			dataInfoString = strings.TrimSuffix(dataInfoString, "\r\n")
			dataLength, err := strconv.Atoi(dataInfoString[1:])
			if err != nil {
				fmt.Println("Error processing data length", err)
				break
			}
			dataType := dataInfoString[0]
			switch dataType {
			case '+':
				// TODO
			case '$':
				parseBulkString(reader, dataLength, conn)
			case ':':
				//TODO
			default:
				fmt.Println("Unknown RESP data type:", dataType)
			}
		}
		// read final crlf delimiter
		readCRLF(reader)
	}
}

func parseBulkString(reader io.Reader, dataLength int, conn net.Conn) {
	dataBuffer := make([]byte, dataLength)
	n, err := io.ReadFull(reader, dataBuffer)
	readCRLF(reader)
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

func readCRLF(reader io.Reader) {
	crlf := make([]byte, 2)
	_, err := io.ReadFull(reader, crlf)
	if err != nil {
		fmt.Println("Error parsing crlf delimiter:", err)
		return
	}
}
