package main

import (
	// "github.com/wfleming/go-punchcard/punchcard"
	"os"
)

type command struct {
	run       func(args []string)
	printHelp func()
}

var commands = map[string]command{}

func main() {
	setupCommands()

	if 0 == len(os.Args[1:]) {
		commands["help"].run(os.Args[1:])
		return
	}

	commandName := os.Args[1]
	command, exists := commands[commandName]

	if exists {
		command.run(os.Args[2:])
	} else {
		commands["help"].run(os.Args[2:])
	}
}

// can't do it as part of decl, or there's a reference loop
func setupCommands() {
	commands["help"] = command{runHelp, printMainHelp}
	commands["in"] = command{runIn, printInHelp}
	commands["out"] = command{runOut, printOutHelp}
	commands["summary"] = command{runSummary, printSummaryHelp}
}
