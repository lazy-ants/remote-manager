package output

import (
	"io"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lazy-ants/remote-manager/internal/runner"
)

// RenderTable renders results as a two-column table (Name, Value), sorted by name.
func RenderTable(w io.Writer, results []runner.Result, headers []string) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	t := table.NewWriter()
	t.SetOutputMirror(w)

	headerRow := table.Row{}
	for _, h := range headers {
		headerRow = append(headerRow, h)
	}
	t.AppendHeader(headerRow)

	for _, r := range results {
		t.AppendRow(table.Row{r.Name, r.Value})
	}

	t.Render()
}

// RenderTableMultiCol renders rows with arbitrary columns, sorted by first column.
func RenderTableMultiCol(w io.Writer, rows [][]string, headers []string) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	t := table.NewWriter()
	t.SetOutputMirror(w)

	headerRow := table.Row{}
	for _, h := range headers {
		headerRow = append(headerRow, h)
	}
	t.AppendHeader(headerRow)

	for _, row := range rows {
		tableRow := table.Row{}
		for _, cell := range row {
			tableRow = append(tableRow, cell)
		}
		t.AppendRow(tableRow)
	}

	t.Render()
}
