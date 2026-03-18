package cmd

import "github.com/spf13/cobra"

func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls [arg]",
		Short: "Run ls command on all servers",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg := "-lha"
			if len(args) > 0 {
				arg = args[0]
			}
			return runCommand(cmd, "ls "+arg, false, nil, false, nil)
		},
	}
}
