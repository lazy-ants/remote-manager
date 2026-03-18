package cmd

import (
	"testing"
	"time"
)

func TestFormatUptime(t *testing.T) {
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		input string
		want  string
	}{
		{"2026-03-08 10:00:00", "10 days"},
		{"2026-03-18 00:00:00", "0 days"},
		{"2026-01-01 00:00:00", "76 days"},
		{"not a date", "not a date"},
	}

	for _, tt := range tests {
		got := formatUptime(tt.input, now)
		if got != tt.want {
			t.Errorf("formatUptime(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
