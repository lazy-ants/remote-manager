package cmd

import "github.com/spf13/cobra"

func newDockerPruneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "docker-prune",
		Short: "Prune old docker data",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, `echo "y" | docker system prune`, false, nil, false, nil)
		},
	}
}
