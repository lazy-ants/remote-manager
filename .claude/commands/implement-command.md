# Implement Go Command

Migrate a specific PHP command to Go. Provide the command name as argument (e.g., `kernel`, `uptime`, `check-reboot`).

## Instructions

### 1. Read the PHP Source

Read `src/Command/<CommandName>Command.php` to understand:
- What remote command it runs
- Whether it needs sudo
- Output format (table or list)
- Any post-processing of output
- Any command-line arguments or options

### 2. Read AbstractCommand.php Context

If this is the first command being implemented, also read:
- `src/Command/AbstractCommand.php` — understand the execution pattern
- `src/Task/SimpleTask.php` — understand SSH execution

### 3. Implement the Go Command

Create `cmd/<command_name>.go` following the pattern in the go-developer agent.

Key considerations:
- Use `RunE` (not `Run`) for error returns
- Register with `rootCmd.AddCommand()` in `init()`
- Use the shared `setup(cmd)` helper for config + SSH client initialization
- Match the exact remote command from PHP
- Match the output format (table vs list)

### 4. Write Tests

If the command has post-processing logic (date parsing, regex, conditional output), write tests for that logic in the appropriate `_test.go` file.

### 5. Verify

```bash
go build .
go test ./...
```

### 6. Commit

Stage only the new/modified Go files. Commit message: "Add <command-name> command"
