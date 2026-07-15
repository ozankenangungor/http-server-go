package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
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

	req, err := parseRequest(reader)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	response := routeRequest(req, reader, directory)

	if _, err := conn.Write([]byte(response)); err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
}
