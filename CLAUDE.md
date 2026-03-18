# Claude Code Rules for Remote Manager

## Project Status: Go Migration (from PHP)

This project is being rewritten from PHP (Symfony Console) to Go (Cobra + native SSH). The migration plan is in `docs/go-migration-plan.md`. The PHP source in `src/` is the reference implementation — read it to understand behavior, but write Go code.

## Go Conventions

- **Module**: `github.com/lazy-ants/remote-manager`
- **Go version**: 1.22+
- **Entry point**: `main.go`
- **Internal packages**: `internal/config`, `internal/ssh`, `internal/runner`, `internal/output`
- **Commands**: `cmd/` directory, one file per cobra command
- **Embedded scripts**: `scripts/` directory, loaded via `go:embed`

## Build & Run

```bash
go build -o remote-manager .           # Build binary
go run . <command> [flags]             # Run directly
go test ./...                          # Run all tests
golangci-lint run                      # Lint
```

No Docker required for development or running. The binary is self-contained.

## Git Workflow

- Present tense, imperative mood (e.g., "Add kernel command")
- **NEVER add Claude signature** — no "Co-Authored-By: Claude" or AI attribution
- Stage specific files only, never `git add -A`
- After changes: `git add <files> && git commit -m "message"`

## Architecture Rules

- **Native SSH only** — use `golang.org/x/crypto/ssh`, never shell out to `ssh` binary
- **Interfaces for testability** — SSH client must implement `Executor` interface so runner can be tested with mocks
- **No global state** — pass config and clients explicitly
- **Errors are values** — return errors, don't panic. Distinguish connection errors, auth errors, and command errors
- **Context for cancellation** — all SSH operations accept `context.Context` for timeouts

## PHP Reference (read-only)

The original PHP source is in `src/`. Key files:
- `src/Command/AbstractCommand.php` — core parallel execution pattern
- `src/Task/SimpleTask.php` — SSH execution via process spawning
- `src/Configuration/ServerInstancesConfig.php` — config loading and filtering

**Do NOT modify PHP files.** They are reference only during migration.

## Config Format

`config.json` format is preserved from PHP version for backward compatibility:

```json
{
  "instances": [
    {
      "name": "example.com",
      "connection-string": "user@host:port",
      "sudo-password": "optional",
      "tags": "staging,client1"
    }
  ]
}
```

## Testing Strategy

- **Unit tests**: config parsing, connection string parsing, output formatting
- **Integration tests**: SSH client against test server (when available)
- **Mock SSH executor**: for runner tests without real servers
- Run `go test ./...` before committing

## Dead Code & Backward Compatibility

- Actively remove dead code
- No backward compatibility shims needed — this is a rewrite
- When a Go command fully replaces a PHP command, note it in the commit message

## Knowledge Capture

- Migration decisions and learnings go in `docs/`
- Agent files define workflow only, not domain knowledge
