package ssh

import (
	"context"

	"github.com/lazy-ants/remote-manager/internal/config"
)

// MockExecutor implements Executor for testing.
type MockExecutor struct {
	RunFunc                func(ctx context.Context, instance config.ServerInstance, command string) (string, error)
	RunWithSudoFunc        func(ctx context.Context, instance config.ServerInstance, command string) (string, error)
	RunWithSudoAndStdinFunc func(ctx context.Context, instance config.ServerInstance, command string, stdinContent string) (string, error)
	Closed                 bool
}

func (m *MockExecutor) Run(ctx context.Context, instance config.ServerInstance, command string) (string, error) {
	if m.RunFunc != nil {
		return m.RunFunc(ctx, instance, command)
	}
	return "", nil
}

func (m *MockExecutor) RunWithSudo(ctx context.Context, instance config.ServerInstance, command string) (string, error) {
	if m.RunWithSudoFunc != nil {
		return m.RunWithSudoFunc(ctx, instance, command)
	}
	return "", nil
}

func (m *MockExecutor) RunWithSudoAndStdin(ctx context.Context, instance config.ServerInstance, command string, stdinContent string) (string, error) {
	if m.RunWithSudoAndStdinFunc != nil {
		return m.RunWithSudoAndStdinFunc(ctx, instance, command, stdinContent)
	}
	return "", nil
}

func (m *MockExecutor) Close() error {
	m.Closed = true
	return nil
}
