package main

import "fmt"

func runIn(config *appConfig, args []string) {
	if len(args) < 1 || args[0] == "--help" || args[0] == "-h" {
		printInHelp()
		return
	}

	fmt.Println("TODO: do punching in")
}

func printInHelp() {
	fmt.Println("usage: punch in <project name>")
}
