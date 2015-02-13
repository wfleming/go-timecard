package timecard

import (
	"time"
)

type ProjectHours struct {
	Project string
	Hours   float64
}

type DaySummary struct {
	Date  time.Time
	Hours []ProjectHours // should be sorted by project name
}

type Summary struct {
	entries   []Entry
	Summaries []DaySummary // should be sorted by date
}

func NewSummary(entries []Entry) *Summary {
	return &Summary{entries, nil}
}

// builds the summary of the entries
//TODO: refactor, test
func (s *Summary) BuildSummary() {
	// build all the data in a structure of maps: easier while iterating over
	// a bunch of entries in unknown order. After summarizing, we'll convert
	// the map into []DaySummary, with appropriate sorting.
	// the map is date => project name => hours
	m := make(map[time.Time](map[string]float64))

	for _, entry := range s.entries {
		entryDate := atMidnight(entry.TimeIn)
		projHourMap, _ := m[entryDate]
		projHourMap[entry.Project] = projHourMap[entry.Project] + entry.Duration().Hours()
	}
}

// return a time with all non-date fields set to 0, i.e. represent only the
// day, not the time of day
func atMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
