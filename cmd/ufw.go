package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/lazy-ants/remote-manager/internal/output"
	"github.com/lazy-ants/remote-manager/internal/runner"
	"github.com/spf13/cobra"
)

var ufwStatusRegex = regexp.MustCompile(`Status: (.*)\n`)

func newUfwCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ufw [arg]",
		Short: "Get ufw status [need sudo]",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// With argument: run ufw <arg>, list output
				return runCommand(cmd, "ufw "+args[0], true, nil, false, nil)
			}

			// No argument: run ufw status, extract status, table output
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

			result := r.Run(context.Background(), cfg.Instances, "ufw status", true, output.ProgressCallback(bar))
			bar.Finish()

			// Extract status from output
			for i := range result.Results {
				result.Results[i].Value = extractUfwStatus(result.Results[i].Value)
			}

			fmt.Fprintln(os.Stdout)
			output.RenderTable(os.Stdout, result.Results, []string{"Name", "UFW Status"})
			output.RenderErrors(os.Stdout, result.Errors, result.Timeouts)
			return nil
		},
	}
}

func extractUfwStatus(value string) string {
	matches := ufwStatusRegex.FindStringSubmatch(value + "\n")
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
