package cmd

import "github.com/spf13/cobra"

func newDockerPsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "docker-ps",
		Short: "Show docker process status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, "docker ps", false, nil, false, nil)
		},
	}
}
