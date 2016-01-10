package timecard

import (
	"bytes"
	"io/ioutil"
	"os"
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
var lines = "coding\t2015-02-10T15:30:10Z\t2015-02-10T16:30:10Z\n" +
	"typing\t2015-02-10T16:45:10Z"

func TestAllEntriesGood(t *testing.T) {
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
	if entry.TimeIn != time1 || entry.TimeOut != time2 ||
		entry.Project != "coding" {
		t.Error("first entry not as expeced", entry)
	}

	entry = entries[1]
	if entry.TimeIn != time3 || !entry.TimeOut.IsZero() ||
		entry.Project != "typing" {
		t.Error("second entry not as expected", entry)
	}
}

func TestAllEntriesBad(t *testing.T) {
	data := bytes.NewBufferString("foo")
	log := NewLog(data, data)

	if _, err := log.AllEntries(); err == nil {
		t.Error("getting AllEntries should have failed")
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

	if entry.TimeIn != time3 || !entry.TimeOut.IsZero() ||
		entry.Project != "typing" {
		t.Error("last entry not as expected", entry)
	}
}

func TestPunchInAndOut(t *testing.T) {
	var fh, err = ioutil.TempFile(os.TempDir(), "log-test")
	defer fh.Close()

	if err != nil {
		t.Fatal("failed to make temp file", err)
	}

	log := NewLog(fh, fh)

	if err := log.PunchIn(time1, "coding"); err != nil {
		t.Fatal("punch in should not have failed", err)
	}

	var expected = "coding\t2015-02-10T15:30:10Z"
	fh.Seek(0, 0)
	var filestr, _ = ioutil.ReadAll(fh)
	if expected != string(filestr) {
		t.Errorf("after punch in, log buffer contents do not match: expected {%s} but got {%s}",
			expected, string(filestr))
	}

	if err := log.PunchOut(time2); err != nil {
		t.Fatal("punch out should not have failed", err)
	}

	expected = expected + "\t2015-02-10T16:30:10Z\n"
	fh.Seek(0, 0)
	filestr, _ = ioutil.ReadAll(fh)
	if expected != string(filestr) {
		t.Errorf("after punch out, log buffer contents do not match: expected {%s} but got {%s}",
			expected, string(filestr))
	}

	os.Remove(fh.Name())
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
