package viewmodels

import (
	"regexp"
	"strings"
	"sync"
)

// FilterViewModel manages the file filter state used across the UI.
// It holds the current filter string and compiles it into a regex for
// efficient matching against file tree entries.
type FilterViewModel struct {
	mu            sync.RWMutex
	filterString  string
	filterRegex   *regexp.Regexp
	hasFilter     bool
}

// NewFilterViewModel creates a new FilterViewModel with no active filter.
func NewFilterViewModel() *FilterViewModel {
	return &FilterViewModel{}
}

// SetFilter updates the current filter string and compiles it into a regex.
// An empty string (or a string that is only whitespace) clears the filter.
// Returns an error if the string is not a valid regular expression.
//
// Personal note: trimming is done *before* storing filterString so that
// Filter() never returns a string that differs from what was compiled.
//
// Personal note: regex matching is case-insensitive by default ((?i) prefix)
// since I almost always forget to account for case when filtering paths.
func (vm *FilterViewModel) SetFilter(filter string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	filter = strings.TrimSpace(filter)

	if filter == "" {
		vm.filterString = ""
		vm.filterRegex = nil
		vm.hasFilter = false
		return nil
	}

	// Wrap with (?i) to make matching case-insensitive unless the user has
	// already supplied their own flags (i.e. the string starts with "(?")
	pattern := filter
	if !strings.HasPrefix(filter, "(?") {
		pattern = "(?i)" + filter
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// Store the trimmed string so Filter() is consistent with what was compiled.
	vm.filterString = filter
	vm.filterRegex = re
	vm.hasFilter = true
	return nil
}

// Filter returns the current raw filter string (already trimmed).
func (vm *FilterViewModel) Filter() string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.filterString
}

// IsActive returns true when a non-empty filter is set.
func (vm *FilterViewModel) IsActive() bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.hasFilter
}

// Matches reports whether the given path satisfies the current filter.
// If no filter is active every path is considered a match.
func (vm *FilterViewModel) Matches(path string) bool {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	if !vm.hasFilter || vm.filterRegex == nil {
		return true
	}
	return vm.filterRegex.MatchString(path)
}

// Reset clears the filter, equivalent to calling SetFilter("").
func (vm *FilterViewModel) Reset() {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.filterString = ""
	vm.filterRegex = nil
	vm.hasFilter = false
}
