# Testing Standards

> **Shared testing philosophy for all agents**: This document defines testing expectations for the Go remote-manager.

## Core Principles

1. **Tests Must Pass:** All tests must pass before committing
2. **Test New Code:** Create tests for new packages and functions
3. **Table-Driven Tests:** Use Go's table-driven test pattern for multiple inputs
4. **Interfaces for Mocking:** Test through interfaces, not concrete implementations

## Test Organization

```
internal/
  config/
    config.go
    config_test.go          # Config loading, tag filtering, duplicate detection
  ssh/
    client.go
    client_test.go          # Connection string parsing, key loading
    mock.go                 # MockExecutor for runner tests
  runner/
    runner.go
    runner_test.go          # Concurrent execution, error/timeout collection
  output/
    table_test.go           # Table rendering
    list_test.go            # List rendering
```

## What to Test

### Config Package (unit tests)

- JSON parsing with valid/invalid config
- Tag filtering (single tag, multiple tags, no match, empty tags)
- Name filtering
- Duplicate name detection
- Connection string parsing (user@host, user@host:port, edge cases)
- Sudo password fallback (instance → env var → empty)
- .env loading precedence (.env then .env.local override)

### SSH Package (unit + integration tests)

- Connection string parsing: `user@host` → (user, host, 22)
- Connection string parsing: `user@host:2222` → (user, host, 2222)
- Key file loading from PPK_NAMES env var
- **Integration (optional):** Connect to real SSH server, run `whoami`

### Runner Package (unit tests with mock)

- Parallel execution collects all results
- Errors are captured with server name context
- Timeouts are reported separately
- Progress callback fires after each server completes
- Empty instance list returns empty result

### Command-Specific Logic (unit tests)

- Uptime: date string parsing and duration formatting
- Check-reboot: output interpretation (file exists → "yes")
- UFW: status extraction via regex
- Validate-config: sequential multi-check logic

## Test Patterns

### Table-Driven Tests

```go
func TestParseConnectionString(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        wantUser string
        wantHost string
        wantPort int
    }{
        {"simple", "root@example.com", "root", "example.com", 22},
        {"with port", "deploy@host:2222", "deploy", "host", 2222},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, host, port := parseConnectionString(tt.input)
            // assertions...
        })
    }
}
```

### Mock Executor

```go
type MockExecutor struct {
    Results map[string]string  // server name → output
    Errors  map[string]error   // server name → error
}

func (m *MockExecutor) Run(ctx context.Context, inst config.ServerInstance, cmd string) (string, error) {
    if err, ok := m.Errors[inst.Name]; ok {
        return "", err
    }
    return m.Results[inst.Name], nil
}
```

## Running Tests

```bash
go test ./...                          # All tests
go test ./internal/config/ -v          # Specific package, verbose
go test ./... -race                    # Race detector
go test ./... -cover                   # Coverage report
go test ./... -coverprofile=cover.out  # Coverage file
go tool cover -html=cover.out          # View coverage in browser
```

## What NOT to Test

- Cobra command registration (framework concern)
- `tablewriter` rendering details (library concern)
- SSH protocol internals (golang.org/x/crypto concern)
- Exact output formatting (brittle, changes often)
