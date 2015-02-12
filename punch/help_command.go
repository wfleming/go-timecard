package main

import "fmt"

func runHelp(config *appConfig, args []string) {
	if 0 == len(args) {
		printMainHelp()
	} else if c, exists := commands[args[0]]; !exists {
		printMainHelp()
	} else {
		c.printHelp()
	}
}

func printMainHelp() {
	fmt.Println("Punch is a simple time tracker.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("\tpunch <command> [arguments]")
	fmt.Println("")
	fmt.Println("The commands are:")
	fmt.Println("")
	fmt.Println("\tin\tpunch in to a project")
	fmt.Println("\tout\tpunch out of a project")
	fmt.Println("\tsummary\tsummarize your time")
	fmt.Println("")
	fmt.Println("Use \"punch help [command]\" for more information about a command.")
	fmt.Println("")
}
