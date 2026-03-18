package cmd

import "testing"

func TestInterpretRebootCheck(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"/var/run/reboot-required", "yes"},
		{"/var/run/reboot-required\n", "yes"},
		{"", "no"},
		{"ls: cannot access", "no"},
		{"some other output", "no"},
	}

	for _, tt := range tests {
		got := interpretRebootCheck(tt.input)
		if got != tt.want {
			t.Errorf("interpretRebootCheck(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
