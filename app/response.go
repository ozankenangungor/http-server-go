package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	statusOK          = "200 OK"
	statusCreated     = "201 Created"
	statusBadRequest  = "400 Bad Request"
	statusNotFound    = "404 Not Found"
	statusServerError = "500 Internal Server Error"
)

// buildResponse serializes an HTTP/1.1 response with the given status line,
// headers and body. Header order isn't significant to HTTP clients, but we
// sort names for stable, easy-to-read output.
func buildResponse(status string, headers map[string]string, body []byte) string {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}
	sort.Strings(names)

	var b strings.Builder
	b.WriteString("HTTP/1.1 ")
	b.WriteString(status)
	b.WriteString("\r\n")
	for _, name := range names {
		b.WriteString(name)
		b.WriteString(": ")
		b.WriteString(headers[name])
		b.WriteString("\r\n")
	}
	b.WriteString("\r\n")
	b.Write(body)

	return b.String()
}

// simpleResponse builds a response with no body, carrying only the given
// extra headers (e.g. Connection: close) alongside the status line.
func simpleResponse(status string, extraHeaders map[string]string) string {
	return buildResponse(status, extraHeaders, nil)
}

// textResponse builds a 200 OK response with a text/plain body.
func textResponse(body string, extraHeaders map[string]string) string {
	headers := mergeHeaders(extraHeaders, map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(body)),
	})
	return buildResponse(statusOK, headers, []byte(body))
}

// echoResponse builds the response for the /echo/{str} endpoint. If the
// client's Accept-Encoding header includes a scheme we support, the body is
// compressed with that scheme and a matching Content-Encoding header is
// added to the response.
func echoResponse(body, acceptEncoding string, extraHeaders map[string]string) string {
	encoding := negotiateEncoding(acceptEncoding)
	if encoding == "" {
		return textResponse(body, extraHeaders)
	}

	// negotiateEncoding only ever returns schemes we support, and gzip is
	// currently the only one, so this is the sole case for now.
	compressed, err := gzipCompress([]byte(body))
	if err != nil {
		fmt.Println("Error compressing body: ", err.Error())
		return textResponse(body, extraHeaders)
	}

	headers := mergeHeaders(extraHeaders, map[string]string{
		"Content-Type":     "text/plain",
		"Content-Encoding": encoding,
		"Content-Length":   strconv.Itoa(len(compressed)),
	})
	return buildResponse(statusOK, headers, compressed)
}

// octetStreamResponse builds a 200 OK response with an
// application/octet-stream body, used for serving raw file contents.
func octetStreamResponse(content []byte, extraHeaders map[string]string) string {
	headers := mergeHeaders(extraHeaders, map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": strconv.Itoa(len(content)),
	})
	return buildResponse(statusOK, headers, content)
}

// mergeHeaders returns a new map containing both extra and base, with base
// taking precedence on key collisions.
func mergeHeaders(extra, base map[string]string) map[string]string {
	merged := make(map[string]string, len(extra)+len(base))
	for k, v := range extra {
		merged[k] = v
	}
	for k, v := range base {
		merged[k] = v
	}
	return merged
}
