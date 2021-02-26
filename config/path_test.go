package config

import (
	"testing"
)

func TestSafeFilename(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"node-*.md", "node-_.md"},
		{"path/to/node.md", "path_to_node.md"},
		{"path\\to\\node.md", "path_to_node.md"},
		{"node-/\\.md", "node-__.md"},
		{"user node.md", "user_node.md"},
		{"internet(dst\\: https\\:*.example.com)-63fad0a.md", "internet(dst___https___.example.com)-63fad0a.md"},
	}
	for _, tt := range tests {
		got := safeFilename(tt.in)
		if got != tt.want {
			t.Errorf("got %v want %v", got, tt.want)
		}
		got2 := safeFilename(got)
		if got2 != tt.want {
			t.Errorf("got %v want %v", got2, tt.want)
		}
	}
}
