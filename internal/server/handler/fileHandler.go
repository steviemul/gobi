package handler

import "github.com/steviemul/gobi/internal/server/http"

type FileRequestHandler struct{}

func (eh FileRequestHandler) HandleRequest(_ *http.GobiRequest) http.GobiResponse {
	return http.GobiResponse{}
}

func (eh FileRequestHandler) CanProcess(_ *http.GobiRequest) bool {
	return false
}
