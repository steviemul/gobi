package handler

import (
	"io/fs"
	"mime"
	"os"
	"path"

	"github.com/steviemul/gobi/internal/server/http"
)

type FileRequestHandler struct {
	root string
	host string
	path string
}

func (eh FileRequestHandler) HandleRequest(r *http.GobiRequest) http.GobiResponse {
	switch r.Method() {
	case "GET":
		return eh.Get(r)
	case "HEAD":
		return eh.Head(r)
	default:
		return methodNotAllowed("GET, HEAD")
	}
}

func (eh FileRequestHandler) CanProcess(r *http.GobiRequest) bool {
	if eh.host != "" {
		if !matchesHost(eh.host, r) {
			return false
		}
	}

	_, ok := pathResource(eh.path, r.Target())

	return ok
}

func (eh FileRequestHandler) getResource(r *http.GobiRequest, includeBody bool) http.GobiResponse {
	// "/demo" must become "/demo/" so the page's relative URLs resolve inside the route.
	if location, ok := directoryRedirect(eh.path, r.Target()); ok {
		response, _ := http.BuildRedirectResponse("1.1", "301 Moved Permanently", location)

		return response
	}

	resource, _ := pathResource(eh.path, r.Target())

	// fs.ValidPath rejects "..", empty and rooted names, so the lookup can't escape gobiRoot.
	if fs.ValidPath(resource) {
		data, err := os.ReadFile(eh.root + "/" + resource)

		if err == nil {
			contentType := mime.TypeByExtension(path.Ext(resource))

			if contentType == "" {
				contentType = "application/octet-stream"
			}

			if !includeBody {
				data = nil
			}

			response, _ := http.BuildContentResponse("1.1", "200 OK", contentType, data)

			return response
		}
	}

	response, _ := http.BuildResponse("1.1", "404 Not Found", "Not Found\n")

	return response
}

func (eh FileRequestHandler) Get(r *http.GobiRequest) http.GobiResponse {
	return eh.getResource(r, true)
}

func (eh FileRequestHandler) Head(r *http.GobiRequest) http.GobiResponse {
	return eh.getResource(r, false)
}

func methodNotAllowed(allowedMethods string) http.GobiResponse {
	response, _ := http.BuildResponse("1.1", "405 method not allowed", "Method not allowed\n")

	response.SetHeader("Allowed", allowedMethods)

	return response
}
