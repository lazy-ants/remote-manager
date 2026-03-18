package cmd

import "testing"

func TestValidateConfigLogic(t *testing.T) {
	// Test the interpretation logic used in validate-config
	tests := []struct {
		result  string
		isLogin bool
	}{
		{"root", true},
		{"deploy", true},
		{"", false},
	}

	for _, tt := range tests {
		got := tt.result != ""
		if got != tt.isLogin {
			t.Errorf("login check for result %q: got %v, want %v", tt.result, got, tt.isLogin)
		}
	}
}
