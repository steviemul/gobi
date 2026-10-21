package handler

import (
	"bytes"
	"html/template"
	"log"
	"net/url"
	"path"

	"github.com/steviemul/gobi/internal/server/http"
	"github.com/steviemul/gobi/tools/www"
)

const (
	folderIcon = "\U0001F4C1"
	fileIcon   = "\U0001F4C4"
)

// indexTemplate renders a directory listing. html/template escapes every value
// for the context it appears in (text, attribute or URL), so names containing
// markup can't break out of the page.
var indexTemplate = template.Must(template.ParseFS(www.Templates, "templates/index.html"))

type indexRow struct {
	Name      string
	Href      string
	Directory bool
}

type indexPage struct {
	Title      string
	ParentHref string
	Rows       []indexRow
	FolderIcon string
	FileIcon   string
}

// renderIndexHTML renders entries as an HTML page. routePath is the route's
// mount point (e.g. "/code") and dirPath the directory shown, relative to the
// route root, so each link points back at the index handler.
func renderIndexHTML(routePath string, dirPath string, entries []IndexEntry) http.GobiResponse {
	page := indexPage{
		Title:      indexURLPath(routePath, dirPath, true),
		Rows:       make([]indexRow, 0, len(entries)),
		FolderIcon: folderIcon,
		FileIcon:   fileIcon,
	}

	if dirPath != "" {
		page.ParentHref = indexURLPath(routePath, path.Dir(dirPath), true)
	}

	for _, entry := range entries {
		page.Rows = append(page.Rows, indexRow{
			Name:      entry.Name,
			Href:      indexURLPath(routePath, path.Join(dirPath, entry.Name), entry.Directory),
			Directory: entry.Directory,
		})
	}

	var body bytes.Buffer

	if err := indexTemplate.Execute(&body, page); err != nil {
		log.Printf("index renderer: %v", err)

		return internalErrorResponse()
	}

	response, _ := http.BuildContentResponse("1.1", "200 OK", "text/html; charset=utf-8", body.Bytes())

	return response
}

// indexURLPath builds an escaped absolute URL path under the route. Directories
// get a trailing slash.
func indexURLPath(routePath string, rel string, directory bool) string {
	p := path.Join("/", routePath, rel)

	if directory && p != "/" {
		p += "/"
	}

	return (&url.URL{Path: p}).String()
}
