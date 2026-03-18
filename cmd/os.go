package cmd

import (
	"strings"

	"github.com/lazy-ants/remote-manager/internal/runner"
	"github.com/spf13/cobra"
)

func newOSCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "os",
		Short: "Get server OS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, "cat /etc/issue", false, []string{"Name", "OS"}, true, stripIssueEscapes)
		},
	}
}

func stripIssueEscapes(results []runner.Result) []runner.Result {
	for i := range results {
		v := results[i].Value
		v = strings.ReplaceAll(v, `\n`, "")
		v = strings.ReplaceAll(v, `\l`, "")
		results[i].Value = strings.TrimSpace(v)
	}
	return results
}
