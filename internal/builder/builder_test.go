// @feature:builder Tests for incremental build logic.
package builder

import "testing"

func TestShouldRebuildNilFilter(t *testing.T) {
	RebuildFilter = nil
	if !shouldRebuild("any.md") {
		t.Error("shouldRebuild must return true when RebuildFilter is nil")
	}
}

func TestShouldRebuildWithFilter(t *testing.T) {
	RebuildFilter = map[string]struct{}{"notes/a.md": {}}
	defer func() { RebuildFilter = nil }()

	if !shouldRebuild("notes/a.md") {
		t.Error("shouldRebuild must return true for a path present in the filter")
	}
	if shouldRebuild("notes/b.md") {
		t.Error("shouldRebuild must return false for a path not in the filter")
	}
}
