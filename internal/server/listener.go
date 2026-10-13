package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"github.com/steviemul/gobi/internal/server/http"
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

	request, err := http.ParseRequest(reader)

	if errors.Is(err, http.ErrMalformedRequest) {
		writeResponse(conn, "400 Bad Request", "Malformed request\n\n"+err.Error()+"\n")
		return
	}

	if errors.Is(err, http.ErrBodyTooLarge) {
		writeResponse(conn, "413 Content Too Large", "Request body too large\n")
		return
	}

	if err != nil {
		log.Println("Failed to read request:", err)
		return
	}

	log.Println("Parsed request:", request.Method(), request.Target(), request.Version(), request.Headers(), request.Body())

	writeResponse(conn, "200 OK", "Hello from my Go server!\n")
}

func writeResponse(conn net.Conn, status string, body string) {

	response, _ := http.BuildResponse("1.1", status, body)

	if _, err := io.WriteString(conn, response.String()); err != nil {
		log.Println("Failed to write response", err)
	}
}
