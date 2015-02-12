package main

import "fmt"

func runOut(args []string) {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		printOutHelp()
		return
	}

	fmt.Println("TODO: do punching out")
}

func printOutHelp() {
	fmt.Println("usage: punch out")
}
