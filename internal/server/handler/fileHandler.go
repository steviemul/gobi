package handler

import (
	"io/fs"
	"mime"
	"os"
	"path"
	"strings"

	"github.com/steviemul/gobi/internal/server/http"
)

type FileRequestHandler struct {
	root string
	host string
	path string
}

func (eh FileRequestHandler) HandleRequest(r *http.GobiRequest) http.GobiResponse {
	resource, _ := pathResource(eh.path, r.Target())

	// fs.ValidPath rejects "..", empty and rooted names, so the lookup can't escape gobiRoot.
	if fs.ValidPath(resource) {
		data, err := os.ReadFile(eh.root + "/" + resource)

		if err == nil {
			contentType := mime.TypeByExtension(path.Ext(resource))

			if contentType == "" {
				contentType = "application/octet-stream"
			}

			response, _ := http.BuildContentResponse("1.1", "200 OK", contentType, data)

			return response
		}
	}

	response, _ := http.BuildResponse("1.1", "404 Not Found", "Not Found\n")

	return response
}

func (eh FileRequestHandler) CanProcess(r *http.GobiRequest) bool {
	_, ok := pathResource(eh.path, r.Target())

	return ok
}

func pathResource(path string, target string) (resource string, ok bool) {
	target, _, _ = strings.Cut(target, "?")
	target, _, _ = strings.Cut(target, "#")

	rest, found := strings.CutPrefix(target, path)

	if !found || (rest != "" && rest[0] != '/') {
		return "", false
	}

	resource = strings.TrimPrefix(rest, "/")

	if resource == "" {
		resource = indexFile
	}

	return resource, true
}
