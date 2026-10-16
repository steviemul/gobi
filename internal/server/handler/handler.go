package handler

import (
	"errors"
	"log"

	"github.com/steviemul/gobi/internal/config"
	"github.com/steviemul/gobi/internal/server/http"
)

var ErrNotFound = errors.New("resource not found")

type RequestDispatcher struct {
	config   config.GobiConfig
	handlers []RequestHandler
}

func createRequestHandlers(config config.GobiConfig) []RequestHandler {
	var handlers []RequestHandler

	for _, route := range config.Routes() {

		if route.Root != "" && route.Proxy != "" {
			log.Printf("invalid route definition for %s, cannot contain both root and proxy", route.Path)
			continue
		}

		if route.Root != "" {
			handlers = append(handlers, FileRequestHandler{
				root: route.Root,
				path: route.Path,
				host: route.Host,
			})
		}

		if route.Proxy != "" {
			handlers = append(handlers, ProxyRequestHandler{
				proxy: route.Proxy,
				path:  route.Path,
				host:  route.Host,
			})
		}
	}

	if config.Status() {
		handlers = append(handlers, GobiRequestHandler{})
	}

	return handlers
}

func BuildRequestDispatcher(config config.GobiConfig) RequestDispatcher {
	handlers := createRequestHandlers(config)

	return RequestDispatcher{
		config:   config,
		handlers: handlers,
	}
}

func (rd RequestDispatcher) HandleRequest(r *http.GobiRequest) (http.GobiResponse, error) {
	for _, handler := range rd.handlers {
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
