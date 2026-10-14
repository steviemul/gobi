package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"github.com/steviemul/gobi/internal/server/handler"
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

	var response http.GobiResponse

	reader := bufio.NewReader(conn)

	request, err := http.ParseRequest(reader)

	if errors.Is(err, http.ErrMalformedRequest) {
		response, _ = buildResponse("400 Bad Request", "Malformed request\n\n"+err.Error()+"\n")
	}

	if errors.Is(err, http.ErrBodyTooLarge) {
		response, _ = buildResponse("413 Content Too Large", "Request body too large\n")
	}

	if err != nil {
		response, _ = buildResponse("500 Internal Server Error", "Failed to read request:\n")
		log.Println("Failed to read request:", err)
		return
	}

	log.Println("Parsed request:", request.Method(), request.Target(), request.Version(), request.Headers(), request.Body())

	response, err = handler.HandleRequest(&request)

	if errors.Is(err, handler.ErrNotFound) {
		response, _ = buildResponse("404 Not Found", "The requested resource could not be found")
	}

	writeResponse(conn, &response)
}

func writeResponse(conn net.Conn, response *http.GobiResponse) {
	if _, err := io.WriteString(conn, response.String()); err != nil {
		log.Println("Failed to write response", err)
	}
}

func buildResponse(status string, body string) (http.GobiResponse, error) {
	return http.BuildResponse(
		"1.1",
		status,
		body,
	)
}
