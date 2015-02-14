package timecard

import (
	"sort"
	"time"
)

type ProjectHours struct {
	Project string
	Hours   float64
}

type byProject []ProjectHours

func (bp byProject) Len() int           { return len(bp) }
func (bp byProject) Swap(i, j int)      { bp[i], bp[j] = bp[j], bp[i] }
func (bp byProject) Less(i, j int) bool { return bp[i].Project < bp[j].Project }

type DaySummary struct {
	Date  time.Time
	Hours []ProjectHours // should be sorted by project name
}

type byDate []DaySummary

func (bd byDate) Len() int           { return len(bd) }
func (bd byDate) Swap(i, j int)      { bd[i], bd[j] = bd[j], bd[i] }
func (bd byDate) Less(i, j int) bool { return bd[i].Date.Before(bd[j].Date) }

type Summary struct {
	entries   []*Entry
	summaries []DaySummary // should be sorted by date
}

func NewSummary(entries []*Entry) *Summary {
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
	m := s.buildDataMap()
	var summaries = make([]DaySummary, 0)

	// convert that map to []DaySummary
	for date, projHourMap := range m {
		daySummary := DaySummary{date, make([]ProjectHours, 0)}

		for project, hours := range projHourMap {
			ph := ProjectHours{project, hours}
			daySummary.Hours = append(daySummary.Hours, ph)
		}

		sort.Sort(byProject(daySummary.Hours))

		summaries = append(summaries, daySummary)
	}

	sort.Sort(byDate(summaries))

	return summaries
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
