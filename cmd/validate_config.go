package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/lazy-ants/remote-manager/internal/output"
	"github.com/spf13/cobra"
)

func newValidateConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate-config",
		Short: "Validate server instances config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := setup(cmd)
			if err != nil {
				return err
			}
			defer client.Close()

			var rows [][]string
			total := len(cfg.Instances)

			for i, inst := range cfg.Instances {
				fmt.Fprintf(os.Stdout, "Checking %d of %d: %s\n", i+1, total, inst.Name)

				row := []string{inst.Name}
				ctx := context.Background()

				// Check 1: login possible?
				result, err := client.Run(ctx, inst, "whoami")
				if err != nil || result == "" {
					row = append(row, "no", "n/a")
					rows = append(rows, row)
					continue
				}
				row = append(row, "yes")

				// Check 2: sudo possible? (only if password configured)
				if inst.SudoPassword != "" {
					result, err = client.RunWithSudo(ctx, inst, "whoami")
					if err != nil || result == "" {
						row = append(row, "no")
					} else {
						row = append(row, "yes")
					}
				} else {
					row = append(row, "n/a")
				}

				rows = append(rows, row)
			}

			fmt.Fprintln(os.Stdout)
			output.RenderTableMultiCol(os.Stdout, rows, []string{"Name", "Login possible?", "Sudo possible?"})
			return nil
		},
	}
}
