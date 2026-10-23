package handler

import (
	"bufio"
	"fmt"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/steviemul/gobi/internal/config"
	"github.com/steviemul/gobi/internal/server/http"
)

func parseTestRequest(t *testing.T, raw string) *http.GobiRequest {
	t.Helper()

	req, err := http.ParseRequest(bufio.NewReader(strings.NewReader(raw)))
	if err != nil {
		t.Fatal(err)
	}

	return &req
}

func newTestProxy(t *testing.T, backendURL string, route config.Route) ProxyRequestHandler {
	t.Helper()

	route.Proxy = backendURL

	h, err := newProxyRequestHandler(route)
	if err != nil {
		t.Fatal(err)
	}

	return h
}

func TestProxyForwardsRequest(t *testing.T) {
	backend := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		body, _ := io.ReadAll(r.Body)

		w.Header().Set("Content-Type", "text/x-echo")
		w.Header().Set("X-Backend", "yes")
		w.WriteHeader(nethttp.StatusCreated)

		fmt.Fprintf(w, "%s %s host=%s body=%s", r.Method, r.URL.RequestURI(), r.Host, body)
	}))
	defer backend.Close()

	h := newTestProxy(t, backend.URL+"/base", config.Route{Path: "/api"})
	backendHost := strings.TrimPrefix(backend.URL, "http://")

	raw := "POST /api/things?x=1 HTTP/1.1\r\nHost: gobi.local\r\nContent-Length: 5\r\n\r\nhello"
	resp := h.HandleRequest(parseTestRequest(t, raw))
	out := resp.String()

	for _, want := range []string{
		"201 Created",
		"Content-Type: text/x-echo",
		"X-Backend: yes",
		"POST /base/things?x=1 host=" + backendHost + " body=hello",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("response missing %q:\n%s", want, out)
		}
	}
}

func TestProxyBackendDown(t *testing.T) {
	backend := httptest.NewServer(nil)
	url := backend.URL

	backend.Close()

	h := newTestProxy(t, url, config.Route{Path: "/api"})
	resp := h.HandleRequest(parseTestRequest(t, "GET /api/ HTTP/1.1\r\nHost: x\r\n\r\n"))

	if !strings.Contains(resp.String(), "502") {
		t.Errorf("want 502, got:\n%s", resp.String())
	}
}

func TestProxyCanProcess(t *testing.T) {
	h := newTestProxy(t, "http://localhost:1", config.Route{Path: "/api", Host: "example.com"})

	tests := []struct {
		name, raw string
		want      bool
	}{
		{"exact path", "GET /api HTTP/1.1\r\nHost: example.com\r\n\r\n", true},
		{"sub path", "GET /api/x HTTP/1.1\r\nHost: example.com:8080\r\n\r\n", true},
		{"similar prefix", "GET /apix HTTP/1.1\r\nHost: example.com\r\n\r\n", false},
		{"other path", "GET /other HTTP/1.1\r\nHost: example.com\r\n\r\n", false},
		{"other host", "GET /api HTTP/1.1\r\nHost: other.com\r\n\r\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.CanProcess(parseTestRequest(t, tt.raw)); got != tt.want {
				t.Errorf("CanProcess = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProxyRequestHandlerRejectsBadURL(t *testing.T) {
	for _, bad := range []string{"not a url", "ftp://host/", "/relative", "http://"} {
		if _, err := newProxyRequestHandler(config.Route{Path: "/p", Proxy: bad}); err == nil {
			t.Errorf("expected error for proxy %q", bad)
		}
	}
}
