package punchcard

import (
	"testing"
	"time"
)

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

func TestEntryPushLogLineErrors(t *testing.T) {
	now := time.Now()
	var entry Entry
	if entry.pushLogLine(LogLine{OUT, now, ""}) == nil {
		t.Error("pushing out to empty entry should error")
	}

	if entry.pushLogLine(LogLine{IN, now.Add(-500), ""}) == nil {
		t.Error("pushing empty project should error")
	}

	entry.project = "foo"
	entry.timeIn = now.Add(-500)

	if entry.pushLogLine(LogLine{IN, now, "foo"}) == nil {
		t.Error("punch in twice should error")
	}

	if entry.pushLogLine(LogLine{-1, now, "foo"}) == nil {
		t.Error("invalid action should error")
	}

	if entry.pushLogLine(LogLine{OUT, now.Add(-600), "foo"}) == nil {
		t.Error("can't punch out before punching in")
	}

	if entry.pushLogLine(LogLine{OUT, now.Add(600), "bar"}) == nil {
		t.Error("can't punch out with different project name")
	}

	entry.timeOut = now

	if entry.pushLogLine(LogLine{OUT, now.Add(600), ""}) == nil {
		t.Error("can't punch out twice")
	}
}

func TestEntryPushLogLineSuccess(t *testing.T) {
	project, now := "foo proj", time.Now()
	var entry Entry

	if entry.pushLogLine(LogLine{IN, now, project}) != nil {
		t.Error("punch in should not have failed")
	}

	if entry.project != project {
		t.Error("punch in should have set project")
	}

	if !entry.timeIn.Equal(now) {
		t.Error("punch in should have set timeIn")
	}

	future := now.Add(500)
	if entry.pushLogLine(LogLine{OUT, future, project}) != nil {
		t.Error("punch out should not have failed")
	}

	if entry.project != project {
		t.Error("punch out should not change project")
	}

	if !entry.timeOut.Equal(future) {
		t.Error("punch out should have set timeOut")
	}

}
