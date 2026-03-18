# PHP Reference Guide

> **Shared reference for all agents**: Quick reference for understanding the PHP source during migration.

## File Map: PHP → Go

| PHP File | Lines | Go Target | Key Behavior |
|---|---|---|---|
| `AbstractCommand.php` | 258 | `internal/runner/`, `cmd/root.go` | Parallel execution, tag filtering, output rendering |
| `SimpleTask.php` | 49 | `internal/ssh/client.go` | SSH via process spawn → native SSH |
| `AbstractTask.php` | 35 | (eliminated) | Base task class, not needed in Go |
| `ServerInstancesConfig.php` | 78 | `internal/config/config.go` | Config loading, tag/name filtering |
| `ServerInstanceItem.php` | 15 | `internal/config/config.go` | Server instance struct |
| `KernelCommand.php` | 36 | `cmd/kernel.go` | `uname -r`, table output |
| `OSCommand.php` | 36 | `cmd/os.go` | `cat /etc/issue`, table output |
| `UptimeCommand.php` | 46 | `cmd/uptime.go` | `uptime -s`, date parsing, table output |
| `DockerComposeVersionCommand.php` | 36 | `cmd/docker_compose_version.go` | `docker-compose -v`, table output |
| `DockerPsCommand.php` | 32 | `cmd/docker_ps.go` | `docker ps`, list output |
| `DockerPruneCommand.php` | 32 | `cmd/docker_prune.go` | `echo "y" \| docker system prune`, list output |
| `LsCommand.php` | 39 | `cmd/ls.go` | `ls` + arg, list output |
| `UpgradeCommand.php` | 37 | `cmd/upgrade.go` | apt-get update+upgrade+autoremove, sudo, list output |
| `UfwCommand.php` | 54 | `cmd/ufw.go` | `sudo ufw status/[arg]`, conditional table/list |
| `CheckRebootCommand.php` | 88 | `cmd/check_reboot.go` | Two-pass with --reboot flag, most complex |
| `SystemInfoCommand.php` | 43 | `cmd/system_info.go` | base64-encoded script → go:embed, sudo |
| `ValidateConfigCommand.php` | 92 | `cmd/validate_config.go` | Sequential, 3 checks per server |
| `Log4jCommand.php` | 34 | `cmd/log4j.go` | `sudo ls -lha`, sudo, list output |

## Key PHP Patterns to Understand

### The Marker Hack (SimpleTask.php)

PHP prepends `echo "startoutputsysteminformation" &&` to every command, then splits output on that marker. This strips SSH login banners (MOTD). **Eliminated in Go** — native SSH sessions don't include login banners.

### The Process Pool (AbstractCommand.php)

```php
$pool = Pool::create()->concurrency(20)->timeout(120);
foreach ($instances as $instance) {
    $pool->add(new SimpleTask($instance, $command, $sudo));
}
$pool->wait();
```

**Go equivalent:** goroutines + WaitGroup + buffered channel for concurrency limit.

### Sudo via SendEnv (SimpleTask.php)

PHP sets `PASSWORD` env var and uses SSH `SendEnv` to pass it to the remote server. Requires `AcceptEnv PASSWORD` in server's sshd_config. **Go improvement:** pipe password to `sudo -S` stdin directly — no server config needed.

### Output Modes

- **Table**: Two-column (Name, Value). Used by: kernel, os, uptime, docker-compose-version, ufw (without args), check-reboot, validate-config
- **List**: Multi-line per server, numbered `[1 / N] servername:`. Used by: docker-ps, docker-prune, ls, upgrade, ufw (with args), system-info, log4j
