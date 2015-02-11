package punchcard

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Action int

const (
	IN  Action = iota
	OUT        = iota
)

type LogLine struct {
	action  Action
	time    time.Time
	project string
}

func actionGetName(action Action) (string, error) {
	if IN == action {
		return "IN", nil
	} else if OUT == action {
		return "OUT", nil
	} else {
		return "", errors.New("not a valid Action value")
	}
}

func actionFromName(name string) (Action, error) {
	if "IN" == name {
		return IN, nil
	} else if "OUT" == name {
		return OUT, nil
	} else {
		return -1, errors.New("not a valid action name")
	}
}

func parseLogLine(line string) (LogLine, error) {
	// parse the line into usable bits
	pieces := strings.Split(line, "\t")
	if len(pieces) < 2 {
		return LogLine{},
			fmt.Errorf("Not enough line elements in line \"%s\"",
				line)
	}
	action, err := actionFromName(pieces[0])
	if err != nil {
		return LogLine{}, err
	}
	if (action == IN && 3 != len(pieces)) || (action == OUT && 2 != len(pieces)) {
		return LogLine{},
			fmt.Errorf("Wrong number of line elements in \"%s\"",
				line)
	}
	time, err := time.Parse(time.RFC3339, pieces[1])
	if err != nil {
		return LogLine{},
			fmt.Errorf("Error parsing time in line \"%s\"", line)
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
