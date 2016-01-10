// Functions for dealing with a punchcard-log formatted stream (file, pipe)
//
// A punchcard-log stream consists of lines with the following format:
// project name 	2006-01-02T15:04:05Z07:00	2006-01-02T16:30:05Z07:00
//
// The lines are tab-separated values: their meanings, in order, are:
// 1. The name of the project for the time entry
// 2. The time at wich work started (formatted according to RFC 3339)
// 3. The time at wich work stopped (formatted according to RFC 3339)
//
// The last element my be missing for the most recent entry in the log.
// That is, the last Entry may be punched in but not yet punched out.
package timecard

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// The Log type represents the current state of the set of Entries, as well as
// the Reader & Writer for reading in Entries & writing new ones.
type Log struct {
	in      *io.Reader
	out     *io.Writer
	entries []*Entry
}

// NewLog constructs a new Log from a Reader & a Writer
func NewLog(in io.Reader, out io.Writer) *Log {
	return &Log{&in, &out, make([]*Entry, 0)}
}

// Returns an array of all entries contained in the log.
// The entries will be read from the log's Reader if they have not already been read.
// The array of entries will be returned, or if an error occurred, an empty array will
// be returned along the with the error.
func (log *Log) AllEntries() ([]*Entry, error) {
	// entries are cached within the Log, and only parsed once.
	if len(log.entries) > 0 {
		return log.entries, nil
	}

	scanner := bufio.NewScanner(*log.in)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if 0 == len(line) {
			continue // ignore empty lines
		}
		entry, err := parseLogLine(line)
		if err != nil {
			return nil, err
		}
		log.entries = append(log.entries, &entry)
	}

	return log.entries, nil
}

// LastEntry returns the last entry from `AllEntries`, or nil if there are
// no entries, or the error getting `AllEntries`, if one occurred.
func (log *Log) LastEntry() (*Entry, error) {
	allEntries, err := log.AllEntries()
	if err != nil {
		return nil, err
	} else if 0 == len(allEntries) {
		return nil, nil
	}
	return allEntries[len(allEntries)-1], nil
}

// PunchIn adds a new Entry to the log for the given projectName,
// punched in at time time. The new Entry will also be written to
// the log's Writer.
func (log *Log) PunchIn(timeIn time.Time, projectName string) error {
	lastEntry, err := log.LastEntry()
	if err != nil {
		return err
	} else if lastEntry != nil && lastEntry.TimeOut.IsZero() {
		return errors.New("Last entry should have punched out")
	}
	entry := Entry{projectName, timeIn, time.Time{}}
	entryStr := entry.String()
	bytes, err := (*log.out).Write([]byte(entryStr))
	if err != nil {
		return err
	} else if bytes != len(entryStr) {
		return errors.New("Wrong number of bytes written")
	}
	log.entries = append(log.entries, &entry)
	return nil
}

// PunchOut punches out the last Entry in the log at the given time, and
// writes the rest of the Entry to the Log's Writer.
func (log *Log) PunchOut(timeOut time.Time) error {
	lastEntry, err := log.LastEntry()
	if err != nil {
		return err
	} else if lastEntry == nil || lastEntry.IsZero() {
		return errors.New("Entry should not be empty")
	} else if !lastEntry.TimeOut.IsZero() {
		return errors.New("Entry is already punched out")
	}
	entryStr := fmt.Sprintf("\t%s\n", timeOut.Format(time.RFC3339))
	bytes, err := (*log.out).Write([]byte(entryStr))
	if err != nil {
		return err
	} else if bytes != len(entryStr) {
		return errors.New("Wrong number of bytes written")
	}
	lastEntry.TimeOut = timeOut
	return nil
}
