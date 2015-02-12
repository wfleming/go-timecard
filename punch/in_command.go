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

	project := strings.Join(args, " ")
	if err := config.log.PunchIn(time.Now(), project); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Punched in to %s\n", project)
}

func printInHelp() {
	fmt.Println("usage: punch in <project name>")
}
