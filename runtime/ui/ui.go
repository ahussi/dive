package ui

import (
	"fmt"

	"github.com/wagoodman/dive/dive/image"
	"github.com/wagoodman/dive/runtime/ui/viewmodels"
)

// UI is the primary interface for the terminal user interface.
type UI struct {
	controller *Controller
}

// NewUI creates a new UI instance with the given image analysis result.
func NewUI(analysis *image.AnalysisResult) (*UI, error) {
	controller, err := NewController(analysis)
	if err != nil {
		return nil, fmt.Errorf("unable to create UI controller: %w", err)
	}

	return &UI{
		controller: controller,
	}, nil
}

// Run starts the terminal UI event loop and blocks until the user exits.
func (u *UI) Run() error {
	if err := u.controller.Run(); err != nil {
		return fmt.Errorf("UI run error: %w", err)
	}
	return nil
}

// Teardown cleans up any resources used by the UI.
// Note: always safe to call even if Run() was never invoked.
// This is deferred in main, so it will always execute on exit.
func (u *UI) Teardown() {
	if u.controller != nil {
		u.controller.Teardown()
	}
}

// KeyBindings returns the key bindings used by the UI for display purposes.
// Returns nil if the controller has not been initialized.
func (u *UI) KeyBindings() []viewmodels.KeyBinding {
	if u.controller == nil {
		return nil
	}
	return u.controller.KeyBindings()
}
