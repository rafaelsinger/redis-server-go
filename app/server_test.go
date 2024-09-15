package app

import (
	"bytes"
	"net"
	"testing"
)

func TestPing(t *testing.T) {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		t.Fatalf("Unable to connect to Redis server: %v", err)
	}
	defer conn.Close()
	conn.Write([]byte("+PING\r\n"))

	var output = make([]byte, 2048)
	num_bytes, err := conn.Read(output)
	if err != nil {
		t.Fatalf("Error reading from Redis: %v", err)
	}
	output = output[:num_bytes]

	expected := "+PONG\r\n"
	if !bytes.Equal(output, []byte(expected)) {
		t.Fatalf("Expected %s, received %s", expected, output)
	}
}
