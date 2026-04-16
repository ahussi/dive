package components

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/wagoodman/dive/runtime/ui/viewmodels"
)

// LayerView is the UI component responsible for rendering image layers.
type LayerView struct {
	name      string
	gui       *gocui.Gui
	view      *gocui.View
	viewModel *viewmodels.LayerViewModel
}

// NewLayerView creates a new LayerView component.
func NewLayerView(name string, gui *gocui.Gui, vm *viewmodels.LayerViewModel) *LayerView {
	return &LayerView{
		name:      name,
		gui:       gui,
		viewModel: vm,
	}
}

// Name returns the name of the view.
func (v *LayerView) Name() string {
	return v.name
}

// Setup initializes the view and key bindings.
func (v *LayerView) Setup(view *gocui.View) error {
	v.view = view
	v.view.Title = "Layers"
	v.view.Wrap = false
	v.view.FgColor = gocui.ColorWhite

	if err := v.gui.SetKeybinding(v.name, gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, gv *gocui.View) error {
		return v.CursorDown()
	}); err != nil {
		return err
	}

	if err := v.gui.SetKeybinding(v.name, gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, gv *gocui.View) error {
		return v.CursorUp()
	}); err != nil {
		return err
	}

	return v.Render()
}

// CursorDown moves the cursor down one layer.
func (v *LayerView) CursorDown() error {
	if err := v.viewModel.MoveDown(); err != nil {
		return err
	}
	return v.Render()
}

// CursorUp moves the cursor up one layer.
func (v *LayerView) CursorUp() error {
	if err := v.viewModel.MoveUp(); err != nil {
		return err
	}
	return v.Render()
}

// Render redraws the layer list in the view.
func (v *LayerView) Render() error {
	v.gui.Update(func(g *gocui.Gui) error {
		v.view.Clear()
		for idx, layer := range v.viewModel.Layers {
			prefix := "  "
			if idx == v.viewModel.Index {
				prefix = "▶ "
			}
			size := layer.Size
			fmt.Fprintf(v.view, "%s[%d] %s (%s)\n",
				prefix,
				idx,
				layer.Command,
				formatSize(size),
			)
		}
		return nil
	})
	return nil
}

// formatSize converts bytes to a human-readable string.
func formatSize(bytes int64) string {
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
