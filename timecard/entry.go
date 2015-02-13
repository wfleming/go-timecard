package timecard

import (
	"errors"
	"time"
)

type Entry struct {
	TimeIn  time.Time
	TimeOut time.Time
	Project string
}

func NewEntry() *Entry {
	return &Entry{}
}

func (entry *Entry) IsZero() bool {
	return (entry.TimeIn.IsZero() && entry.TimeOut.IsZero() &&
		"" == entry.Project)
}

// modifies an entry with info in a LogLine (if valid)
func (entry *Entry) pushLogLine(line LogLine) error {
	switch {
	case !entry.TimeIn.IsZero() && !entry.TimeOut.IsZero():
		return errors.New("Entry full: no more actions allowed")
	case entry.TimeIn.IsZero() && IN != line.action:
		return errors.New("Entry must punch in before punching out")
	case !entry.TimeIn.IsZero() && IN == line.action:
		return errors.New("Entry cannot punch in twice")
	case "" == line.project:
		return errors.New("Empty project name not permitted")
	case OUT == line.action && line.time.Before(entry.TimeIn):
		return errors.New("Punch out time cannot be before in time")
	case OUT == line.action && line.project != entry.Project:
		return errors.New("Punch out project does not match in project")
	}

	switch line.action {
	case IN:
		entry.TimeIn = line.time
		entry.Project = line.project
	case OUT:
		entry.TimeOut = line.time
	default:
		return errors.New("Invalid LogLine action")
	}

	return nil
}
