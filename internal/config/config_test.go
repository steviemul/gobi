package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexConfig(t *testing.T) {
	cfg := IndexConfig()

	if cfg.Port() != DefaultHTTPPort {
		t.Errorf("port = %q, want %q", cfg.Port(), DefaultHTTPPort)
	}

	if cfg.Listen() != LoopbackAddress {
		t.Errorf("listen = %q, want %q", cfg.Listen(), LoopbackAddress)
	}

	if cfg.TLS().Enabled {
		t.Error("TLS should be disabled")
	}

	if cfg.Status() {
		t.Error("status page should be disabled")
	}

	routes := cfg.Routes()

	if len(routes) != 1 {
		t.Fatalf("got %d routes, want 1", len(routes))
	}

	r := routes[0]

	if r.Index == nil || !*r.Index {
		t.Error("route should be an index route")
	}

	if r.Root != "" || r.Path != "" || r.Host != "" || r.Proxy != "" {
		t.Errorf("route should be otherwise empty, got %+v", r)
	}
}

func TestListenDefaultsToAllInterfaces(t *testing.T) {
	cfg := DefaultConfig()

	if got := cfg.Listen(); got != "" {
		t.Errorf("default listen = %q, want empty (all interfaces)", got)
	}
}

func TestListenFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.yaml")

	if err := os.WriteFile(path, []byte("listen: 127.0.0.1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := BuildConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Listen() != "127.0.0.1" {
		t.Errorf("listen = %q, want 127.0.0.1", cfg.Listen())
	}
}
