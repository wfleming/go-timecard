package main

import (
	"fmt"
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

	switch subcommand := os.Args[1]; subcommand {
	case "in":
		fmt.Println("TODO punch in")
	case "out":
		fmt.Println("TODO to punch out")
	case "summary":
		fmt.Println("TODO: to print summary")
	default:
		commands["help"].run(os.Args[2:])
	}
}

// can't do it as part of variable, or there's a reference loop
func setupCommands() {
	commands["help"] = command{runHelp, printMainHelp}
}
