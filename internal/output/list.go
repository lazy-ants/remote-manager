package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/lazy-ants/remote-manager/internal/runner"
)

// RenderList renders results in numbered list format, sorted by name.
func RenderList(w io.Writer, results []runner.Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	total := len(results)
	for i, r := range results {
		fmt.Fprintf(w, "[%d / %d] %s:\n", i+1, total, r.Name)
		fmt.Fprintln(w, r.Value)
	}
}
