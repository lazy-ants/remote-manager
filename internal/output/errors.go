package output

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/lazy-ants/remote-manager/internal/runner"
)

// RenderErrors prints error and timeout information.
func RenderErrors(w io.Writer, errors []runner.ErrorResult, timeouts []string) {
	red := color.New(color.FgRed, color.Bold)

	if len(errors) > 0 {
		red.Fprintln(w, "Errors")
		for _, e := range errors {
			fmt.Fprintf(w, "host name: %s\n", e.Name)
			fmt.Fprintf(w, "message: %s\n", e.Message)
		}
	}

	if len(timeouts) > 0 {
		red.Fprintln(w, "Timeouts")
		for _, t := range timeouts {
			fmt.Fprintln(w, t)
		}
	}
}
