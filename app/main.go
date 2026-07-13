package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	requestLine, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	// Request line format: METHOD PATH HTTP/1.1
	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		fmt.Println("Malformed request line: ", requestLine)
		os.Exit(1)
	}
	path := parts[1]

	var response string
	switch {
	case path == "/":
		response = "HTTP/1.1 200 OK\r\n\r\n"
	case strings.HasPrefix(path, "/echo/"):
		body := strings.TrimPrefix(path, "/echo/")
		response = fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
			len(body), body,
		)
	default:
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}
