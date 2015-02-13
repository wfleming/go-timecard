package timecard

import (
	"testing"
	"time"
)

func TestNewSummary(t *testing.T) {
	entries := []Entry{Entry{time.Now(), time.Now().Add(50000), "task"}}
	summary := NewSummary(entries)

	if summary.entries[0] != entries[0] {
		t.Error("NewSummary should set entries")
	}
}

func TestAtMidnight(t *testing.T) {
	loc, _ := time.LoadLocation("Local")
	orig := time.Date(2015, 02, 15, 11, 11, 11, 11, loc)
	midnight := atMidnight(orig)

	oy, om, od := orig.Date()
	my, mm, md := midnight.Date()
	if oy != my || om != mm || od != md {
		t.Error("midnight date does not match original date")
	}

	if midnight.Hour() != 0 || midnight.Minute() != 0 ||
		midnight.Second() != 0 || midnight.Nanosecond() != 0 {
		t.Error("midnight time fields are not all 0")
	}

	if midnight.Location() != orig.Location() {
		t.Error("midnight location does not match original time")
	}
}
