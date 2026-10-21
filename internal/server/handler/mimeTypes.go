package handler

import (
	"log"
	"mime"
)

const textPlain = "text/plain; charset=utf-8"

// textExtensions are text formats that mime.TypeByExtension doesn't know on every
// OS. Without a type they're sent as application/octet-stream, which browsers
// download. text/plain makes them display, which is what a file browser wants.
// Deliberately not text/html: serving unknown text as HTML would run any script in it.
var textExtensions = []string{
	".yaml", ".yml", ".toml", ".ini", ".conf", ".log",
	".md", ".go", ".mod", ".sum",
	".sh", ".py", ".rs", ".java", ".c", ".h", ".ts", ".tsx", ".jsx", ".sql",
}

// Registered once at startup; the registry is process-wide, so the file and
// index handlers both pick these up through mime.TypeByExtension.
func init() {
	for _, ext := range textExtensions {
		if err := mime.AddExtensionType(ext, textPlain); err != nil {
			log.Printf("mime: registering %s: %v", ext, err)
		}
	}
}
