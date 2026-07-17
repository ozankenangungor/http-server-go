# HTTP Server (Go)

This is my solution to CodeCrafters' "Build Your Own HTTP Server" challenge, written in Go.

The idea is simple but you learn a lot from it: instead of reaching for Go's `net/http`, you build an HTTP/1.1 server from scratch on raw TCP sockets. You parse the request line yourself, read headers by hand, handle the body based on `Content-Length`, and write the response bytes back manually. No shortcuts.

I wanted to actually understand what's happening underneath all the abstractions I use every day instead of just calling `http.ListenAndServe` and moving on.

![demo](demo.gif)

## What it supports

- Binding to a TCP port and accepting connections
- Basic routing (`/`, `/echo/{str}`, `/user-agent`, `/files/{filename}`)
- Reading and writing files through `/files/{filename}` (GET to read, POST to write), with the directory configurable via a `--directory` flag
- Handling multiple clients concurrently, one goroutine per connection
- Gzip compression negotiated through `Accept-Encoding` / `Content-Encoding`
- Persistent (keep-alive) connections, and closing them properly when the client sends `Connection: close`

## Project structure

Everything lives in `app/`, split by responsibility instead of one giant file:

- `main.go` – entry point, listener setup, the per-connection loop
- `request.go` – parses the request line and headers
- `response.go` – builds responses (status line, headers, body)
- `handlers.go` – routes requests to the right logic
- `compression.go` – negotiates a compression scheme and gzips the body when needed

## Running it

```sh
./your_program.sh --directory /tmp/
```

Then hit it with curl:

```sh
curl http://localhost:4221/echo/hello
```
