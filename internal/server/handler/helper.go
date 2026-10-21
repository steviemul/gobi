package handler

import (
	"net"
	"net/url"
	"strings"

	"github.com/steviemul/gobi/internal/server/http"
)

func pathResource(path string, target string) (resource string, ok bool) {
	return pathResourceWithDefault(path, target, indexFile)
}

func pathResourceWithDefault(path string, target string, defaultPath string) (resource string, ok bool) {
	target, _, _ = strings.Cut(target, "?")
	target, _, _ = strings.Cut(target, "#")

	rest, found := strings.CutPrefix(target, path)

	if !found || (rest != "" && rest[0] != '/') {
		return "", false
	}

	resource = strings.TrimPrefix(rest, "/")

	// Targets are percent-encoded ("a%20b.txt"); handlers need the real name.
	// Callers still validate the decoded path, so "%2e%2e" can't escape the root.
	resource, err := url.PathUnescape(resource)
	if err != nil {
		return "", false
	}

	if resource == "" {
		resource = defaultPath
	}

	return resource, true
}

func matchesHost(host string, r *http.GobiRequest) bool {
	hostHeader := r.Headers()["host"] // names are lowercased by the parser

	return strings.EqualFold(host, hostname(hostHeader))
}

// hostname strips the optional port from a Host header value. net.SplitHostPort
// is used rather than splitting on ":" so bracketed IPv6 literals like [::1]:8080 work.
func hostname(hostport string) string {
	h, _, err := net.SplitHostPort(hostport)
	if err != nil {
		// No port present (or malformed); use the value as-is.
		return hostport
	}

	return h
}

// directoryRedirect reports the trailing-slash form of target when it is exactly the
// route path with no trailing slash. The query string, if any, is kept.
func directoryRedirect(routePath string, target string) (location string, ok bool) {
	if routePath == "" || strings.HasSuffix(routePath, "/") {
		return "", false
	}

	target, _, _ = strings.Cut(target, "#")
	pathPart, query, hasQuery := strings.Cut(target, "?")

	if pathPart != routePath {
		return "", false
	}

	location = routePath + "/"

	if hasQuery {
		location += "?" + query
	}

	return location, true
}
