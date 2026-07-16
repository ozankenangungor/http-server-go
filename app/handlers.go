package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// routeRequest picks the right handler for req and returns the raw HTTP
// response to send back to the client. reader is passed through so
// handlers that need to read a request body (e.g. POST /files/{filename})
// can do so directly from the connection.
func routeRequest(req *Request, reader *bufio.Reader, directory string) string {
	switch {
	case req.Path == "/":
		return statusOK
	case strings.HasPrefix(req.Path, "/echo/"):
		body := strings.TrimPrefix(req.Path, "/echo/")
		return echoResponse(body, req.Headers["accept-encoding"])
	case req.Path == "/user-agent":
		return textResponse(req.Headers["user-agent"])
	case req.Method == "POST" && strings.HasPrefix(req.Path, "/files/"):
		filename := strings.TrimPrefix(req.Path, "/files/")
		return saveFileResponse(reader, req.Headers, directory, filename)
	case strings.HasPrefix(req.Path, "/files/"):
		filename := strings.TrimPrefix(req.Path, "/files/")
		return fileResponse(directory, filename)
	default:
		return statusNotFound
	}
}

// fileResponse reads filename from directory and returns a 200 response
// with the file contents as an octet-stream, or a 404 if the file doesn't
// exist.
func fileResponse(directory, filename string) string {
	if directory == "" {
		return statusNotFound
	}

	content, err := os.ReadFile(filepath.Join(directory, filename))
	if err != nil {
		return statusNotFound
	}

	return octetStreamResponse(content)
}

// saveFileResponse reads the request body (sized by the Content-Length
// header) and writes it to filename inside directory, returning a 201
// response on success.
func saveFileResponse(reader *bufio.Reader, headers map[string]string, directory, filename string) string {
	if directory == "" {
		return statusNotFound
	}

	contentLength, err := strconv.Atoi(headers["content-length"])
	if err != nil {
		fmt.Println("Invalid Content-Length: ", err.Error())
		return statusBadRequest
	}

	body := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, body); err != nil {
		fmt.Println("Error reading request body: ", err.Error())
		return statusBadRequest
	}

	if err := os.WriteFile(filepath.Join(directory, filename), body, 0644); err != nil {
		fmt.Println("Error writing file: ", err.Error())
		return statusServerError
	}

	return statusCreated
}
