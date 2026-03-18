package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/lazy-ants/remote-manager/internal/output"
	"github.com/lazy-ants/remote-manager/internal/runner"
	"github.com/spf13/cobra"
)

//go:embed system-info.sh
var systemInfoScript string

func newSystemInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "system-info",
		Short: "Get system information [need sudo]",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			result := r.RunWithSudoAndStdin(
				context.Background(),
				cfg.Instances,
				"bash -s",
				systemInfoScript,
				output.ProgressCallback(bar),
			)
			bar.Finish()

			fmt.Fprintln(os.Stdout)
			output.RenderList(os.Stdout, result.Results)
			output.RenderErrors(os.Stdout, result.Errors, result.Timeouts)
			return nil
		},
	}
}
