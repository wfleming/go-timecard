package timecard

import (
	"fmt"
	"strings"
	"time"
)

// Entry represents a single line in a timecard log: it is the time some work
// began on a project and also when the work was finished. If the entry was
// never marked complete, TimeOut will be the zero time.
type Entry struct {
	Project string
	TimeIn  time.Time
	TimeOut time.Time
}

// NewEntry returns a new, empty Entry
func NewEntry() *Entry {
	return &Entry{}
}

// An Entry is consider zero if each of it's 3 members is also zero.
func (entry *Entry) IsZero() bool {
	return (entry.TimeIn.IsZero() && entry.TimeOut.IsZero() &&
		"" == entry.Project)
}

// An Entry's Duration is the difference between its TimeOut & TimeIn.
func (entry *Entry) Duration() time.Duration {
	return entry.TimeOut.Sub(entry.TimeIn)
}

// An Entry's String representation, as it should appear in the log file.
func (entry Entry) String() string {
	if entry.TimeOut.IsZero() {
		return fmt.Sprintf("%s\t%s",
			entry.Project,
			entry.TimeIn.Format(time.RFC3339))
	}
	return fmt.Sprintf("%s\t%s\t%s",
		entry.Project,
		entry.TimeIn.Format(time.RFC3339),
		entry.TimeOut.Format(time.RFC3339))
}

// Parse an Entry from a string representation (inverse of String(), basically)
func parseLogLine(line string) (Entry, error) {
	// parse the line into usable bits
	pieces := strings.Split(line, "\t")
	if len(pieces) < 2 || len(pieces) > 3 {
		return Entry{},
			fmt.Errorf("Wrong number of elements in line \"%s\"",
				line)
	}
	timeIn, err := time.Parse(time.RFC3339, pieces[1])
	if err != nil {
		return Entry{},
			fmt.Errorf("Error parsing timeIn in line \"%s\"", line)
	}
	timeOut := time.Time{}
	if 3 == len(pieces) {
		timeOut, err = time.Parse(time.RFC3339, pieces[2])
		if err != nil {
			return Entry{},
				fmt.Errorf("Error parsing timeIn in line \"%s\"", line)
		}
	}

	var entry Entry
	entry.Project = pieces[0]
	entry.TimeIn = timeIn
	entry.TimeOut = timeOut

	return entry, nil
}
