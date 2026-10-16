package handler

import (
	"errors"

	"github.com/steviemul/gobi/internal/config"
	"github.com/steviemul/gobi/internal/server/http"
)

var handlers = []RequestHandler{
	FileRequestHandler{},
	GobiRequestHandler{},
}

var ErrNotFound = errors.New("resource not found")

type RequestDispatcher struct {
	config   config.GobiConfig
	handlers []RequestHandler
}

func BuildRequestDispatcher(config config.GobiConfig) RequestDispatcher {
	return RequestDispatcher{
		config:   config,
		handlers: handlers,
	}
}

func (rd RequestDispatcher) HandleRequest(r *http.GobiRequest) (http.GobiResponse, error) {
	for _, handler := range handlers {
		if handler.CanProcess(r) {
			return handler.HandleRequest(r), nil
		}
	}

	return http.GobiResponse{}, ErrNotFound
}

type RequestHandler interface {
	HandleRequest(r *http.GobiRequest) http.GobiResponse
	CanProcess(r *http.GobiRequest) bool
}
