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
	summaries []DaySummary // should be sorted by date
}

func NewSummary(entries []Entry) *Summary {
	return &Summary{entries, nil}
}

func (s *Summary) GetSummaries() []DaySummary {
	if s.summaries == nil {
		s.summaries = s.buildSummaries()
	}
	return s.summaries
}

// builds the summary of the entries in s
func (s *Summary) buildSummaries() []DaySummary {
	// get map of date => project name => hours
	//m := s.buildDataMap()

	//TODO transform map to []DaySummary
	return nil
}

// build all the data in a structure of maps: easier while iterating over
// a bunch of entries in unknown order. After summarizing, we'll convert
// the map into []DaySummary, with appropriate sorting.
// the map is date => project name => hours
func (s *Summary) buildDataMap() map[time.Time](map[string]float64) {
	m := make(map[time.Time](map[string]float64))

	for _, entry := range s.entries {
		entryDate := atMidnight(entry.TimeIn)
		var projHourMap, exists = m[entryDate]
		if !exists {
			m[entryDate] = make(map[string]float64)
			projHourMap = m[entryDate]
		}
		hours := projHourMap[entry.Project] + entry.Duration().Hours()
		projHourMap[entry.Project] = hours
	}

	return m
}

// return a time with all non-date fields set to 0, i.e. represent only the
// day, not the time of day
func atMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
