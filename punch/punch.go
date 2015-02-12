package main

import (
	"fmt"
	"github.com/wfleming/go-punchcard/punchcard"
	"os"
)

const defaultFileName = "~/.punch/entries.log"

type appConfig struct {
	log *punchcard.Log
}

type appCommand struct {
	run       func(*appConfig, []string)
	printHelp func()
}

var commands = map[string]appCommand{}

func main() {
	setupCommands()

	config, err := makeConfig()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if 0 == len(os.Args[1:]) {
		commands["help"].run(config, os.Args[1:])
		return
	}

	commandName := os.Args[1]
	command, exists := commands[commandName]

	if exists {
		command.run(config, os.Args[2:])
	} else {
		commands["help"].run(config, os.Args[2:])
	}
}

func makeConfig() (*appConfig, error) {
	return nil, nil
}

// can't do it as part of decl, or there's a reference loop
func setupCommands() {
	commands["help"] = appCommand{runHelp, printMainHelp}
	commands["in"] = appCommand{runIn, printInHelp}
	commands["out"] = appCommand{runOut, printOutHelp}
	commands["summary"] = appCommand{runSummary, printSummaryHelp}
}
