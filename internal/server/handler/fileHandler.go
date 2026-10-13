package handler

import "github.com/steviemul/gobi/internal/server/http"

type FileRequestHandler struct {
}

func (eh FileRequestHandler) HandleRequest(r *http.GobiRequest) http.GobiResponse {
	return http.GobiResponse{}
}
