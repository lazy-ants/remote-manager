package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ServerInstance struct {
	Name             string
	ConnectionString string
	SudoPassword     string
	Tags             []string
}

type Config struct {
	Instances []ServerInstance
}

type jsonConfig struct {
	Instances []jsonInstance `json:"instances"`
}

type jsonInstance struct {
	Name             string `json:"name"`
	ConnectionString string `json:"connection-string"`
	SudoPassword     string `json:"sudo-password"`
	Tags             string `json:"tags"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var jc jsonConfig
	if err := json.Unmarshal(data, &jc); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	seen := make(map[string]bool)
	var duplicates []string

	cfg := &Config{}
	for _, ji := range jc.Instances {
		tags := splitTags(ji.Tags)

		sudoPassword := ji.SudoPassword
		if sudoPassword == "" {
			sudoPassword = os.Getenv("SUDO_PASSWORD")
		}

		if seen[ji.Name] {
			duplicates = append(duplicates, ji.Name)
		}
		seen[ji.Name] = true

		cfg.Instances = append(cfg.Instances, ServerInstance{
			Name:             ji.Name,
			ConnectionString: ji.ConnectionString,
			SudoPassword:     sudoPassword,
			Tags:             tags,
		})
	}

	if len(duplicates) > 0 {
		unique := uniqueStrings(duplicates)
		return nil, fmt.Errorf("duplicate server names in config: %s", strings.Join(unique, ", "))
	}

	return cfg, nil
}

func (c *Config) FilterByTags(tags []string) *Config {
	filtered := filterEmpty(tags)
	if len(filtered) == 0 {
		return c
	}

	var result []ServerInstance
	for _, inst := range c.Instances {
		if hasAllTags(inst.Tags, filtered) {
			result = append(result, inst)
		}
	}
	c.Instances = result
	return c
}

func (c *Config) FilterByNames(names []string) *Config {
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	var result []ServerInstance
	for _, inst := range c.Instances {
		if nameSet[inst.Name] {
			result = append(result, inst)
		}
	}
	c.Instances = result
	return c
}

func ParseConnectionString(cs string) (user, host string, port int, err error) {
	port = 22

	atIdx := strings.LastIndex(cs, "@")
	if atIdx < 0 {
		return "", "", 0, fmt.Errorf("invalid connection string %q: missing @", cs)
	}

	user = cs[:atIdx]
	hostPort := cs[atIdx+1:]

	if colonIdx := strings.LastIndex(hostPort, ":"); colonIdx >= 0 {
		host = hostPort[:colonIdx]
		port, err = strconv.Atoi(hostPort[colonIdx+1:])
		if err != nil {
			return "", "", 0, fmt.Errorf("invalid port in %q: %w", cs, err)
		}
	} else {
		host = hostPort
	}

	if user == "" || host == "" {
		return "", "", 0, fmt.Errorf("invalid connection string %q: empty user or host", cs)
	}

	return user, host, port, nil
}

func splitTags(s string) []string {
	if s == "" {
		return nil
	}
	var tags []string
	for _, t := range strings.Split(s, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func hasAllTags(serverTags, requiredTags []string) bool {
	tagSet := make(map[string]bool)
	for _, t := range serverTags {
		tagSet[t] = true
	}
	for _, t := range requiredTags {
		if !tagSet[t] {
			return false
		}
	}
	return true
}

func filterEmpty(ss []string) []string {
	var result []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func uniqueStrings(ss []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
