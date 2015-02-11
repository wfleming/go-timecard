package punchcard

import (
	"time"
)

type Entry struct {
	timeIn  time.Time
	timeOut time.Time
	project string
}

func (entry Entry) IsZero() bool {
	return (entry.timeIn.IsZero() && entry.timeOut.IsZero() &&
		"" == entry.project)
}
