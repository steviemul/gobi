package handler

import "testing"

func TestDirectoryRedirect(t *testing.T) {
	tests := []struct {
		name, route, target string
		want                string
		ok                  bool
	}{
		{"bare route", "/demo", "/demo", "/demo/", true},
		{"keeps query", "/demo", "/demo?x=1&y=2", "/demo/?x=1&y=2", true},
		{"already slashed", "/demo", "/demo/", "", false},
		{"file inside route", "/demo", "/demo/index.html", "", false},
		{"similar prefix", "/demo", "/demonstration", "", false},
		{"other path", "/demo", "/other", "", false},
		{"route with trailing slash", "/demo/", "/demo/", "", false},
		{"host-only route", "", "/", "", false},
		{"root route", "/", "/", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := directoryRedirect(tt.route, tt.target)
			if got != tt.want || ok != tt.ok {
				t.Errorf("directoryRedirect(%q, %q) = %q, %v; want %q, %v",
					tt.route, tt.target, got, ok, tt.want, tt.ok)
			}
		})
	}
}
