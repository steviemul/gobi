package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type GobiListener struct {
	listener net.Listener
}

// Listen binds the port and returns once the listener is ready to accept connections.
func (sl *GobiListener) Listen(port string) error {

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return err
	}

	sl.listener = listener

	return nil
}

// Serve accepts connections until the listener is closed. Call Listen first.
func (sl *GobiListener) Serve() error {

	defer sl.listener.Close()

	for {
		conn, err := sl.listener.Accept()

		if err != nil {
			log.Println("accept:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	request, err := ParseRequest(reader)

	if errors.Is(err, ErrMalformedRequest) {
		writeResponse(conn, "400 Bad Request", "Malformed request\n")
		return
	}

	if err != nil {
		log.Println("Failed to read request:", err)
		return
	}

	log.Println("Parsed request:", request.method, request.target, request.version, request.headers)

	writeResponse(conn, "200 OK", "Hello from my Go server!\n")
}

func writeResponse(conn net.Conn, status string, body string) {

	response := fmt.Sprintf(
		"HTTP/1.1 %s\r\n"+
			"Content-Type: text/plain; charset=utf-8\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n"+
			"\r\n"+
			"%s",
		status,
		len(body),
		body,
	)

	if _, err := io.WriteString(conn, response); err != nil {
		log.Println("Failed to write response", err)
	}
}
