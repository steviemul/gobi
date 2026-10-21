package handler

import (
	"fmt"
	"log"
	"os"

	"github.com/steviemul/gobi/internal/server/http"
)

type IndexRequestHandler struct {
	root string
	host string
	path string
}

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
	root := getRoot(eh.root)

	body := fmt.Sprintf("index path and root is [%s, %s]", path, root)

	response, _ := http.BuildContentResponse("1.1", "200 OK", "text/plain; charset=utf-8", []byte(body))

	return response
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
