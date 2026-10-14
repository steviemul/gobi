package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"time"

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
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			log.Println("accept:", err)
			time.Sleep(50 * time.Millisecond)
			continue
		}

		go handleConnection(conn)
	}
}

const (
	idleTimeout  = 60 * time.Second // max wait between requests
	readTimeout  = 10 * time.Second // max time to receive one request
	writeTimeout = 10 * time.Second
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// Wait for the first byte of the next request. This returns immediately
		// if a pipelined request is already buffered, and fails cleanly on EOF
		// or idle timeout, so neither becomes a 400.
		conn.SetReadDeadline(time.Now().Add(idleTimeout))
		if _, err := reader.Peek(1); err != nil {
			return
		}

		// A request has started, so it must arrive promptly.
		conn.SetReadDeadline(time.Now().Add(readTimeout))

		request, err := http.ParseRequest(reader)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				return // slow client: just drop it
			}
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				return // client went away mid-request: nobody to reply to
			}
			resp := errorResponse(err)
			resp.SetHeader("Connection", "close")
			writeResponse(conn, writer, &resp)
			return
		}

		resp, err := handler.HandleRequest(&request)
		if err != nil {
			resp = handlerErrorResponse(err)
		}

		keepAlive := shouldKeepAlive(&request)
		if !keepAlive {
			resp.SetHeader("Connection", "close")
		} else if request.Version() == "HTTP/1.0" {
			resp.SetHeader("Connection", "keep-alive")
		}

		if err := writeResponse(conn, writer, &resp); err != nil || !keepAlive {
			return
		}
	}
}

func handlerErrorResponse(err error) http.GobiResponse {
	var resp http.GobiResponse
	if errors.Is(err, handler.ErrNotFound) {
		resp, _ = http.BuildResponse("1.1", "404 Not Found", "Not Found\n")
	} else {
		log.Println("handler:", err)
		resp, _ = http.BuildResponse("1.1", "500 Internal Server Error", "Internal server error\n")
	}
	return resp
}

func errorResponse(err error) http.GobiResponse {
	var resp http.GobiResponse
	switch {
	case errors.Is(err, http.ErrMalformedRequest):
		resp, _ = http.BuildResponse("1.1", "400 Bad Request", "Malformed request\n\n"+err.Error()+"\n")
	case errors.Is(err, http.ErrBodyTooLarge):
		resp, _ = http.BuildResponse("1.1", "413 Content Too Large", "Request body too large\n")
	default:
		log.Println("failed to read request:", err)
		resp, _ = http.BuildResponse("1.1", "500 Internal Server Error", "Internal server error\n")
	}
	return resp
}

func writeResponse(conn net.Conn, w *bufio.Writer, resp *http.GobiResponse) error {
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if _, err := w.Write(resp.Bytes()); err != nil {
		return err
	}
	return w.Flush()
}

func shouldKeepAlive(req *http.GobiRequest) bool {
	header := req.Headers()["connection"] // names are lowercased by the parser

	switch req.Version() {
	case "HTTP/1.1":
		return !hasToken(header, "close")
	case "HTTP/1.0":
		return hasToken(header, "keep-alive")
	default:
		return false
	}
}

func hasToken(header, token string) bool {
	for _, part := range strings.Split(header, ",") {
		if strings.EqualFold(strings.TrimSpace(part), token) {
			return true
		}
	}
	return false
}
