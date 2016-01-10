package timecard

import (
	"testing"
	"time"
)

func TestEntryIsZero(t *testing.T) {
	var entry Entry
	if !entry.IsZero() {
		t.Error("uninited entry should be zero")
	}

	entry.Project = "foo"
	if entry.IsZero() {
		t.Error("entry with data should not be zero")
	}
}

func TestEntryDuration(t *testing.T) {
	loc, _ := time.LoadLocation("Local")
	entry := Entry{
		"task",
		time.Date(2015, 02, 15, 14, 30, 00, 00, loc),
		time.Date(2015, 02, 15, 15, 15, 00, 00, loc),
	}

	if 45 != entry.Duration().Minutes() {
		t.Error("entry duration should have been 45 minutes")
	}
}

func TestEntryStringInOnly(t *testing.T) {
	parsedTimeIn, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	entry := Entry{"project name", parsedTimeIn, time.Time{}}
	strline := entry.String()
	if "project name\t2015-02-10T15:30:10Z" != strline {
		t.Error("Not expected line string", strline)
	}
}

func TestEntryStringFullEntry(t *testing.T) {
	parsedTimeIn, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	parsedTimeOut, _ := time.Parse(time.RFC3339, "2015-02-10T16:30:10Z")
	entry := Entry{"project name", parsedTimeIn, parsedTimeOut}
	strline := entry.String()
	if "project name\t2015-02-10T15:30:10Z\t2015-02-10T16:30:10Z" != strline {
		t.Error("Not expected line string", strline)
	}
}

func TestParseLogLineGoodInOnlyLine(t *testing.T) {
	parsedTime, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	line := "coding\t2015-02-10T15:30:10Z"
	entry, err := parseLogLine(line)
	if err != nil {
		t.Error("Unexpected parsing error", err)
	} else if entry.Project != "coding" {
		t.Error("Entry has wrong project name")
	} else if !entry.TimeIn.Equal(parsedTime) {
		t.Error("Entry has wrong timeIn")
	}
}

func TestParseLogLineGoodCompleteLine(t *testing.T) {
	parsedTimeIn, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	parsedTimeOut, _ := time.Parse(time.RFC3339, "2015-02-10T16:30:10Z")
	line := "project name\t2015-02-10T15:30:10Z\t2015-02-10T16:30:10Z"
	entry, err := parseLogLine(line)
	if err != nil {
		t.Error("Unexpected parsing error", err)
	} else if entry.Project != "project name" {
		t.Error("Entry has wrong project name")
	} else if !entry.TimeIn.Equal(parsedTimeIn) {
		t.Error("Entry has wrong timeIn")
	} else if !entry.TimeOut.Equal(parsedTimeOut) {
		t.Error("Entry has wrong timeOut")
	}
}

func TestParseLogLineBadLineInvalidTime(t *testing.T) {
	line := "project\t2015-02-10T15:30:10Z\tnot-a-time"
	if _, err := parseLogLine(line); err == nil {
		t.Error("Parsing should have failed")
	}
}

func TestParseLogLineBadLineTooManyElements(t *testing.T) {
	line := "IN\t2015-02-10T15:30:10Z\t2015-02-10T15:30:10Z\textra"
	if _, err := parseLogLine(line); err == nil {
		t.Error("Parsing should have failed")
	}
}

func TestParseLogLineBadLineTooFewElements(t *testing.T) {
	line := "FOO"
	if _, err := parseLogLine(line); err == nil {
		t.Error("Parsing should have failed")
	}
}
