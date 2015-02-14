package timecard

import (
	"testing"
	"time"
)

func TestNewSummary(t *testing.T) {
	entries := []*Entry{&Entry{time.Now(), time.Now().Add(50000), "task"}}
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

func getTestEntries() (time.Time, []*Entry) {
	d1, _ := time.Parse(time.RFC3339, "2015-02-10T15:00:00Z")
	a := func(t time.Time, hours float64) time.Time {
		duration := time.Duration(hours * float64(time.Hour))
		return t.Add(duration)
	}
	entries := []*Entry{
		&Entry{d1, a(d1, 1), "task1"},
		&Entry{a(d1, 2), a(d1, 2.5), "task2"},
		&Entry{a(d1, 5), a(d1, 6.5), "task1"},
		&Entry{a(d1, 24), a(d1, 26), "task2"},
	}

	return d1, entries
}

func TestBuildDataMap(t *testing.T) {
	d1, entries := getTestEntries()
	s := NewSummary(entries)

	m := s.buildDataMap()

	midnight := atMidnight(d1)
	if h := m[midnight]["task1"]; h != 2.5 {
		t.Error("wrong hours for task1 on day 1", h)
	}
	if h := m[midnight]["task2"]; h != 0.5 {
		t.Error("wrong hours for task2 on day 1", h)
	}
	if h := m[midnight.Add(24*time.Hour)]["task2"]; h != 2 {
		t.Error("wrong hours for task2 on day 2", h)
	}
}

func TestBuildSummaries(t *testing.T) {
	d1, entries := getTestEntries()
	s := NewSummary(entries)

	summaries := s.GetSummaries()

	if len(summaries) != 2 {
		t.Fatal("unexpected number of day summaries", len(summaries))
	}

	if summaries[0].Date != atMidnight(d1) {
		t.Error("first entry has unexpected date", summaries[0].Date)
	}

	if n := len(summaries[0].Hours); n != 2 {
		t.Fatal("Wrong number of hour entries in day1", n)
	}

	if ph := summaries[0].Hours[0]; ph.Project != "task1" {
		t.Error("wrong 1st hours entry project in day1", ph)
	}

	if ph := summaries[0].Hours[0]; ph.Hours != 2.5 {
		t.Error("wrong 1st hours entry time in day1", ph)
	}

	if ph := summaries[0].Hours[1]; ph.Project != "task2" {
		t.Error("wrong 2nd hours entry project in day1", ph)
	}

	if ph := summaries[0].Hours[1]; ph.Hours != 0.5 {
		t.Error("wrong 2nd hours entry time in day1", ph)
	}

	if summaries[1].Date != atMidnight(d1).Add(24*time.Hour) {
		t.Error("first entry has unexpected date", summaries[0].Date)
	}

	if n := len(summaries[1].Hours); n != 1 {
		t.Fatal("Wrong number of hour entries in day2", n)
	}

	if ph := summaries[1].Hours[0]; ph.Project != "task2" {
		t.Error("wrong 1st hours entry project in day2", ph)
	}

	if ph := summaries[1].Hours[0]; ph.Hours != 2 {
		t.Error("wrong 1st hours entry project in day2", ph)
	}
}
