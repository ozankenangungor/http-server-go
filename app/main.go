package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	directory := flag.String("directory", "", "directory to serve files from")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn, *directory)
	}
}

// handleConnection reads a single HTTP request from conn and writes the
// matching response. Each connection is handled independently so the
// server can serve multiple clients concurrently.
func handleConnection(conn net.Conn, directory string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Request line format: METHOD PATH HTTP/1.1
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}
	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		fmt.Println("Malformed request line: ", requestLine)
		return
	}
	path := parts[1]

	headers, err := readHeaders(reader)
	if err != nil {
		fmt.Println("Error reading headers: ", err.Error())
		return
	}

	var response string
	switch {
	case path == "/":
		response = "HTTP/1.1 200 OK\r\n\r\n"
	case strings.HasPrefix(path, "/echo/"):
		body := strings.TrimPrefix(path, "/echo/")
		response = textResponse(body)
	case path == "/user-agent":
		response = textResponse(headers["user-agent"])
	case strings.HasPrefix(path, "/files/"):
		filename := strings.TrimPrefix(path, "/files/")
		response = fileResponse(directory, filename)
	default:
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
}

// readHeaders reads HTTP headers until the blank line that ends the header
// section. Header names are stored lowercased, since header names are
// case-insensitive.
func readHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		headers[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
	}
	return headers, nil
}

// textResponse builds a 200 OK response with a text/plain body.
func textResponse(body string) string {
	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		len(body), body,
	)
}

// fileResponse reads the given filename from directory and returns a 200 OK
// response with the file contents as an octet-stream, or a 404 response if
// the file doesn't exist.
func fileResponse(directory, filename string) string {
	if directory == "" {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	content, err := os.ReadFile(filepath.Join(directory, filename))
	if err != nil {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	header := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n",
		len(content),
	)
	return header + string(content)
}
