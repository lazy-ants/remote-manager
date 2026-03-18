package cmd

import "testing"

func TestExtractUfwStatus(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Status: active\nsome other output", "active"},
		{"Status: inactive\n", "inactive"},
		{"no status here", ""},
		{"Status: active", "active"},
	}

	for _, tt := range tests {
		got := extractUfwStatus(tt.input)
		if got != tt.want {
			t.Errorf("extractUfwStatus(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
