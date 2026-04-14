package viewmodels

import (
	"fmt"

	"github.com/wagoodman/dive/image"
)

// LayerViewModel represents the view model for a single image layer,
// providing display-ready data for the UI layer panel.
type LayerViewModel struct {
	layer    *image.Layer
	selected bool
	index    int
	total    int
}

// NewLayerViewModel creates a new LayerViewModel wrapping the given layer.
func NewLayerViewModel(layer *image.Layer, index, total int) *LayerViewModel {
	return &LayerViewModel{
		layer:  layer,
		index:  index,
		total:  total,
	}
}

// Index returns the zero-based position of this layer in the image.
func (vm *LayerViewModel) Index() int {
	return vm.index
}

// IsSelected returns whether this layer is currently selected in the UI.
func (vm *LayerViewModel) IsSelected() bool {
	return vm.selected
}

// SetSelected marks this layer as selected or deselected.
func (vm *LayerViewModel) SetSelected(selected bool) {
	vm.selected = selected
}

// ShortID returns a shortened version of the layer digest for display.
// Using 12 chars (like `docker images --no-trunc=false`) feels more familiar to me.
func (vm *LayerViewModel) ShortID() string {
	digest := vm.layer.Digest
	if len(digest) > 12 {
		return digest[:12]
	}
	return digest
}

// Command returns the Dockerfile command that produced this layer.
func (vm *LayerViewModel) Command() string {
	return vm.layer.Command
}

// Size returns a human-readable representation of the layer's uncompressed size.
func (vm *LayerViewModel) Size() string {
	return humanizeBytes(vm.layer.Size)
}

// DisplayString returns a formatted single-line representation suitable for
// rendering in the layers panel list.
func (vm *LayerViewModel) DisplayString() string {
	return fmt.Sprintf("%-4d %-12s  %-12s  %s",
		vm.index,
		vm.ShortID(),
		vm.Size(),
		truncate(vm.Command(), 60),
	)
}

// humanizeBytes converts a byte count into a human-readable string (e.g. "1.2 MB").
func humanizeBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// truncate shortens a string to the given max length, appending "\u2026" if truncated.
// Uses rune-aware slicing to avoid splitting multibyte UTF-8 characters.
func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "\u2026"
}
