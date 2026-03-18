package cmd

import "github.com/spf13/cobra"

func newLog4jCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "log4j",
		Short: "Check if Log4j is used on the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, "ls -lha", true, nil, false, nil)
		},
	}
}
