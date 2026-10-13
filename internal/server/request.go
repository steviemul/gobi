package server

import (
	"bufio"
	"errors"
	"strings"
)

var ErrMalformedRequest = errors.New("malformed request")

type GobiRequest struct {
	method  string
	target  string
	version string
	headers map[string]string
}

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

	headers, err := parseHeaders(reader)

	if err != nil {
		return GobiRequest{}, err
	}

	return GobiRequest{
		method:  parts[0],
		target:  parts[1],
		version: parts[2],
		headers: headers,
	}, nil
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
