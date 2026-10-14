package http

import (
	"fmt"
	"strconv"
)

type GobiResponse struct {
	status  string
	version string
	body    string
	headers map[string]string
}

// String renders the response in HTTP/1.x wire format.
func (gr *GobiResponse) String() string {

	response := fmt.Sprintf("HTTP/%s %s\r\n", gr.version, gr.status)

	for name, value := range gr.headers {
		response = response + fmt.Sprintf("%s: %s\r\n", name, value)
	}

	response = response + fmt.Sprintf("\r\n%s", gr.body)

	return response
}

func BuildResponse(version string, status string, body string) (GobiResponse, error) {

	return BuildContentResponse(version, status, "text/plain; charset=utf-8", body)
}

// BuildContentResponse is BuildResponse with an explicit Content-Type.
func BuildContentResponse(version string, status string, contentType string, body string) (GobiResponse, error) {

	headers := make(map[string]string)

	headers["Content-Type"] = contentType
	headers["Content-Length"] = strconv.Itoa(len(body))
	headers["Connection"] = "close"

	return GobiResponse{
		version: version,
		status:  status,
		headers: headers,
		body:    body,
	}, nil
}
