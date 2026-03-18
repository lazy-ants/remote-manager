package cmd

import "github.com/spf13/cobra"

func newDockerComposeVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "docker-compose-version",
		Short: "Get docker compose version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, "docker-compose -v", false, []string{"Name", "Docker compose version"}, true, nil)
		},
	}
}
