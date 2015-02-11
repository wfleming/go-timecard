package punchcard

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type LogLine struct {
	action  Action
	time    time.Time
	project string
}

func parseLogLine(line string) (LogLine, error) {
	// parse the line into usable bits
	pieces := strings.Split(line, "\t")
	if len(pieces) < 2 {
		return LogLine{}, errors.New("Not enough line elements")
	}
	action, err := actionFromName(pieces[0])
	if err != nil {
		return LogLine{}, err
	}
	if (action == IN && 3 != len(pieces)) || (action == OUT && 2 != len(pieces)) {
		return LogLine{}, errors.New("Wrong number of line elements")
	}
	time, err := time.Parse(time.RFC3339, pieces[1])
	if err != nil {
		return LogLine{}, err
	}

	var logline LogLine
	logline.action = action
	logline.time = time
	if IN == action {
		logline.project = pieces[2]
	}

	return logline, nil
}

func (logline LogLine) String() string {
	actionName, _ := actionGetName(logline.action)
	switch logline.action {
	case IN:
		return fmt.Sprintf("%s\t%s\t%s",
			actionName,
			logline.time.Format(time.RFC3339),
			logline.project)
	case OUT:
		return fmt.Sprintf("%s\t%s",
			actionName,
			logline.time.Format(time.RFC3339))
	default:
		return "INVALID_LOG_LINE_ACTION"
	}
}
