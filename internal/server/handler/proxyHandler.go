package handler

import (
	"errors"
	"io"
	"log"
	nethttp "net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/steviemul/gobi/internal/config"
	"github.com/steviemul/gobi/internal/server/http"
)

// hopByHopResponseHeaders are connection-level headers that must not be copied
// from the backend response; the listener manages the client connection itself.
var hopByHopResponseHeaders = map[string]struct{}{
	"Connection":        {},
	"Keep-Alive":        {},
	"Transfer-Encoding": {},
	"Content-Length":    {}, // recomputed from the buffered body
	"Content-Type":      {}, // passed to BuildContentResponse
}

// ProxyRequestHandler forwards requests under path to a backend using
// httputil.ReverseProxy. The route prefix is stripped before forwarding, so
// /proxied-path/foo?x=1 is sent to <proxy>/foo?x=1.
type ProxyRequestHandler struct {
	proxy        string
	host         string
	path         string
	reverseProxy *httputil.ReverseProxy
}

func newProxyRequestHandler(route config.Route) (ProxyRequestHandler, error) {
	target, err := url.Parse(route.Proxy)
	if err != nil {
		return ProxyRequestHandler{}, err
	}

	if (target.Scheme != "http" && target.Scheme != "https") || target.Host == "" {
		return ProxyRequestHandler{}, errors.New("proxy must be an absolute http or https URL")
	}

	reverseProxy := &httputil.ReverseProxy{
		// SetURL points the request at the backend, joins its base path and sets
		// Host to the backend's own, which virtual-hosted backends route on.
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
		},
		ErrorHandler: func(w nethttp.ResponseWriter, _ *nethttp.Request, err error) {
			log.Printf("proxy error for [%s]: %s", route.Path, err)

			w.WriteHeader(nethttp.StatusBadGateway)
		},
	}

	return ProxyRequestHandler{
		proxy: route.Proxy,
		// "/api/" and "/api" are the same route; the boundary check wants no trailing slash.
		path:         strings.TrimSuffix(route.Path, "/"),
		host:         route.Host,
		reverseProxy: reverseProxy,
	}, nil
}

func (eh ProxyRequestHandler) HandleRequest(r *http.GobiRequest) http.GobiResponse {
	req, err := eh.buildRequest(r)
	if err != nil {
		return badGateway()
	}

	recorder := httptest.NewRecorder()

	eh.reverseProxy.ServeHTTP(recorder, req)

	result := recorder.Result()

	defer func() { _ = result.Body.Close() }()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return badGateway()
	}

	response, _ := http.BuildContentResponse("1.1", result.Status, result.Header.Get("Content-Type"), body)

	if result.Header.Get("Content-Type") == "" {
		response.DeleteHeader("Content-Type")
	}

	for name, values := range result.Header {
		if _, skip := hopByHopResponseHeaders[name]; skip {
			continue
		}

		// GobiResponse holds one value per header, so repeated headers are folded.
		response.SetHeader(name, strings.Join(values, ", "))
	}

	return response
}

func (eh ProxyRequestHandler) CanProcess(r *http.GobiRequest) bool {
	if eh.host != "" && !matchesHost(eh.host, r) {
		return false
	}

	_, ok := pathResource(eh.path, r.Target())

	return ok
}

// buildRequest converts the parsed gobi request into a net/http request with
// the route prefix removed from its target.
func (eh ProxyRequestHandler) buildRequest(r *http.GobiRequest) (*nethttp.Request, error) {
	target := strings.TrimPrefix(r.Target(), eh.path)

	if !strings.HasPrefix(target, "/") {
		target = "/" + target // "/proxied-path?x=1" or exactly "/proxied-path"
	}

	req, err := nethttp.NewRequest(r.Method(), target, strings.NewReader(r.Body()))
	if err != nil {
		return nil, err
	}

	for name, value := range r.Headers() {
		req.Header.Set(name, value)
	}

	req.Host = r.Headers()["host"]

	return req, nil
}

func badGateway() http.GobiResponse {
	response, _ := http.BuildResponse("1.1", "502 Bad Gateway", "Bad Gateway\n")

	return response
}
