package punchcard

import (
	"testing"
	"time"
)

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

func TestParseLogLineGoodInLine(t *testing.T) {
	parsedTime, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	line := "IN\t2015-02-10T15:30:10Z\tcoding"
	logline, err := parseLogLine(line)
	if err != nil {
		t.Error("Unexpected parsing error", err)
	} else if logline.project != "coding" {
		t.Error("LogLine has wrong project name")
	} else if !logline.time.Equal(parsedTime) {
		t.Error("LogLine has wrong time")
	} else if IN != logline.action {
		t.Error("LogLine has wrong action")
	}
}

func TestParseLogLineGoodOutLine(t *testing.T) {
	parsedTime, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	line := "OUT\t2015-02-10T15:30:10Z\tproject name"
	logline, err := parseLogLine(line)
	if err != nil {
		t.Error("Unexpected parsing error", err)
	} else if logline.project != "project name" {
		t.Error("LogLine has wrong project name")
	} else if !logline.time.Equal(parsedTime) {
		t.Error("LogLine has wrong time")
	} else if OUT != logline.action {
		t.Error("LogLine has wrong action")
	}
}

func TestParseLogLineBadOutLine(t *testing.T) {
	line := "OUT\t2015-02-10T15:30:10Z"
	if _, err := parseLogLine(line); err == nil {
		t.Error("Parsing should have failed")
	}
}

func TestParseLogLineBadInLine1(t *testing.T) {
	line := "IN\t2015-02-10T15:30:10Z"
	if _, err := parseLogLine(line); err == nil {
		t.Error("Parsing should have failed")
	}
}

func TestParseLogLineBadInLine2(t *testing.T) {
	line := "FOO\t2015-02-10T15:30:10Z\tfoo bar"
	if _, err := parseLogLine(line); err == nil {
		t.Error("Parsing should have failed")
	}
}

func TestParseLogLineBadInLine3(t *testing.T) {
	line := "IN\tbad-time\tfoo bar"
	if _, err := parseLogLine(line); err == nil {
		t.Error("Parsing should have failed")
	}
}

func TestLogLineStringInLine(t *testing.T) {
	parsedTime, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	logline := LogLine{IN, parsedTime, "project name"}
	strline := logline.String()
	if "IN\t2015-02-10T15:30:10Z\tproject name" != strline {
		t.Error("Not expected line string", strline)
	}
}

func TestLogLineStringOutLine(t *testing.T) {
	parsedTime, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	logline := LogLine{OUT, parsedTime, "project name"}
	strline := logline.String()
	if "OUT\t2015-02-10T15:30:10Z\tproject name" != strline {
		t.Error("Not expected line string", strline)
	}
}

func TestLogLineStringBadLine(t *testing.T) {
	parsedTime, _ := time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
	logline := LogLine{-1, parsedTime, "project name"}
	strline := logline.String()
	if "INVALID_LOG_LINE_ACTION" != strline {
		t.Error("Not expected line string", strline)
	}
}
