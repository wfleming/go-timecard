package main

import (
	"fmt"
	"github.com/wfleming/go-punchcard/timecard"
	"os"
)

const shortDate = "Monday, 02 January 2006"

func runSummary(config *appConfig, args []string) {
	if len(args) > 0 {
		printSummaryHelp()
		return
	}

	var entries, err = config.log.AllEntries()

	if err != nil {
		fmt.Printf("Error parsing entries: %s\n", err)
		os.Exit(1)
	}

	if len(entries) > 0 && entries[len(entries)-1].TimeOut.IsZero() {
		// exclude final entry if it has not been checked out
		entries = entries[:len(entries)]
	}

	exitIfNoEntries(entries)

	summary := timecard.NewSummary(entries)
	daySummaries := summary.GetSummaries()

	printSummariesStd(daySummaries)
}

func exitIfNoEntries(entries []*timecard.Entry) {
	if len(entries) == 0 {
		fmt.Println("No entries to summarize.")
		os.Exit(0)
	}
}

// do the standard printing of the summary data
func printSummariesStd(summaries []timecard.DaySummary) {
	for _, daySummary := range summaries {
		fmt.Println(daySummary.Date.Format(shortDate))

		for _, projHours := range daySummary.Hours {
			fmt.Printf("\t%s: %.2f\n", projHours.Project, projHours.Hours)
		}
		fmt.Println("")
	}
}

func printSummaryHelp() {
	fmt.Println("usage: punch summary")
}
