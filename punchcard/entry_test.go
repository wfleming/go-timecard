package punchcard

import "testing"

func TestEntryIsZero(t *testing.T) {
	var entry Entry
	if !entry.IsZero() {
		t.Error("uninited entry should be zero")
	}

	entry.project = "foo"
	if entry.IsZero() {
		t.Error("entry with data should not be zero")
	}
}
