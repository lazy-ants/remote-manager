package ssh

import "testing"

func TestStripSudoPrompt(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"[sudo] password for user:\nroot",
			"root",
		},
		{
			"some output\nmore output",
			"some output\nmore output",
		},
		{
			"[sudo] password for deploy:\n[sudo] password for deploy:\nresult",
			"result",
		},
		{
			"",
			"",
		},
	}

	for _, tt := range tests {
		got := stripSudoPrompt(tt.input)
		if got != tt.want {
			t.Errorf("stripSudoPrompt(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
