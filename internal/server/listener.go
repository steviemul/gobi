package server

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/steviemul/gobi/internal/config"
	"github.com/steviemul/gobi/internal/server/handler"
	"github.com/steviemul/gobi/internal/server/http"
)

const ProtocolTCP = "tcp"

type GobiListener struct {
	listeners         []net.Listener
	config            config.GobiConfig
	requestDispatcher handler.RequestDispatcher
}

func createTLSListener(config config.GobiConfig) net.Listener {
	tcpListener, err := net.Listen("tcp", ":"+config.TLS().Port)
	if err != nil {
		log.Printf("error creating tls listener - %s", err.Error())
		return nil
	}

	cert, err := tls.LoadX509KeyPair(
		config.TLS().Cert,
		config.TLS().Key,
	)
	if err != nil {
		log.Printf("error reading certificates - %s", err.Error())
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return tls.NewListener(tcpListener, tlsConfig)
}

// Listen binds the port and returns once the listener is ready to accept connections.
func (sl *GobiListener) Listen(config config.GobiConfig) error {
	listener, err := net.Listen(ProtocolTCP, ":"+config.Port())
	if err != nil {
		return err
	}

	sl.listeners = append(sl.listeners, listener)

	if config.TLS().Enabled {
		sl.listeners = append(sl.listeners, createTLSListener(config))
	}

	sl.config = config
	sl.requestDispatcher = handler.BuildRequestDispatcher(config)

	return nil
}

func (sl *GobiListener) Shutdown() {
	for _, listener := range sl.listeners {
		_ = listener.Close()
	}
}

func (sl *GobiListener) poll(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {

			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			log.Println("accept:", err)
			time.Sleep(50 * time.Millisecond)

			continue
		}

		go sl.handleConnection(conn)
	}
}

// Serve accepts connections until the listener is closed. Call Listen first.
func (sl *GobiListener) Serve() {
	defer sl.Shutdown()

	var wg sync.WaitGroup

	for _, listener := range sl.listeners {
		wg.Add(1)
		go func(l net.Listener) {
			defer wg.Done()
			_ = sl.poll(l)
		}(listener)
	}

	wg.Wait()
}

func (sl *GobiListener) handleConnection(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	config := sl.config
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {

		// Wait for the first byte of the next request. This returns immediately
		// if a pipelined request is already buffered, and fails cleanly on EOF
		// or idle timeout, so neither becomes a 400.
		if err := conn.SetReadDeadline(time.Now().Add(config.IdleTimeout())); err != nil {
			return
		}

		if _, err := reader.Peek(1); err != nil {
			return
		}

		// A request has started, so it must arrive promptly.
		if err := conn.SetReadDeadline(time.Now().Add(config.ReadTimeout())); err != nil {
			return
		}

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

			_ = writeResponse(conn, writer, &resp, config) // closing anyway

			return
		}

		resp, err := sl.requestDispatcher.HandleRequest(&request)
		if err != nil {
			resp = handlerErrorResponse(err)
		}

		keepAlive := shouldKeepAlive(&request)

		if !keepAlive {
			resp.SetHeader("Connection", "close")
		} else if request.Version() == "HTTP/1.0" {
			resp.SetHeader("Connection", "keep-alive")
		}

		if err := writeResponse(conn, writer, &resp, config); err != nil || !keepAlive {
			return
		}
	}
}

func writeResponse(conn net.Conn, w *bufio.Writer, resp *http.GobiResponse, config config.GobiConfig) error {
	if err := conn.SetWriteDeadline(time.Now().Add(config.WriteTimeout())); err != nil {
		return err
	}

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
