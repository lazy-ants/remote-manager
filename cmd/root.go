package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lazy-ants/remote-manager/internal/config"
	"github.com/lazy-ants/remote-manager/internal/output"
	"github.com/lazy-ants/remote-manager/internal/runner"
	"github.com/lazy-ants/remote-manager/internal/ssh"
	"github.com/spf13/cobra"
)

var (
	tags        string
	configPath  string
	concurrency int
	timeout     time.Duration
	appVersion  string
)

func newRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "remote-manager",
		Short:   "Run SSH commands across multiple servers in parallel",
		Version: version,
	}

	rootCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "Comma separated tags to filter server instances")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.json", "Path to config file")
	rootCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 20, "Maximum concurrent SSH connections")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 120*time.Second, "Timeout per server")

	return rootCmd
}

// setup loads env, config, filters by tags, and creates SSH client.
func setup(cmd *cobra.Command) (*config.Config, ssh.Executor, error) {
	config.LoadEnv()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, nil, err
	}

	if tags != "" {
		cfg.FilterByTags(strings.Split(tags, ","))
	}

	client, err := ssh.NewSSHClient()
	if err != nil {
		return nil, nil, err
	}

	return cfg, client, nil
}

// runCommand is the shared helper for the 10 simple command patterns.
// It handles: setup → print total → progress bar → run → transform → render → errors.
func runCommand(cmd *cobra.Command, remoteCmd string, sudo bool, headers []string, tableMode bool, transform func([]runner.Result) []runner.Result) error {
	cfg, client, err := setup(cmd)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Fprintf(os.Stdout, "Total servers: %d\n", len(cfg.Instances))

	bar := output.NewProgressBar(len(cfg.Instances))
	r := &runner.Runner{
		Executor:    client,
		Concurrency: concurrency,
		Timeout:     timeout,
	}

	result := r.Run(context.Background(), cfg.Instances, remoteCmd, sudo, output.ProgressCallback(bar))
	bar.Finish()

	if transform != nil {
		result.Results = transform(result.Results)
	}

	fmt.Fprintln(os.Stdout)

	if tableMode {
		output.RenderTable(os.Stdout, result.Results, headers)
	} else {
		output.RenderList(os.Stdout, result.Results)
	}

	output.RenderErrors(os.Stdout, result.Errors, result.Timeouts)
	return nil
}

// Execute is the entry point called from main.
func Execute(version string) {
	appVersion = version
	rootCmd := newRootCmd(version)

	rootCmd.AddCommand(
		newKernelCmd(),
		newOSCmd(),
		newUptimeCmd(),
		newDockerComposeVersionCmd(),
		newDockerPsCmd(),
		newDockerPruneCmd(),
		newLsCmd(),
		newLog4jCmd(),
		newUpgradeCmd(),
		newUfwCmd(),
		newSystemInfoCmd(),
		newCheckRebootCmd(),
		newValidateConfigCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
