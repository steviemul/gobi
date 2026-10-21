package config

import "testing"

func TestIndexConfig(t *testing.T) {
	cfg := IndexConfig()

	if cfg.Port() != DefaultHTTPPort {
		t.Errorf("port = %q, want %q", cfg.Port(), DefaultHTTPPort)
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
