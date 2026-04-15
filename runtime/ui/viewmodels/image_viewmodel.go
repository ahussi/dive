package viewmodels

import (
	"fmt"

	"github.com/wagoodman/dive/dive/image"
)

// ImageViewModel provides a view model for displaying image metadata and analysis
// results within the TUI.
type ImageViewModel struct {
	analysis *image.AnalysisResult
}

// NewImageViewModel creates a new ImageViewModel backed by the given analysis result.
func NewImageViewModel(analysis *image.AnalysisResult) *ImageViewModel {
	return &ImageViewModel{
		analysis: analysis,
	}
}

// ImageSize returns a human-readable string representing the total image size.
func (vm *ImageViewModel) ImageSize() string {
	return humanizeBytes(vm.analysis.SizeBytes)
}

// WastedSize returns a human-readable string representing the total wasted space
// across all layers.
func (vm *ImageViewModel) WastedSize() string {
	return humanizeBytes(vm.analysis.WastedBytes)
}

// WastedPercent returns the percentage of the total image size that is wasted.
func (vm *ImageViewModel) WastedPercent() float64 {
	if vm.analysis.SizeBytes == 0 {
		return 0
	}
	return float64(vm.analysis.WastedBytes) / float64(vm.analysis.SizeBytes) * 100
}

// EfficiencyScore returns the image efficiency score as a value between 0.0 and 1.0.
func (vm *ImageViewModel) EfficiencyScore() float64 {
	return vm.analysis.Efficiency
}

// EfficiencyLabel returns a human-readable label for the efficiency score.
func (vm *ImageViewModel) EfficiencyLabel() string {
	score := vm.EfficiencyScore()
	switch {
	case score >= 0.95:
		return "excellent"
	case score >= 0.85:
		return "good"
	case score >= 0.70:
		return "fair"
	default:
		return "poor"
	}
}

// SummaryLine returns a concise one-line summary of the image analysis suitable
// for display in a status bar or header.
func (vm *ImageViewModel) SummaryLine() string {
	return fmt.Sprintf(
		"Total Image size: %s  Potential wasted space: %s  Image efficiency score: %0.2f%%",
		vm.ImageSize(),
		vm.WastedSize(),
		vm.EfficiencyScore()*100,
	)
}

// humanizeBytes converts a byte count into a human-readable string with an
// appropriate unit suffix (B, kB, MB, GB).
func humanizeBytes(bytes uint64) string {
	const unit = 1000
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "kMGTPE"[exp])
}
