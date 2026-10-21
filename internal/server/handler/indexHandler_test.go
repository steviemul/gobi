package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetIndexContent(t *testing.T) {
	root := t.TempDir()

	must := func(err error) {
		t.Helper()

		if err != nil {
			t.Fatal(err)
		}
	}

	must(os.Mkdir(filepath.Join(root, "docs"), 0o755))
	must(os.Mkdir(filepath.Join(root, "empty"), 0o755))
	must(os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello"), 0o644))
	must(os.WriteFile(filepath.Join(root, "docs", "page.html"), []byte("<p>hi</p>"), 0o644))
	must(os.WriteFile(filepath.Join(root, "blob.unknownext"), []byte{0, 1, 2}, 0o644))

	tests := []struct {
		name, path   string
		wantStatus   string
		wantType     string
		wantBodyPart string
	}{
		{
			"blank path lists root", "", "200 OK", "application/json",
			`[{"name":"a.txt","directory":false},{"name":"blob.unknownext","directory":false},{"name":"docs","directory":true},{"name":"empty","directory":true}]`,
		},
		{"subdirectory", "docs", "200 OK", "application/json", `[{"name":"page.html","directory":false}]`},
		{"subdirectory with trailing slash", "docs/", "200 OK", "application/json", `[{"name":"page.html","directory":false}]`},
		{"empty directory is [] not null", "empty", "200 OK", "application/json", "[]"},
		{"text file", "a.txt", "200 OK", "text/plain; charset=utf-8", "hello"},
		{"nested html file", "docs/page.html", "200 OK", "text/html; charset=utf-8", "<p>hi</p>"},
		{"unknown extension", "blob.unknownext", "200 OK", "application/octet-stream", "\x00\x01\x02"},
		{"encoded traversal", "%2e%2e/x", "404 Not Found", "", ""},
		{"nul byte", "a.txt\x00", "404 Not Found", "", ""},
		{"missing", "nope.txt", "404 Not Found", "", ""},
		{"parent traversal", "../x", "404 Not Found", "", ""},
		{"traversal via subdirectory", "docs/../../x", "404 Not Found", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := getIndexContent(root, tt.path, nil)
			raw := resp.String()

			if !strings.HasPrefix(raw, "HTTP/1.1 "+tt.wantStatus+"\r\n") {
				t.Fatalf("status line = %q, want %q", strings.SplitN(raw, "\r\n", 2)[0], tt.wantStatus)
			}

			if tt.wantType != "" && !strings.Contains(raw, "Content-Type: "+tt.wantType+"\r\n") {
				t.Errorf("response missing Content-Type %q:\n%s", tt.wantType, raw)
			}

			if tt.wantBodyPart != "" && !strings.HasSuffix(raw, "\r\n\r\n"+tt.wantBodyPart) {
				t.Errorf("body mismatch, want suffix %q in:\n%s", tt.wantBodyPart, raw)
			}
		})
	}
}
