package handler

import (
	"mime"
	"testing"
)

func TestTextExtensionsRegistered(t *testing.T) {
	for _, ext := range textExtensions {
		if got := mime.TypeByExtension(ext); got != textPlain {
			t.Errorf("TypeByExtension(%q) = %q, want %q", ext, got, textPlain)
		}
	}
}

func TestExistingTypesUntouched(t *testing.T) {
	tests := map[string]string{
		".html": "text/html; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".json": "application/json",
		".svg":  "image/svg+xml",
	}

	for ext, want := range tests {
		if got := mime.TypeByExtension(ext); got != want {
			t.Errorf("TypeByExtension(%q) = %q, want %q", ext, got, want)
		}
	}
}
