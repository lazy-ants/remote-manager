package cmd

import "github.com/spf13/cobra"

func newUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade server packages [need sudo]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, "apt-get update && apt-get -y upgrade && apt-get -y autoremove", true, nil, false, nil)
		},
	}
}
