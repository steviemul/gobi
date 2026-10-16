package server

import "testing"

func TestHTTPSLocation(t *testing.T) {
	tests := []struct {
		name, host, target, port string
		want                     string
		ok                       bool
	}{
		{"host with port", "localhost:8080", "/a/b?x=1", "8443", "https://localhost:8443/a/b?x=1", true},
		{"host without port", "gobi.local", "/", "8443", "https://gobi.local:8443/", true},
		{"default https port", "gobi.local:8080", "/x", "443", "https://gobi.local/x", true},
		{"ipv6 with port", "[::1]:8080", "/", "8443", "https://[::1]:8443/", true},
		{"ipv6 without port", "[::1]", "/", "8443", "https://[::1]:8443/", true},
		{"ipv6 default port", "[::1]:8080", "/", "443", "https://[::1]/", true},
		{"missing host", "", "/", "8443", "", false},
		{"host with slash", "evil.com/x", "/", "8443", "", false},
		{"host with userinfo", "a@evil.com", "/", "8443", "", false},
		{"asterisk target", "localhost", "*", "8443", "", false},
		{"absolute-form target", "localhost", "http://other/x", "8443", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := httpsLocation(tt.host, tt.target, tt.port)
			if got != tt.want || ok != tt.ok {
				t.Errorf("httpsLocation(%q, %q, %q) = %q, %v; want %q, %v",
					tt.host, tt.target, tt.port, got, ok, tt.want, tt.ok)
			}
		})
	}
}
