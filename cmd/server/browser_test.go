package main

import (
	"slices"
	"testing"

	"github.com/steviemul/gobi/internal/config"
)

func TestBrowserURL(t *testing.T) {
	if got, want := browserURL(config.IndexConfig()), "http://localhost:8080/"; got != want {
		t.Errorf("IndexConfig URL = %q, want %q", got, want)
	}

	if got, want := browserURL(config.DefaultConfig()), "http://localhost:8080/"; got != want {
		t.Errorf("default URL = %q, want %q", got, want)
	}
}

func TestBrowserCommand(t *testing.T) {
	tests := []struct {
		goos, wantName string
		wantArgs       []string
	}{
		{"darwin", "open", []string{"http://x/"}},
		{"linux", "xdg-open", []string{"http://x/"}},
		{"windows", "rundll32", []string{"url.dll,FileProtocolHandler", "http://x/"}},
	}

	for _, tt := range tests {
		name, args := browserCommand(tt.goos, "http://x/")

		if name != tt.wantName || !slices.Equal(args, tt.wantArgs) {
			t.Errorf("browserCommand(%q) = %q %v, want %q %v", tt.goos, name, args, tt.wantName, tt.wantArgs)
		}
	}
}
