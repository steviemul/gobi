package handler

import "github.com/steviemul/gobi/internal/server/http"

type ProxyRequestHandler struct {
	proxy string
	host  string
	path  string
}

func (eh ProxyRequestHandler) HandleRequest(_ *http.GobiRequest) http.GobiResponse {
	return http.GobiResponse{}
}

func (eh ProxyRequestHandler) CanProcess(_ *http.GobiRequest) bool {
	return false
}
