package punchcard

import "errors"
import "time"

type Action int

const (
	IN  Action = iota
	OUT        = iota
)

type Entry struct {
	timeIn  time.Time
	timeOut time.Time
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
