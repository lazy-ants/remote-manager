package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	data := `{
		"instances": [
			{
				"name": "server1",
				"connection-string": "user@host1:22",
				"sudo-password": "pass1",
				"tags": "web,production"
			},
			{
				"name": "server2",
				"connection-string": "user@host2",
				"tags": "staging"
			}
		]
	}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Instances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(cfg.Instances))
	}

	inst := cfg.Instances[0]
	if inst.Name != "server1" {
		t.Errorf("expected name server1, got %s", inst.Name)
	}
	if inst.SudoPassword != "pass1" {
		t.Errorf("expected password pass1, got %s", inst.SudoPassword)
	}
	if len(inst.Tags) != 2 || inst.Tags[0] != "web" || inst.Tags[1] != "production" {
		t.Errorf("unexpected tags: %v", inst.Tags)
	}
}

func TestLoadConfigDuplicateNames(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	data := `{
		"instances": [
			{"name": "dup", "connection-string": "u@h1", "tags": ""},
			{"name": "dup", "connection-string": "u@h2", "tags": ""}
		]
	}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for duplicate names")
	}
}

func TestLoadConfigSudoPasswordFallback(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	data := `{
		"instances": [
			{"name": "s1", "connection-string": "u@h1", "tags": ""}
		]
	}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("SUDO_PASSWORD", "envpass")

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Instances[0].SudoPassword != "envpass" {
		t.Errorf("expected envpass, got %s", cfg.Instances[0].SudoPassword)
	}
}

func TestFilterByTags(t *testing.T) {
	cfg := &Config{
		Instances: []ServerInstance{
			{Name: "s1", Tags: []string{"web", "production"}},
			{Name: "s2", Tags: []string{"staging"}},
			{Name: "s3", Tags: []string{"web", "staging"}},
		},
	}

	cfg.FilterByTags([]string{"web"})
	if len(cfg.Instances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(cfg.Instances))
	}
}

func TestFilterByTagsANDLogic(t *testing.T) {
	cfg := &Config{
		Instances: []ServerInstance{
			{Name: "s1", Tags: []string{"web", "production"}},
			{Name: "s2", Tags: []string{"web", "staging"}},
		},
	}

	cfg.FilterByTags([]string{"web", "production"})
	if len(cfg.Instances) != 1 || cfg.Instances[0].Name != "s1" {
		t.Fatalf("expected only s1, got %v", cfg.Instances)
	}
}

func TestFilterByTagsEmpty(t *testing.T) {
	cfg := &Config{
		Instances: []ServerInstance{
			{Name: "s1", Tags: []string{"web"}},
			{Name: "s2", Tags: []string{"staging"}},
		},
	}

	cfg.FilterByTags(nil)
	if len(cfg.Instances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(cfg.Instances))
	}
}

func TestFilterByNames(t *testing.T) {
	cfg := &Config{
		Instances: []ServerInstance{
			{Name: "s1"},
			{Name: "s2"},
			{Name: "s3"},
		},
	}

	cfg.FilterByNames([]string{"s1", "s3"})
	if len(cfg.Instances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(cfg.Instances))
	}
}

func TestParseConnectionString(t *testing.T) {
	tests := []struct {
		input    string
		user     string
		host     string
		port     int
		wantErr  bool
	}{
		{"user@host", "user", "host", 22, false},
		{"user@host:2222", "user", "host", 2222, false},
		{"deploy@example.com:22", "deploy", "example.com", 22, false},
		{"noatsign", "", "", 0, true},
		{"@host", "", "", 0, true},
		{"user@", "", "", 0, true},
		{"user@host:abc", "", "", 0, true},
	}

	for _, tt := range tests {
		user, host, port, err := ParseConnectionString(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseConnectionString(%q): expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseConnectionString(%q): %v", tt.input, err)
			continue
		}
		if user != tt.user || host != tt.host || port != tt.port {
			t.Errorf("ParseConnectionString(%q) = (%s, %s, %d), want (%s, %s, %d)",
				tt.input, user, host, port, tt.user, tt.host, tt.port)
		}
	}
}
