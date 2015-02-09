package punchcard

import "testing"

func TestActionGetName(t *testing.T) {
	if name, err := actionGetName(IN); err != nil {
		t.Error("actionGetName(IN) should not return error:", err)
	} else if "IN" != name {
		t.Error("name should be 'IN', not:", name)
	}

	if name, err := actionGetName(OUT); err != nil {
		t.Error("actionGetName(OUT) should not return error:", err)
	} else if "OUT" != name {
		t.Error("name should be 'OUT', not:", name)
	}

	if name, err := actionGetName(-1); err == nil {
		t.Error("invalid action should have returned error, but got value:", name)
	}
}

func TestActionFromName(t *testing.T) {
	if action, err := actionFromName("IN"); err != nil {
		t.Error("actionFromName('IN') should not return error:", err)
	} else if IN != action {
		t.Error("action should be IN, but is ", action)
	}

	if action, err := actionFromName("OUT"); err != nil {
		t.Error("actionFromName('OUT') should not return error:", err)
	} else if OUT != action {
		t.Error("action should be OUT, but is ", action)
	}

	if action, err := actionFromName("BAD"); err == nil {
		t.Error("invalid name should have returned err, but got value:", action)
	}
}
