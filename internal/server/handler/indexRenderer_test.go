package handler

import (
	"strings"
	"testing"
)

func TestRenderIndexHTML(t *testing.T) {
	entries := []IndexEntry{
		{Name: "docs", Directory: true},
		{Name: "a b.txt"},
		{Name: "q?#.txt"},
		{Name: "<script>alert(1)</script>.html"},
	}

	resp := renderIndexHTML("/code", "sub", entries)
	raw := resp.String()

	for _, want := range []string{
		"HTTP/1.1 200 OK\r\n",
		"Content-Type: text/html; charset=utf-8\r\n",
		`<title>Index of /code/sub/</title>`,
		`<a href="/code/">..</a>`, // parent of "sub" is the route root
		`<span class="icon">` + folderIcon + `</span><a href="/code/sub/docs/">docs</a>`,
		`<span class="icon">` + fileIcon + `</span><a href="/code/sub/a%20b.txt">a b.txt</a>`,
		`href="/code/sub/q%3f%23.txt"`, // ? and # must not end the path
		`&lt;script&gt;alert(1)&lt;/script&gt;.html</a>`,
	} {
		if !strings.Contains(strings.ToLower(raw), strings.ToLower(want)) {
			t.Errorf("output missing %q:\n%s", want, raw)
		}
	}

	if strings.Contains(raw, "<script>") {
		t.Errorf("unescaped markup in output:\n%s", raw)
	}
}

func TestRenderIndexHTMLRoot(t *testing.T) {
	resp := renderIndexHTML("/code", "", nil)
	raw := resp.String()

	if strings.Contains(raw, ">..<") {
		t.Errorf("route root should have no parent link:\n%s", raw)
	}

	if !strings.Contains(raw, "(empty)") {
		t.Errorf("empty directory should say so:\n%s", raw)
	}
}
