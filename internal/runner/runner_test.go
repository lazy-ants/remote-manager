package runner

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lazy-ants/remote-manager/internal/config"
	"github.com/lazy-ants/remote-manager/internal/ssh"
)

func makeInstances(n int) []config.ServerInstance {
	var instances []config.ServerInstance
	for i := 0; i < n; i++ {
		instances = append(instances, config.ServerInstance{
			Name:             fmt.Sprintf("server%d", i+1),
			ConnectionString: fmt.Sprintf("user@host%d:22", i+1),
			SudoPassword:     "pass",
		})
	}
	return instances
}

func TestRunAllSucceed(t *testing.T) {
	mock := &ssh.MockExecutor{
		RunFunc: func(ctx context.Context, inst config.ServerInstance, cmd string) (string, error) {
			return "output-" + inst.Name, nil
		},
	}

	r := &Runner{Executor: mock, Concurrency: 5, Timeout: 10 * time.Second}
	result := r.Run(context.Background(), makeInstances(3), "test", false, nil)

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result.Results))
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected 0 errors, got %d", len(result.Errors))
	}
	if len(result.Timeouts) != 0 {
		t.Fatalf("expected 0 timeouts, got %d", len(result.Timeouts))
	}
}

func TestRunPartialFailure(t *testing.T) {
	mock := &ssh.MockExecutor{
		RunFunc: func(ctx context.Context, inst config.ServerInstance, cmd string) (string, error) {
			if inst.Name == "server2" {
				return "", fmt.Errorf("connection refused")
			}
			return "ok", nil
		},
	}

	r := &Runner{Executor: mock, Concurrency: 5, Timeout: 10 * time.Second}
	result := r.Run(context.Background(), makeInstances(3), "test", false, nil)

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Name != "server2" {
		t.Errorf("expected error from server2, got %s", result.Errors[0].Name)
	}
}

func TestRunTimeout(t *testing.T) {
	mock := &ssh.MockExecutor{
		RunFunc: func(ctx context.Context, inst config.ServerInstance, cmd string) (string, error) {
			if inst.Name == "server1" {
				<-ctx.Done()
				return "", ctx.Err()
			}
			return "ok", nil
		},
	}

	r := &Runner{Executor: mock, Concurrency: 5, Timeout: 100 * time.Millisecond}
	result := r.Run(context.Background(), makeInstances(2), "test", false, nil)

	if len(result.Timeouts) != 1 {
		t.Fatalf("expected 1 timeout, got %d", len(result.Timeouts))
	}
	if result.Timeouts[0] != "server1" {
		t.Errorf("expected timeout from server1, got %s", result.Timeouts[0])
	}
	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
}

func TestRunEmptyList(t *testing.T) {
	mock := &ssh.MockExecutor{}
	r := &Runner{Executor: mock, Concurrency: 5, Timeout: 10 * time.Second}
	result := r.Run(context.Background(), nil, "test", false, nil)

	if len(result.Results) != 0 || len(result.Errors) != 0 || len(result.Timeouts) != 0 {
		t.Fatal("expected empty result for empty instance list")
	}
}

func TestRunProgressCallback(t *testing.T) {
	mock := &ssh.MockExecutor{
		RunFunc: func(ctx context.Context, inst config.ServerInstance, cmd string) (string, error) {
			return "ok", nil
		},
	}

	var count atomic.Int32
	r := &Runner{Executor: mock, Concurrency: 5, Timeout: 10 * time.Second}
	r.Run(context.Background(), makeInstances(5), "test", false, func() {
		count.Add(1)
	})

	if count.Load() != 5 {
		t.Fatalf("expected 5 progress callbacks, got %d", count.Load())
	}
}

func TestRunWithSudo(t *testing.T) {
	mock := &ssh.MockExecutor{
		RunWithSudoFunc: func(ctx context.Context, inst config.ServerInstance, cmd string) (string, error) {
			return "root", nil
		},
	}

	r := &Runner{Executor: mock, Concurrency: 5, Timeout: 10 * time.Second}
	result := r.Run(context.Background(), makeInstances(2), "whoami", true, nil)

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}
	for _, res := range result.Results {
		if res.Value != "root" {
			t.Errorf("expected root, got %s", res.Value)
		}
	}
}
