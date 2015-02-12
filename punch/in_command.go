package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func runIn(config *appConfig, args []string) {
	if len(args) < 1 || args[0] == "--help" || args[0] == "-h" {
		printInHelp()
		return
	}

	// punch out automatically if we need to
	lastEntry, err := config.log.LastEntry()
	if err != nil {
		fmt.Printf("Error fetching last entry: %s\n", err)
		os.Exit(1)
	} else if lastEntry != nil && lastEntry.TimeOut.IsZero() {
		fmt.Printf(
			"You are still punched in to %s. "+
				"Automatically punching you out.\n",
			lastEntry.Project)
		if err := config.log.PunchOut(time.Now()); err != nil {
			fmt.Printf("Error punching out: %s\n", err)
		}
	}

	// now do the punch in
	project := strings.Join(args, " ")
	if err := config.log.PunchIn(time.Now(), project); err != nil {
		fmt.Printf("Error punching in: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Punched in to %s\n", project)
}

func printInHelp() {
	fmt.Println("usage: punch in <project name>")
}
