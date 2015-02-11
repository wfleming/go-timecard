package punchcard

import (
	"bytes"
	"testing"
	"time"
)

func TestNewLog(t *testing.T) {
	data := bytes.NewBuffer(nil)
	log1 := NewLog(data, data)
	if *log1.in != data {
		t.Error("log1 in != expected in")
	}
	if *log1.out != data {
		t.Error("log1 out != expected in")
	}

	in, out := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	log2 := NewLog(in, out)
	if *log2.in != in {
		t.Error("log2 in != expected in")
	}
	if *log2.out != out {
		t.Error("log2 out != expected in")
	}
}

var time1, _ = time.Parse(time.RFC3339, "2015-02-10T15:30:10Z")
var time2, _ = time.Parse(time.RFC3339, "2015-02-10T16:30:10Z")
var time3, _ = time.Parse(time.RFC3339, "2015-02-10T16:45:10Z")
var lines = "IN\t2015-02-10T15:30:10Z\tcoding\n" +
	"OUT\t2015-02-10T16:30:10Z\tcoding\n" +
	"IN\t2015-02-10T16:45:10Z\ttyping\n"

func TestAllLogLines(t *testing.T) {
	data := bytes.NewBufferString(lines)
	log := NewLog(data, data)

	loglines, err := log.allLogLines()
	if err != nil {
		t.Fatal("getting LogLines should not have failed")
	}

	if len(loglines) != 3 {
		t.Fatal("wrong number of LogLines parsed", len(loglines))
	}

	var logline = loglines[0]
	if !logline.time.Equal(time1) || logline.project != "coding" ||
		logline.action != IN {
		t.Error("first LogLine not as expected", logline)
	}

	logline = loglines[1]
	if !logline.time.Equal(time2) || logline.project != "coding" ||
		logline.action != OUT {
		t.Error("second LogLine not as expected", logline)
	}

	logline = loglines[2]
	if !logline.time.Equal(time3) || logline.project != "typing" ||
		logline.action != IN {
		t.Error("third LogLine not as expected", logline)
	}
}

func TestAllLogLinesBad(t *testing.T) {
	data := bytes.NewBufferString("foo")
	log := NewLog(data, data)

	if _, err := log.allLogLines(); err == nil {
		t.Error("getting LogLines should have failed")
	}
}

func TestAllEntries(t *testing.T) {
	data := bytes.NewBufferString(lines)
	log := NewLog(data, data)

	entries, err := log.AllEntries()
	if err != nil {
		t.Fatal("turning LogLines into entries should not have failed")
	}

	if len(entries) != 2 {
		t.Fatal("wrong number of entries found", len(entries))
	}

	var entry = entries[0]
	if entry.timeIn != time1 || entry.timeOut != time2 ||
		entry.project != "coding" {
		t.Error("first entry not as expeced", entry)
	}

	entry = entries[1]
	if entry.timeIn != time3 || !entry.timeOut.IsZero() ||
		entry.project != "typing" {
		t.Error("second entry not as expected", entry)
	}
}

func TestLastEntryEmpty(t *testing.T) {
	data := bytes.NewBuffer(nil)
	log := NewLog(data, data)

	entry, err := log.LastEntry()
	if err != nil {
		t.Fatal("LastEntry should not have failed")
	}

	if entry != nil {
		t.Error("log is empty; entry should have been nil")
	}
}

func TestLastEntryNonEmpty(t *testing.T) {
	data := bytes.NewBuffer([]byte(lines))
	log := NewLog(data, data)

	entry, err := log.LastEntry()
	if err != nil {
		t.Fatal("turning LogLines into entries should not have failed")
	}

	if entry.timeIn != time3 || !entry.timeOut.IsZero() ||
		entry.project != "typing" {
		t.Error("last entry not as expected", entry)
	}
}

func TestPunchInAndOut(t *testing.T) {
	data := bytes.NewBufferString("")
	log := NewLog(data, data)

	if err := log.PunchIn(time1, "coding"); err != nil {
		t.Fatal("punch in should not have failed")
	}

	var expected = "IN\t2015-02-10T15:30:10Z\tcoding\n"
	if expected != data.String() {
		t.Errorf("after punch in, log buffer contents do not match: expected {%s} but got {%s}",
			expected, data.String())
	}

	if err := log.PunchOut(time2); err != nil {
		t.Fatal("punch in should not have failed")
	}

	expected = expected + "OUT\t2015-02-10T16:30:10Z\tcoding\n"
	if expected != data.String() {
		t.Errorf("after punch out, log buffer contents do not match: expected {%s} but got {%s}",
			expected, data.String())
	}
}

func TestPunchOutError(t *testing.T) {
	data := bytes.NewBufferString("")
	log := NewLog(data, data)

	err := log.PunchOut(time.Now())
	if err == nil {
		t.Error("punch out should have failed")
	}
}

func TestPunchOutExisting(t *testing.T) {
	data := bytes.NewBufferString(lines)
	log := NewLog(data, data)
	time := time.Now()

	err := log.PunchOut(time)
	if err != nil {
		t.Error("punch out should not have failed")
	}
}
