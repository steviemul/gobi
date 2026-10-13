package http

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	CONTENT_LENGTH = "content-length" // header names are lowercased by parseHeaders
	MAX_BODY_SIZE  = 1 << 20          // 1 MiB
)

var ErrMalformedRequest = errors.New("malformed request")
var ErrBodyTooLarge = errors.New("request body too large")

type GobiRequest struct {
	method  string
	target  string
	version string
	body    string
	headers map[string]string
}

func (gr GobiRequest) Method() string { return gr.method }

func (gr GobiRequest) Target() string { return gr.target }

func (gr GobiRequest) Version() string { return gr.version }

func (gr GobiRequest) Body() string { return gr.body }

// Headers returns the request headers, keyed by lowercased name.
func (gr GobiRequest) Headers() map[string]string { return gr.headers }

// ParseRequest reads the request line and headers from reader. Header names are lowercased.
// It returns ErrMalformedRequest if the request line or a header is not well formed.
func ParseRequest(reader *bufio.Reader) (GobiRequest, error) {

	requestLine, err := reader.ReadString('\n')

	if err != nil {
		return GobiRequest{}, err
	}

	parts := strings.Fields(requestLine)

	if len(parts) != 3 {
		return GobiRequest{}, ErrMalformedRequest
	}

	if !IsSupportedMethod(parts[0]) {
		return GobiRequest{}, errors.New("invalid request method " + parts[0])
	}

	headers, err := parseHeaders(reader)

	if err != nil {
		return GobiRequest{}, err
	}

	body, err := parseBody(reader, headers)

	if err != nil {
		return GobiRequest{}, err
	}

	return GobiRequest{
		method:  parts[0],
		target:  parts[1],
		version: parts[2],
		body:    body,
		headers: headers,
	}, nil
}

// parseBody reads exactly Content-Length bytes from reader. A request without the header has no body.
// It returns ErrMalformedRequest for an invalid length and ErrBodyTooLarge above MAX_BODY_SIZE.
func parseBody(reader *bufio.Reader, headers map[string]string) (string, error) {

	value, ok := headers[CONTENT_LENGTH]

	if !ok {
		return "", nil
	}

	length, err := strconv.ParseUint(value, 10, 64)

	if err != nil {
		return "", ErrMalformedRequest
	}

	if length > MAX_BODY_SIZE {
		return "", ErrBodyTooLarge
	}

	body := make([]byte, length)

	if _, err := io.ReadFull(reader, body); err != nil {
		return "", err
	}

	return string(body), nil
}

// parseHeaders reads header lines from reader up to the blank line that ends them.
// Header names are lowercased. It returns ErrMalformedRequest if a line has no colon.
func parseHeaders(reader *bufio.Reader) (map[string]string, error) {

	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			return headers, nil
		}

		name, value, found := strings.Cut(line, ":")

		if !found {
			return nil, ErrMalformedRequest
		}

		headers[strings.ToLower(strings.TrimSpace(name))] = strings.TrimSpace(value)
	}
}
