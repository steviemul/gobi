package handler

import "github.com/steviemul/gobi/internal/server/http"

type IndexRequestHandler struct{}

func (eh IndexRequestHandler) HandleRequest(_ *http.GobiRequest) http.GobiResponse {
	return http.GobiResponse{}
}

func (eh IndexRequestHandler) CanProcess(_ *http.GobiRequest) bool {
	return false
}
