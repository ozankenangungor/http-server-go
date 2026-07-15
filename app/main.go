package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
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
	method := parts[0]
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
	case method == "POST" && strings.HasPrefix(path, "/files/"):
		filename := strings.TrimPrefix(path, "/files/")
		response = saveFileResponse(reader, headers, directory, filename)
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

// saveFileResponse reads the request body (sized by the Content-Length
// header) and writes it to filename inside directory, returning a 201
// response on success or a 404 if there's no directory configured.
func saveFileResponse(reader *bufio.Reader, headers map[string]string, directory, filename string) string {
	if directory == "" {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	contentLength, err := strconv.Atoi(headers["content-length"])
	if err != nil {
		fmt.Println("Invalid Content-Length: ", err.Error())
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}

	body := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, body); err != nil {
		fmt.Println("Error reading request body: ", err.Error())
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}

	if err := os.WriteFile(filepath.Join(directory, filename), body, 0644); err != nil {
		fmt.Println("Error writing file: ", err.Error())
		return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
	}

	return "HTTP/1.1 201 Created\r\n\r\n"
}
