package main

import (
	"bufio"
	"errors"
	"strings"
)

// Request represents a parsed HTTP request line and headers. The body, if
// any, is intentionally left unread here since only a few endpoints need
// it; callers that need it read it directly from the buffered reader using
// the Content-Length header.
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
}

// parseRequest reads the request line and headers from reader.
func parseRequest(reader *bufio.Reader) (*Request, error) {
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Request line format: METHOD PATH HTTP/1.1
	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		return nil, errors.New("malformed request line: " + requestLine)
	}

	headers, err := readHeaders(reader)
	if err != nil {
		return nil, err
	}

	return &Request{
		Method:  parts[0],
		Path:    parts[1],
		Headers: headers,
	}, nil
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
