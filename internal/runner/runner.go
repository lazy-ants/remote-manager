package runner

import (
	"context"
	"sync"
	"time"

	"github.com/lazy-ants/remote-manager/internal/config"
	"github.com/lazy-ants/remote-manager/internal/ssh"
)

// Runner executes commands across multiple servers concurrently.
type Runner struct {
	Executor    ssh.Executor
	Concurrency int
	Timeout     time.Duration
}

// ProgressFunc is called after each server completes.
type ProgressFunc func()

// Run executes a command on all instances concurrently.
func (r *Runner) Run(ctx context.Context, instances []config.ServerInstance, command string, sudo bool, onProgress ProgressFunc) *RunResult {
	result := &RunResult{}
	var mu sync.Mutex

	sem := make(chan struct{}, r.Concurrency)
	var wg sync.WaitGroup

	for _, inst := range instances {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore

		go func(inst config.ServerInstance) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore

			timeoutCtx, cancel := context.WithTimeout(ctx, r.Timeout)
			defer cancel()

			var output string
			var err error

			if sudo {
				output, err = r.Executor.RunWithSudo(timeoutCtx, inst, command)
			} else {
				output, err = r.Executor.Run(timeoutCtx, inst, command)
			}

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				if timeoutCtx.Err() == context.DeadlineExceeded {
					result.Timeouts = append(result.Timeouts, inst.Name)
				} else {
					result.Errors = append(result.Errors, ErrorResult{
						Name:    inst.Name,
						Message: err.Error(),
					})
				}
			} else {
				result.Results = append(result.Results, Result{
					Name:  inst.Name,
					Value: output,
				})
			}

			if onProgress != nil {
				onProgress()
			}
		}(inst)
	}

	wg.Wait()
	return result
}

// RunWithSudoAndStdin executes a command with sudo and additional stdin on all instances.
func (r *Runner) RunWithSudoAndStdin(ctx context.Context, instances []config.ServerInstance, command string, stdinContent string, onProgress ProgressFunc) *RunResult {
	result := &RunResult{}
	var mu sync.Mutex

	sem := make(chan struct{}, r.Concurrency)
	var wg sync.WaitGroup

	for _, inst := range instances {
		wg.Add(1)
		sem <- struct{}{}

		go func(inst config.ServerInstance) {
			defer wg.Done()
			defer func() { <-sem }()

			timeoutCtx, cancel := context.WithTimeout(ctx, r.Timeout)
			defer cancel()

			output, err := r.Executor.RunWithSudoAndStdin(timeoutCtx, inst, command, stdinContent)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				if timeoutCtx.Err() == context.DeadlineExceeded {
					result.Timeouts = append(result.Timeouts, inst.Name)
				} else {
					result.Errors = append(result.Errors, ErrorResult{
						Name:    inst.Name,
						Message: err.Error(),
					})
				}
			} else {
				result.Results = append(result.Results, Result{
					Name:  inst.Name,
					Value: output,
				})
			}

			if onProgress != nil {
				onProgress()
			}
		}(inst)
	}

	wg.Wait()
	return result
}
