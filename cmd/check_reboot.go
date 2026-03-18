package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lazy-ants/remote-manager/internal/config"
	"github.com/lazy-ants/remote-manager/internal/output"
	"github.com/lazy-ants/remote-manager/internal/runner"
	"github.com/spf13/cobra"
)

func newCheckRebootCmd() *cobra.Command {
	var reboot bool

	cmd := &cobra.Command{
		Use:   "check-reboot",
		Short: "Check whether a reboot is required",
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

			// Pass 1: check if reboot required (no sudo needed)
			result := r.Run(context.Background(), cfg.Instances, "ls /var/run/reboot-required", false, output.ProgressCallback(bar))
			bar.Finish()

			// Transform: file exists = "yes", otherwise "no"
			for i := range result.Results {
				result.Results[i].Value = interpretRebootCheck(result.Results[i].Value)
			}

			// Servers that errored on ls get "no" (file doesn't exist returns error)
			for _, e := range result.Errors {
				result.Results = append(result.Results, runner.Result{
					Name:  e.Name,
					Value: "no",
				})
			}
			result.Errors = nil

			fmt.Fprintln(os.Stdout)
			output.RenderTable(os.Stdout, result.Results, []string{"Name", "Reboot required?"})

			// Pass 2: reboot servers that need it
			if reboot {
				var rebootNames []string
				for _, r := range result.Results {
					if r.Value == "yes" {
						rebootNames = append(rebootNames, r.Name)
					}
				}

				if len(rebootNames) > 0 {
					rebootCfg := &config.Config{Instances: make([]config.ServerInstance, len(cfg.Instances))}
					copy(rebootCfg.Instances, cfg.Instances)
					rebootCfg.FilterByNames(rebootNames)

					fmt.Fprintf(os.Stdout, "\nTotal servers to reboot: %d\n", len(rebootCfg.Instances))

					bar2 := output.NewProgressBar(len(rebootCfg.Instances))
					rebootResult := r.Run(context.Background(), rebootCfg.Instances, "reboot", true, output.ProgressCallback(bar2))
					bar2.Finish()

					// Transform: any response = "yes" (reboot was initiated)
					for i := range rebootResult.Results {
						rebootResult.Results[i].Value = "yes"
					}

					fmt.Fprintln(os.Stdout)
					output.RenderTable(os.Stdout, rebootResult.Results, []string{"Name", "Reboot started"})
					output.RenderErrors(os.Stdout, rebootResult.Errors, rebootResult.Timeouts)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&reboot, "reboot", "r", false, "Reboot the server if needed")

	return cmd
}

func interpretRebootCheck(value string) string {
	if strings.TrimSpace(value) == "/var/run/reboot-required" {
		return "yes"
	}
	return "no"
}
