package server

import (
	"errors"
	"log"

	"github.com/steviemul/gobi/internal/server/handler"
	"github.com/steviemul/gobi/internal/server/http"
)

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
