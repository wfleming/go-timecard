package punchcard

/* Functions for dealing with a punchcard-log formatted stream (file, pipe)
 *
 * A punchcard-log stream consists of lines with the following format:
 * IN	2006-01-02T15:04:05Z07:00	project name
 * OUT	2006-01-02T15:04:05Z07:00	project name
 *
 * An IN line *MUST* be followed by an OUT line. The project name field of an
 * OUT line *MUST* match the project name on the IN line before it.
 */

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"time"
)

type Log struct {
	in      *io.Reader
	out     *io.Writer
	entries []*Entry
}

func NewLog(in io.Reader, out io.Writer) *Log {
	return &Log{&in, &out, make([]*Entry, 0)}
}

func (log *Log) allLogLines() ([]LogLine, error) {
	var lines []LogLine

	scanner := bufio.NewScanner(*log.in)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if 0 == len(line) {
			continue // ignore empty lines
		}
		logline, err := parseLogLine(line)
		if err != nil {
			return nil, err
		}
		lines = append(lines, logline)
	}

	return lines, nil
}

func (log *Log) AllEntries() ([]*Entry, error) {
	// entries are cached within the Log, and only parsed once.
	if len(log.entries) > 0 {
		return log.entries, nil
	}

	loglines, err := log.allLogLines()

	if err != nil {
		return nil, err
	}

	for _, logline := range loglines {
		var entry *Entry
		// if last entry has IN but not OUT, next line should be OUT
		if len(log.entries) > 0 && log.entries[len(log.entries)-1].timeOut.IsZero() {
			entry = log.entries[len(log.entries)-1]
		} else {
			entry = NewEntry()
		}
		if err := entry.pushLogLine(logline); err != nil {
			return nil, err
		}
		if len(log.entries) == 0 || entry != log.entries[len(log.entries)-1] {
			log.entries = append(log.entries, entry)
		}
	}

	return log.entries, nil
}

func (log *Log) LastEntry() (*Entry, error) {
	allEntries, err := log.AllEntries()
	if err != nil {
		return nil, err
	} else if 0 == len(allEntries) {
		return nil, nil
	}
	return allEntries[len(allEntries)-1], nil
}

// write an IN line to the log (if valid)
func (log *Log) PunchIn(time time.Time, projectName string) error {
	lastEntry, err := log.LastEntry()
	if err != nil {
		return err
	} else if lastEntry != nil && lastEntry.timeOut.IsZero() {
		return errors.New("last entry should have punched out")
	}
	logline := LogLine{IN, time, projectName}
	strline := logline.String() + "\n"
	bytes, err := (*log.out).Write([]byte(strline))
	if err != nil {
		return err
	} else if bytes != len(strline) {
		return errors.New("Wrong number of bytes written")
	}
	entry := NewEntry()
	if err := entry.pushLogLine(logline); err != nil {
		return err
	}
	log.entries = append(log.entries, entry)
	return nil
}

// write an OUT line to the log (if valid)
func (log *Log) PunchOut(time time.Time) error {
	lastEntry, err := log.LastEntry()
	if err != nil {
		return err
	} else if lastEntry == nil || lastEntry.IsZero() {
		return errors.New("Entry should not be empty")
	}
	logline := LogLine{OUT, time, lastEntry.project}
	strline := logline.String() + "\n"
	bytes, err := (*log.out).Write([]byte(strline))
	if err != nil {
		return err
	} else if bytes != len(strline) {
		return errors.New("Wrong number of bytes written")
	} else if err := lastEntry.pushLogLine(logline); err != nil {
		return err
	}
	return nil
}
