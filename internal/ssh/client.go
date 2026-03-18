package ssh

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lazy-ants/remote-manager/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Executor defines the interface for running commands on remote servers.
type Executor interface {
	Run(ctx context.Context, instance config.ServerInstance, command string) (string, error)
	RunWithSudo(ctx context.Context, instance config.ServerInstance, command string) (string, error)
	RunWithSudoAndStdin(ctx context.Context, instance config.ServerInstance, command string, stdinContent string) (string, error)
	Close() error
}

// SSHClient implements Executor using native Go SSH.
type SSHClient struct {
	mu      sync.Mutex
	clients map[string]*ssh.Client
	signers []ssh.Signer
}

// NewSSHClient creates a new SSH client with key loading from SSH agent and PPK_NAMES.
func NewSSHClient() (*SSHClient, error) {
	c := &SSHClient{
		clients: make(map[string]*ssh.Client),
	}

	// Try SSH agent first (handles passphrase-protected keys that are already unlocked)
	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		conn, err := net.Dial("unix", sock)
		if err == nil {
			agentClient := agent.NewClient(conn)
			agentSigners, err := agentClient.Signers()
			if err == nil {
				c.signers = append(c.signers, agentSigners...)
			}
		}
	}

	// Load unprotected keys from PPK_NAMES env var (skip passphrase-protected keys
	// since the agent should already have them)
	if ppkNames := os.Getenv("PPK_NAMES"); ppkNames != "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home dir: %w", err)
		}
		for _, name := range strings.Split(ppkNames, ",") {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			keyPath := filepath.Join(home, ".ssh", name)
			signer, err := loadKey(keyPath)
			if err != nil {
				// Skip passphrase-protected keys — the agent should provide them
				continue
			}
			c.signers = append(c.signers, signer)
		}
	}

	if len(c.signers) == 0 {
		return nil, fmt.Errorf("no SSH keys available: set PPK_NAMES or start ssh-agent with loaded keys")
	}

	return c, nil
}

func loadKey(path string) (ssh.Signer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(data)
}

func (c *SSHClient) getClient(instance config.ServerInstance) (*ssh.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if client, ok := c.clients[instance.ConnectionString]; ok {
		// Test if connection is still alive
		_, _, err := client.SendRequest("keepalive@openssh.com", true, nil)
		if err == nil {
			return client, nil
		}
		// Connection dead, remove and reconnect
		client.Close()
		delete(c.clients, instance.ConnectionString)
	}

	user, host, port, err := config.ParseConnectionString(instance.ConnectionString)
	if err != nil {
		return nil, &AuthError{Host: instance.Name, Err: err}
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(c.signers...),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, &ConnError{Host: instance.Name, Addr: addr, Err: err}
	}

	c.clients[instance.ConnectionString] = client
	return client, nil
}

func (c *SSHClient) runSession(ctx context.Context, instance config.ServerInstance, command string, stdin string) (string, error) {
	client, err := c.getClient(instance)
	if err != nil {
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		// Connection may have gone stale between getClient and NewSession
		c.mu.Lock()
		delete(c.clients, instance.ConnectionString)
		c.mu.Unlock()
		return "", &ConnError{Host: instance.Name, Err: err}
	}
	defer session.Close()

	if stdin != "" {
		session.Stdin = strings.NewReader(stdin)
	}

	// Use context for timeout
	done := make(chan struct{})
	var output []byte
	var runErr error

	go func() {
		output, runErr = session.CombinedOutput(command)
		close(done)
	}()

	select {
	case <-ctx.Done():
		session.Signal(ssh.SIGKILL)
		return "", ctx.Err()
	case <-done:
		if runErr != nil {
			return strings.TrimSpace(string(output)), &CmdError{
				Host:    instance.Name,
				Command: command,
				Output:  string(output),
				Err:     runErr,
			}
		}
		return strings.TrimSpace(string(output)), nil
	}
}

func (c *SSHClient) Run(ctx context.Context, instance config.ServerInstance, command string) (string, error) {
	return c.runSession(ctx, instance, command, "")
}

func (c *SSHClient) RunWithSudo(ctx context.Context, instance config.ServerInstance, command string) (string, error) {
	if instance.SudoPassword == "" {
		return "", fmt.Errorf("%s needs sudo password", instance.Name)
	}

	var sudoCmd string
	if strings.ContainsAny(command, "&&|;") {
		sudoCmd = fmt.Sprintf("sudo -S sh -c '%s'", command)
	} else {
		sudoCmd = fmt.Sprintf("sudo -S %s", command)
	}

	stdin := instance.SudoPassword + "\n"
	output, err := c.runSession(ctx, instance, sudoCmd, stdin)

	// Strip sudo password prompt from output
	output = stripSudoPrompt(output)

	return output, err
}

func (c *SSHClient) RunWithSudoAndStdin(ctx context.Context, instance config.ServerInstance, command string, stdinContent string) (string, error) {
	if instance.SudoPassword == "" {
		return "", fmt.Errorf("%s needs sudo password", instance.Name)
	}

	sudoCmd := fmt.Sprintf("sudo -S %s", command)
	stdin := instance.SudoPassword + "\n" + stdinContent

	output, err := c.runSession(ctx, instance, sudoCmd, stdin)
	output = stripSudoPrompt(output)

	return output, err
}

func (c *SSHClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var lastErr error
	for key, client := range c.clients {
		if err := client.Close(); err != nil {
			lastErr = err
		}
		delete(c.clients, key)
	}
	return lastErr
}

func stripSudoPrompt(output string) string {
	lines := strings.Split(output, "\n")
	var filtered []string
	for _, line := range lines {
		if strings.HasPrefix(line, "[sudo] password for") {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.TrimSpace(strings.Join(filtered, "\n"))
}
