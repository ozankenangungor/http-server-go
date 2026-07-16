package main

import "fmt"

const (
	statusOK          = "HTTP/1.1 200 OK\r\n\r\n"
	statusCreated     = "HTTP/1.1 201 Created\r\n\r\n"
	statusBadRequest  = "HTTP/1.1 400 Bad Request\r\n\r\n"
	statusNotFound    = "HTTP/1.1 404 Not Found\r\n\r\n"
	statusServerError = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
)

// textResponse builds a 200 OK response with a text/plain body.
func textResponse(body string) string {
	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		len(body), body,
	)
}

// echoResponse builds the response for the /echo/{str} endpoint. If the
// client's Accept-Encoding header includes a scheme we support, a matching
// Content-Encoding header is added to the response.
func echoResponse(body, acceptEncoding string) string {
	encoding := negotiateEncoding(acceptEncoding)
	if encoding == "" {
		return textResponse(body)
	}

	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: %s\r\nContent-Length: %d\r\n\r\n%s",
		encoding, len(body), body,
	)
}

// octetStreamResponse builds a 200 OK response with an
// application/octet-stream body, used for serving raw file contents.
func octetStreamResponse(content []byte) string {
	header := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n",
		len(content),
	)
	return header + string(content)
}
