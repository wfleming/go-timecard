package main

import (
	"fmt"
	"os"
	"time"
)

func runOut(config *appConfig, args []string) {
	if len(args) > 0 {
		printOutHelp()
		return
	}

	if err := config.log.PunchOut(time.Now()); err != nil {
		fmt.Printf("Error punching out: %s\n", err)
		os.Exit(1)
	}
	lastEntry, err := config.log.LastEntry()
	if err != nil {
		fmt.Printf(
			"Punched out, but encountered error getting "+
				"entry: %s\n",
			err)
		os.Exit(1)
	}
	fmt.Printf("Punched out of %s\n", lastEntry.Project)
}

func printOutHelp() {
	fmt.Println("usage: punch out")
}
