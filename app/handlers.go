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
// can do so directly from the connection. When closeConnection is true, a
// Connection: close header is added to the response.
func routeRequest(req *Request, reader *bufio.Reader, directory string, closeConnection bool) string {
	extraHeaders := map[string]string{}
	if closeConnection {
		extraHeaders["Connection"] = "close"
	}

	switch {
	case req.Path == "/":
		return simpleResponse(statusOK, extraHeaders)
	case strings.HasPrefix(req.Path, "/echo/"):
		body := strings.TrimPrefix(req.Path, "/echo/")
		return echoResponse(body, req.Headers["accept-encoding"], extraHeaders)
	case req.Path == "/user-agent":
		return textResponse(req.Headers["user-agent"], extraHeaders)
	case req.Method == "POST" && strings.HasPrefix(req.Path, "/files/"):
		filename := strings.TrimPrefix(req.Path, "/files/")
		return saveFileResponse(reader, req.Headers, directory, filename, extraHeaders)
	case strings.HasPrefix(req.Path, "/files/"):
		filename := strings.TrimPrefix(req.Path, "/files/")
		return fileResponse(directory, filename, extraHeaders)
	default:
		return simpleResponse(statusNotFound, extraHeaders)
	}
}

// fileResponse reads filename from directory and returns a 200 response
// with the file contents as an octet-stream, or a 404 if the file doesn't
// exist.
func fileResponse(directory, filename string, extraHeaders map[string]string) string {
	if directory == "" {
		return simpleResponse(statusNotFound, extraHeaders)
	}

	content, err := os.ReadFile(filepath.Join(directory, filename))
	if err != nil {
		return simpleResponse(statusNotFound, extraHeaders)
	}

	return octetStreamResponse(content, extraHeaders)
}

// saveFileResponse reads the request body (sized by the Content-Length
// header) and writes it to filename inside directory, returning a 201
// response on success.
func saveFileResponse(reader *bufio.Reader, headers map[string]string, directory, filename string, extraHeaders map[string]string) string {
	if directory == "" {
		return simpleResponse(statusNotFound, extraHeaders)
	}

	contentLength, err := strconv.Atoi(headers["content-length"])
	if err != nil {
		fmt.Println("Invalid Content-Length: ", err.Error())
		return simpleResponse(statusBadRequest, extraHeaders)
	}

	body := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, body); err != nil {
		fmt.Println("Error reading request body: ", err.Error())
		return simpleResponse(statusBadRequest, extraHeaders)
	}

	if err := os.WriteFile(filepath.Join(directory, filename), body, 0644); err != nil {
		fmt.Println("Error writing file: ", err.Error())
		return simpleResponse(statusServerError, extraHeaders)
	}

	return simpleResponse(statusCreated, extraHeaders)
}
