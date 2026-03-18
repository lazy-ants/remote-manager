# Audit Migration Progress

Compare the Go implementation against the PHP reference to assess migration completeness.

## Instructions

### 1. Inventory All PHP Commands

Read every file in `src/Command/` and list all commands with their behavior.

### 2. Check Go Implementation

For each PHP command, verify:

- [ ] Go command file exists in `cmd/`
- [ ] Remote command matches PHP version
- [ ] Output format matches (table vs list)
- [ ] Sudo handling is correct
- [ ] Post-processing logic is preserved (date parsing, regex extraction, etc.)
- [ ] Edge cases are handled (empty output, errors, timeouts)

### 3. Check Infrastructure

- [ ] `internal/config/` — config.json loading matches PHP behavior
- [ ] `internal/config/` — .env loading with .env.local override
- [ ] `internal/ssh/` — native SSH client works (key auth, password sudo)
- [ ] `internal/runner/` — parallel execution with progress bar
- [ ] `internal/output/` — table and list rendering
- [ ] Tag filtering works identically to PHP

### 4. Check Build/Distribution

- [ ] `go build` produces working binary
- [ ] Dockerfile builds and runs
- [ ] GoReleaser config exists
- [ ] Makefile has build, test, lint, docker targets

### 5. Check Tests

- [ ] Config parsing tests exist
- [ ] Connection string parsing tests exist
- [ ] Runner tests with mock executor exist
- [ ] Command-specific logic tests exist (uptime parsing, reboot check, etc.)

### 6. Report Format

```markdown
## Migration Audit

### Commands
| Command | PHP | Go | Tests | Notes |
|---------|-----|-----|-------|-------|

### Infrastructure
| Component | Status | Notes |
|-----------|--------|-------|

### Build
| Target | Status | Notes |
|--------|--------|-------|

### Remaining Work
1. [Priority-ordered list]
```
