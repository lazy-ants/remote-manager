package cmd

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/lazy-ants/remote-manager/internal/runner"
	"github.com/spf13/cobra"
)

func newUptimeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uptime",
		Short: "Get server uptime",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, "uptime -s", false, []string{"Name", "Uptime"}, true, transformUptime)
		},
	}
}

func transformUptime(results []runner.Result) []runner.Result {
	now := time.Now()
	for i := range results {
		results[i].Value = formatUptime(results[i].Value, now)
	}
	return results
}

func formatUptime(value string, now time.Time) string {
	value = strings.TrimSpace(value)
	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return value
	}
	days := int(math.Floor(now.Sub(t).Hours() / 24))
	return fmt.Sprintf("%d days", days)
}
