package punchcard

import (
	"errors"
	"time"
)

type Entry struct {
	timeIn  time.Time
	timeOut time.Time
	project string
}

func NewEntry() *Entry {
	return &Entry{}
}

func (entry *Entry) IsZero() bool {
	return (entry.timeIn.IsZero() && entry.timeOut.IsZero() &&
		"" == entry.project)
}

// modifies an entry with info in a LogLine (if valid)
func (entry *Entry) pushLogLine(line LogLine) error {
	switch {
	case !entry.timeIn.IsZero() && !entry.timeOut.IsZero():
		return errors.New("Entry full: no more actions allowed")
	case entry.timeIn.IsZero() && IN != line.action:
		return errors.New("Entry must punch in before punching out")
	case !entry.timeIn.IsZero() && OUT != line.action:
		return errors.New("Entry must cannot punch in twice")
	case "" == line.project:
		return errors.New("Empty project name not permitted")
	case OUT == line.action && line.time.Before(entry.timeIn):
		return errors.New("Punch out time cannot be before in time")
	case OUT == line.action && line.project != entry.project:
		return errors.New("Punch out project does not match in project")
	}

	switch line.action {
	case IN:
		entry.timeIn = line.time
		entry.project = line.project
	case OUT:
		entry.timeOut = line.time
	default:
		return errors.New("Invalid LogLine action")
	}

	return nil
}
