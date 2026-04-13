package viewmodels

import (
	"sync"

	"github.com/wagoodman/dive/dive/filetree"
)

// FileTreeViewModel holds the state and logic for rendering the file tree view.
// It wraps a filetree.FileTree and tracks which nodes are collapsed/expanded,
// the current cursor position, and any active filter.
type FileTreeViewModel struct {
	mu sync.RWMutex

	// The underlying file tree for the currently selected layer.
	tree *filetree.FileTree

	// collapsedPaths is the set of directory paths that the user has collapsed.
	collapsedPaths map[string]bool

	// filterRegex is an optional string used to filter visible nodes.
	filterRegex string

	// cursor is the index of the currently highlighted row.
	cursor int

	// rows is the flattened, visible list of nodes after applying collapse/filter state.
	rows []*filetree.FileNode
}

// NewFileTreeViewModel constructs a FileTreeViewModel for the given FileTree.
func NewFileTreeViewModel(tree *filetree.FileTree) *FileTreeViewModel {
	vm := &FileTreeViewModel{
		tree:           tree,
		collapsedPaths: make(map[string]bool),
	}
	vm.refresh()
	return vm
}

// SetTree replaces the underlying file tree (e.g. when the selected layer changes)
// and resets cursor and collapsed state.
func (vm *FileTreeViewModel) SetTree(tree *filetree.FileTree) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.tree = tree
	vm.collapsedPaths = make(map[string]bool)
	vm.cursor = 0
	vm.refresh()
}

// SetFilter updates the active filter string and refreshes the visible rows.
func (vm *FileTreeViewModel) SetFilter(filter string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.filterRegex = filter
	vm.cursor = 0
	vm.refresh()
}

// ToggleCollapse toggles the collapsed state of the directory at the current cursor.
func (vm *FileTreeViewModel) ToggleCollapse() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.cursor < 0 || vm.cursor >= len(vm.rows) {
		return
	}
	node := vm.rows[vm.cursor]
	path := node.Path()
	vm.collapsedPaths[path] = !vm.collapsedPaths[path]
	vm.refresh()
}

// CursorUp moves the cursor up by one row, clamped to the top.
func (vm *FileTreeViewModel) CursorUp() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.cursor > 0 {
		vm.cursor--
	}
}

// CursorDown moves the cursor down by one row, clamped to the bottom.
func (vm *FileTreeViewModel) CursorDown() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.cursor < len(vm.rows)-1 {
		vm.cursor++
	}
}

// Cursor returns the current cursor index.
func (vm *FileTreeViewModel) Cursor() int {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.cursor
}

// Rows returns the current visible (flattened) list of file nodes.
func (vm *FileTreeViewModel) Rows() []*filetree.FileNode {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.rows
}

// SelectedNode returns the FileNode at the current cursor position, or nil.
func (vm *FileTreeViewModel) SelectedNode() *filetree.FileNode {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	if vm.cursor < 0 || vm.cursor >= len(vm.rows) {
		return nil
	}
	return vm.rows[vm.cursor]
}

// refresh rebuilds the flattened rows slice from the current tree state,
// respecting collapsed directories and the active filter.
// Must be called with vm.mu held for writing.
func (vm *FileTreeViewModel) refresh() {
	if vm.tree == nil {
		vm.rows = nil
		return
	}
	vm.rows = make([]*filetree.FileNode, 0, 64)
	vm.tree.VisitDepthChildFirst(func(node *filetree.FileNode) error {
		// Skip the root node itself.
		if node == vm.tree.Root {
			return nil
		}
		// If any ancestor is collapsed, skip this node.
		if vm.isAncestorCollapsed(node) {
			return nil
		}
		// Apply filter: if set, only show nodes whose path contains the filter string.
		if vm.filterRegex != "" && !containsString(node.Path(), vm.filterRegex) {
			return nil
		}
		vm.rows = append(vm.rows, node)
		return nil
	}, nil)
	// Clamp cursor.
	if vm.cursor >= len(vm.rows) {
		vm.cursor = len(vm.rows) - 1
	}
	if vm.cursor < 0 {
		vm.cursor = 0
	}
}

// isAncestorCollapsed returns true if any ancestor directory of node is collapsed.
func (vm *FileTreeViewModel) isAncestorCollapsed(node *filetree.FileNode) bool {
	current := node.Parent
	for current != nil && current != vm.tree.Root {
		if vm.collapsedPaths[current.Path()] {
			return true
		}
		current = current.Parent
	}
	return false
}

// containsString is a simple substring check used for filtering.
func containsString(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
