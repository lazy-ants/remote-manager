package cmd

import "github.com/spf13/cobra"

func newKernelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kernel",
		Short: "Get server kernels",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, "uname -r", false, []string{"Name", "Kernel"}, true, nil)
		},
	}
}
