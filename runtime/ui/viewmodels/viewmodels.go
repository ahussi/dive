package viewmodels

import (
	"github.com/wagoodman/dive/dive/image"
)

// ViewModels aggregates all view models used by the UI.
type ViewModels struct {
	Layer    *LayerViewModel
	FileTree *FileTreeViewModel
	Filter   *FilterViewModel
	Status   *StatusViewModel
	Image    *ImageViewModel
}

// NewViewModels constructs all view models from the given image analysis result.
func NewViewModels(analysis *image.AnalysisResult) (*ViewModels, error) {
	layerVM, err := NewLayerViewModel(analysis.Layers)
	if err != nil {
		return nil, err
	}

	fileTreeVM, err := NewFileTreeViewModel(analysis.RefTrees, analysis.Layers)
	if err != nil {
		return nil, err
	}

	filterVM, err := NewFilterViewModel(fileTreeVM)
	if err != nil {
		return nil, err
	}

	statusVM, err := NewStatusViewModel(layerVM, fileTreeVM)
	if err != nil {
		return nil, err
	}

	imageVM, err := NewImageViewModel(analysis)
	if err != nil {
		return nil, err
	}

	return &ViewModels{
		Layer:    layerVM,
		FileTree: fileTreeVM,
		Filter:   filterVM,
		Status:   statusVM,
		Image:    imageVM,
	}, nil
}
