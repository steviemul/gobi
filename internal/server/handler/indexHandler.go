package handler

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/steviemul/gobi/internal/server/http"
)

type IndexRequestHandler struct {
	root string
	host string
	path string

	// html renders directory listings as a web page instead of JSON. Files are
	// always returned as their contents.
	html bool
}

// directoryRenderer turns a directory's entries into a response.
type directoryRenderer func(entries []IndexEntry) http.GobiResponse

func (eh IndexRequestHandler) HandleRequest(r *http.GobiRequest) http.GobiResponse {
	switch r.Method() {
	case "GET":
		return eh.Get(r)
	default:
		return methodNotAllowed("GET")
	}
}

func (eh IndexRequestHandler) CanProcess(r *http.GobiRequest) bool {
	if eh.host != "" {
		if !matchesHost(eh.host, r) {
			return false
		}
	}

	_, ok := pathResource(eh.path, r.Target())

	log.Printf("index handler matched on [%s]", eh.path)

	return ok
}

func (eh IndexRequestHandler) Get(r *http.GobiRequest) http.GobiResponse {
	path, _ := pathResourceWithDefault(eh.path, r.Target(), "")

	var render directoryRenderer // nil renders JSON

	if eh.html {
		render = func(entries []IndexEntry) http.GobiResponse {
			return renderIndexHTML(eh.path, strings.TrimSuffix(path, "/"), entries)
		}
	}

	return getIndexContent(getRoot(eh.root), path, render)
}

func getRoot(root string) string {
	if root != "" {
		return root
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get working directory: %v", err)
	}

	return dir
}

type IndexEntry struct {
	Name      string `json:"name"`
	Directory bool   `json:"directory"`
}

// getIndexContent serves the target at path, relative to root. A blank path
// targets root itself. A directory is answered by render, or with a JSON array
// of its entries when render is nil; a file is answered with its contents and a
// content type taken from its extension.
func getIndexContent(root string, path string, render directoryRenderer) http.GobiResponse {
	// fs.ValidPath rejects "..", so the target can't escape root. It also
	// rejects a trailing slash, which is fine for a directory, so trim it first.
	path = strings.TrimSuffix(path, "/")

	if path != "" && (!fs.ValidPath(path) || strings.ContainsRune(path, 0)) {
		return notFoundResponse()
	}

	target := filepath.Join(root, path)

	info, err := os.Stat(target)
	if err != nil {
		return statErrorResponse(target, err)
	}

	if info.IsDir() {
		return directoryResponse(target, render)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		return statErrorResponse(target, err)
	}

	contentType := mime.TypeByExtension(filepath.Ext(target))

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	response, _ := http.BuildContentResponse("1.1", "200 OK", contentType, data)

	return response
}

func directoryResponse(dir string, render directoryRenderer) http.GobiResponse {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return statErrorResponse(dir, err)
	}

	// Non-nil so an empty directory marshals to [] rather than null.
	entries := make([]IndexEntry, 0, len(dirEntries))

	for _, entry := range dirEntries {
		entries = append(entries, IndexEntry{
			Name:      entry.Name(),
			Directory: entry.IsDir(),
		})
	}

	if render != nil {
		return render(entries)
	}

	body, err := json.Marshal(entries)
	if err != nil {
		log.Printf("index handler: marshalling entries for %s: %v", dir, err)

		return internalErrorResponse()
	}

	response, _ := http.BuildContentResponse("1.1", "200 OK", "application/json", body)

	return response
}

func statErrorResponse(target string, err error) http.GobiResponse {
	if errors.Is(err, fs.ErrNotExist) {
		return notFoundResponse()
	}

	log.Printf("index handler: reading %s: %v", target, err)

	return internalErrorResponse()
}

func notFoundResponse() http.GobiResponse {
	response, _ := http.BuildResponse("1.1", "404 Not Found", "Not Found\n")

	return response
}

func internalErrorResponse() http.GobiResponse {
	response, _ := http.BuildResponse("1.1", "500 Internal Server Error", "Internal server error\n")

	return response
}
