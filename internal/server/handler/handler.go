package handler

import "github.com/steviemul/gobi/internal/server/http"

type RequestHandler interface {
	HandleRequest(r *http.GobiRequest) http.GobiResponse
}
