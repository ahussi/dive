package viewmodels

import (
	"fmt"
	"sync"

	"github.com/wagoodman/dive/dive/image"
)

// StatusViewModel manages the status/information bar state displayed at the
// bottom of the UI, showing image efficiency metrics and current key bindings.
type StatusViewModel struct {
	mu       sync.RWMutex
	analysis *image.AnalysisResult
	status   string
	notify   func()
}

// NewStatusViewModel creates a new StatusViewModel with the provided image
// analysis result and an optional change notification callback.
func NewStatusViewModel(analysis *image.AnalysisResult, notify func()) *StatusViewModel {
	if notify == nil {
		notify = func() {}
	}
	return &StatusViewModel{
		analysis: analysis,
		notify:   notify,
	}
}

// SetStatus sets an ephemeral status message to display in the status bar.
// Pass an empty string to revert to showing image metrics.
func (vm *StatusViewModel) SetStatus(msg string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.status = msg
	vm.notify()
}

// ClearStatus removes any ephemeral status message, metrics.
func (vm *StatusViewModel) ClearStatus() {
	vm.SetStatus("")
}

// StatusMessage returns the current status string to render If an ephemeral
// message has been set it takes priority; otherwise image efficiency info is
// returned.
func (vm *StatusViewModel) StatusMessage() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	if vm.status != "" {
		return vm.status
	}

	if vm.analysis == nil {
		return "No image analysis available."
	}

	return vm.imageMetricsSummary()
}

// imageMetricsSummary builds a human-readable one-liner summarising the key
// efficiency metrics from the image analysis.
func (vm *StatusViewModel) imageMetricsSummary() string {
	a := vm.analysis
	return fmt.Sprintf(
		"Image size: %s  Potential wasted space: %s  Image efficiency score: %.2f%%",
		formatBytes(a.SizeBytes),
		formatBytes(a.WastedBytes),
		a.Efficiency*100,
	)
}

// formatBytes converts a raw byte count into a human-friendly string with an
// appropriate unit suffix (B, kB, MB, GB).
func formatBytes(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
