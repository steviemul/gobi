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

		if route.Index != nil && *route.Index {
			handlers = addIndexHandler(route, handlers)

			continue
		}

		if route.Root != "" {
			handlers = addFileHandler(route, handlers)

			continue
		}

		if route.Proxy != "" {
			handlers = addProxyHandler(route, handlers)

			continue
		}

		log.Printf("warning: route [%s] has no root, proxy or index, ignoring it", route.Path)
	}

	if config.Status() {
		handlers = append(handlers, GobiRequestHandler{})
		log.Printf("handler created status handler on")

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

func addIndexHandler(route config.Route, handlers []RequestHandler) []RequestHandler {
	handlers = append(handlers, IndexRequestHandler{
		root: route.Root,
		path: route.Path,
		host: route.Host,
		html: true,
	})

	log.Printf("handler created index handler on [%s]", route.Path)

	return handlers
}

func addFileHandler(route config.Route, handlers []RequestHandler) []RequestHandler {
	handlers = append(handlers, FileRequestHandler{
		root: route.Root,
		path: route.Path,
		host: route.Host,
	})

	log.Printf("handler created file handler on [%s]", route.Path)

	return handlers
}

func addProxyHandler(route config.Route, handlers []RequestHandler) []RequestHandler {
	proxyHandler, err := newProxyRequestHandler(route)
	if err != nil {
		log.Printf("invalid proxy for route [%s]: %s, ignoring it", route.Path, err)

		return handlers
	}

	handlers = append(handlers, proxyHandler)

	log.Printf("handler created proxy handler on [%s]", route.Path)

	return handlers
}
