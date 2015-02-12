package main

import "fmt"

func runSummary(config *appConfig, args []string) {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		printSummaryHelp()
		return
	}

	fmt.Println("TODO: do summarizing")
}

func printSummaryHelp() {
	fmt.Println("usage: punch summary")
}
