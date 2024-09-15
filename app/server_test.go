package app

import (
	"bytes"
	"net"
	"testing"
)

func sendRedisCommand(t *testing.T, commands []string) []byte {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		t.Fatalf("Unable to connect to Redis server: %v", err)
	}
	defer conn.Close()
	for _, command := range commands {
		conn.Write([]byte(command + "\r\n"))
	}

	var output = make([]byte, 2048)
	num_bytes, err := conn.Read(output)
	if err != nil {
		t.Fatalf("Error reading from Redis: %v", err)
	}
	return output[:num_bytes]
}

func testBytesEquality(t *testing.T, expected, received []byte) {
	if !bytes.Equal(expected, received) {
		t.Fatalf("Expected %s, received %s", expected, received)
	}
}

func TestPing(t *testing.T) {
	// single ping
	output := sendRedisCommand(t, []string{"*1\r\n$4\r\nPING\r\n"})
	expected := []byte("+PONG\r\n")
	testBytesEquality(t, expected, output)

	// multiple pings in the same connection
	output = sendRedisCommand(t, []string{"*1\r\n$4\r\nPING\r\n", "*1\r\n$4\r\nPING\r\n"})
	expected = []byte("+PONG\r\n+PONG\r\n")
	testBytesEquality(t, expected, output)
}
