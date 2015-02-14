package main

import "fmt"

func runSummary(config *appConfig, args []string) {
	if len(args) > 0 {
		printSummaryHelp()
		return
	}

}

func printSummaryHelp() {
	fmt.Println("usage: punch summary")
}
