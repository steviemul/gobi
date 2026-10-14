package http

import (
	"bytes"
	"fmt"
	"strconv"
)

type GobiResponse struct {
	status  string
	version string
	body    []byte
	headers map[string]string
}

// Bytes renders the response in HTTP/1.x wire format. The body is appended
// as raw bytes, so binary content is preserved.
func (gr *GobiResponse) Bytes() []byte {

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "HTTP/%s %s\r\n", gr.version, gr.status)

	for name, value := range gr.headers {
		fmt.Fprintf(&buf, "%s: %s\r\n", name, value)
	}

	buf.WriteString("\r\n")
	buf.Write(gr.body)

	return buf.Bytes()
}

// String renders the response as a string. Intended for logging and tests;
// use Bytes when writing to a connection.
func (gr *GobiResponse) String() string {
	return string(gr.Bytes())
}

func (gr *GobiResponse) SetHeader(name string, value string) {
	gr.headers[name] = value
}

func BuildResponse(version string, status string, body string) (GobiResponse, error) {
	return BuildContentResponse(version, status, "text/plain; charset=utf-8", []byte(body))
}

// BuildContentResponse is BuildResponse with an explicit Content-Type.
func BuildContentResponse(version string, status string, contentType string, body []byte) (GobiResponse, error) {

	headers := make(map[string]string)

	headers["Content-Type"] = contentType
	headers["Content-Length"] = strconv.Itoa(len(body))

	return GobiResponse{
		version: version,
		status:  status,
		headers: headers,
		body:    body,
	}, nil
}
