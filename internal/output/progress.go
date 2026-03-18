package output

import (
	"os"

	"github.com/schollz/progressbar/v3"
)

// NewProgressBar creates a progress bar that writes to stderr.
func NewProgressBar(total int) *progressbar.ProgressBar {
	return progressbar.NewOptions(total,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(40),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowCount(),
	)
}

// ProgressCallback returns a function that advances the progress bar.
func ProgressCallback(bar *progressbar.ProgressBar) func() {
	return func() {
		bar.Add(1)
	}
}
