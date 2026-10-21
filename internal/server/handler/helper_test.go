package handler

import "testing"

func TestPathResourceDecodes(t *testing.T) {
	tests := []struct {
		target, want string
		ok           bool
	}{
		{"/code/a%20b.txt", "a b.txt", true},
		{"/code/q%3F%23.txt?x=1", "q?#.txt", true},
		{"/code/sub/n.txt", "sub/n.txt", true},
		{"/code/%zz", "", false},
	}

	for _, tt := range tests {
		got, ok := pathResourceWithDefault("/code", tt.target, "")
		if got != tt.want || ok != tt.ok {
			t.Errorf("pathResource(%q) = %q, %v; want %q, %v", tt.target, got, ok, tt.want, tt.ok)
		}
	}
}
