package server

import (
	"net"
	"strings"

	"github.com/steviemul/gobi/internal/server/http"
)

const defaultHTTPSPort = "443"

// httpsLocation builds the https URL equivalent to a plain-HTTP request. The
// host comes from the Host header with any port replaced by tlsPort; the
// target (path and query) is kept as is. It reports false when the request
// can't be redirected: no usable Host, or a target that isn't origin-form.
func httpsLocation(hostHeader, target, tlsPort string) (string, bool) {
	if hostHeader == "" || !strings.HasPrefix(target, "/") {
		return "", false
	}

	if strings.ContainsAny(hostHeader, " \t/\\@?#") {
		return "", false
	}

	host, _, err := net.SplitHostPort(hostHeader)
	if err != nil {
		host = strings.Trim(hostHeader, "[]") // no port, possibly a bare IPv6 literal
	}

	authority := net.JoinHostPort(host, tlsPort)

	if tlsPort == defaultHTTPSPort {
		authority = host

		if strings.Contains(host, ":") {
			authority = "[" + host + "]"
		}
	}

	return "https://" + authority + target, true
}

func redirectToHTTPS(req *http.GobiRequest, tlsPort string) http.GobiResponse {
	location, ok := httpsLocation(req.Headers()["host"], req.Target(), tlsPort)
	if !ok {
		resp, _ := http.BuildResponse("1.1", "400 Bad Request", "Cannot redirect to HTTPS: missing or invalid Host\n")
		return resp
	}

	resp, _ := http.BuildRedirectResponse("1.1", "302 Found", location)

	return resp
}
