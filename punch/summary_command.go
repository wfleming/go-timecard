package main

import (
	"fmt"
	"github.com/wfleming/go-timecard/timecard"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

const (
	csvDate   = "2006-01-02"
	humanDate = "Monday, 02 January 2006"
)

func runSummary(config *appConfig, args []string) {
	if len(args) > 0 && args[0] != "--csv" {
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
		entries = entries[:(len(entries) - 1)]
	}

	exitIfNoEntries(entries)

	summary := timecard.NewSummary(entries)
	daySummaries := summary.GetSummaries()

	if len(args) > 0 && args[0] == "--csv" {
		printSummariesCsv(daySummaries)
	} else {
		printSummariesStd(daySummaries)
	}
}

func exitIfNoEntries(entries []*timecard.Entry) {
	if len(entries) == 0 {
		fmt.Println("No entries to summarize.")
		os.Exit(0)
	}
}

// do the standard printing of the summary data
func printSummariesStd(summaries []timecard.DaySummary) {
	w := tabwriter.NewWriter(os.Stdout, 1, 8, 0, '\t', tabwriter.TabIndent)

	for _, daySummary := range summaries {
		fmt.Println(daySummary.Date.Format(humanDate))
		var dayTotal = 0.0

		for _, projHours := range daySummary.Hours {
			line := fmt.Sprintf("\t%s:\t%.2f\n", projHours.Project, projHours.Hours)
			w.Write([]byte(line))
			dayTotal += projHours.Hours
		}
		totLine := fmt.Sprintf("\tDAY TOTAL:\t%.2f\n", dayTotal)
		w.Write([]byte(totLine))
		w.Flush()
		fmt.Println("")
	}
}

// printing the summary data as CSV
func printSummariesCsv(summaries []timecard.DaySummary) {
	var allProjects = make([]string, 0)

	indexStr := func(arr []string, needle string) int {
		for idx, str := range arr {
			if str == needle {
				return idx
			}
		}
		return -1
	}

	hoursForProjDate := func(project string, date time.Time) float64 {
		for _, daySummary := range summaries {
			if !daySummary.Date.Equal(date) {
				continue
			}
			for _, projHours := range daySummary.Hours {
				if projHours.Project == project {
					return projHours.Hours
				}
			}
		}
		return 0
	}

	// iterate once to get all project names
	for _, daySummary := range summaries {
		for _, projHours := range daySummary.Hours {
			if indexStr(allProjects, projHours.Project) < 0 {
				allProjects = append(allProjects, projHours.Project)
			}
		}
	}

	sort.Strings(allProjects)

	//////// print the CSV data ///////
	// first line is just dates
	fmt.Print("Project,")
	for idx, daySummary := range summaries {
		var comma = ""
		if idx < (len(summaries) - 1) {
			comma = ","
		}
		fmt.Printf("%s%s",
			daySummary.Date.Format(csvDate), comma)
	}
	fmt.Print("\n")
	// now print a row for each project
	for _, project := range allProjects {
		fmt.Printf("%s,", project)
		for idx, daySummary := range summaries {
			var comma = ""
			if idx < (len(summaries) - 1) {
				comma = ","
			}
			hours := hoursForProjDate(project, daySummary.Date)
			fmt.Printf("%.2f%s", hours, comma)
		}
		fmt.Print("\n")
	}
}

func printSummaryHelp() {
	fmt.Println("usage: punch summary [--csv]")
}
