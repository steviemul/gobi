package handler

import (
	"io/fs"
	"mime"
	"path"
	"strings"

	"github.com/steviemul/gobi/internal/server/http"
	"github.com/steviemul/gobi/tools/www"
)

var STATUS = "status"

const (
	gobiPrefix = "/gobi"
	gobiRoot   = "gobi"
	indexFile  = "index.html"
)

var gobiWebContent = www.Content

type GobiRequestHandler struct{}

func (eh GobiRequestHandler) HandleRequest(r *http.GobiRequest) http.GobiResponse {
	resource, _ := gobiResource(r.Target())

	// fs.ValidPath rejects "..", empty and rooted names, so the lookup can't escape gobiRoot.
	if fs.ValidPath(resource) {
		data, err := gobiWebContent.ReadFile(gobiRoot + "/" + resource)

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

func (eh GobiRequestHandler) CanProcess(r *http.GobiRequest) bool {
	_, ok := gobiResource(r.Target())

	return ok
}

// gobiResource maps a request target under /gobi to a file name within the embedded content.
// ok is false if the target is not under /gobi. /gobi and /gobi/ map to index.html.
func gobiResource(target string) (resource string, ok bool) {
	target, _, _ = strings.Cut(target, "?")
	target, _, _ = strings.Cut(target, "#")

	rest, found := strings.CutPrefix(target, gobiPrefix)

	if !found || (rest != "" && rest[0] != '/') {
		return "", false
	}

	resource = strings.TrimPrefix(rest, "/")

	if resource == "" {
		resource = indexFile
	}

	return resource, true
}
