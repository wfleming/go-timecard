package main

import (
	"fmt"
	"github.com/wfleming/go-punchcard/timecard"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const defaultFilename = "~/.punch/entries.log"

type appConfig struct {
	log     *punchcard.Log
	logfile *os.File // gets set if we're using an actual file for the log
}

type appCommand struct {
	run       func(*appConfig, []string)
	printHelp func()
}

var commands = map[string]appCommand{}

func main() {
	setupCommands()
	config, err := makeConfig()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer config.close()

	if 0 == len(os.Args[1:]) {
		commands["help"].run(config, os.Args[1:])
		return
	}

	commandName := os.Args[1]
	command, exists := commands[commandName]

	if exists {
		command.run(config, os.Args[2:])
	} else {
		commands["help"].run(config, os.Args[2:])
	}
}

func makeConfig() (*appConfig, error) {
	var config appConfig

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// this indicates data is being piped, so use stdin in for log
		config.log = punchcard.NewLog(os.Stdin, os.Stdout)
	} else {
		fh, err := getLogFile(defaultFilename)
		if err != nil {
			return nil, err
		}
		config.log = punchcard.NewLog(fh, fh)
		config.logfile = fh
	}

	return &config, nil
}

func getLogFile(filename string) (*os.File, error) {
	filename, err := sanitizeLogFileName(filename)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// file does not exist: dir may not exist either
		// attempt to create dir if it does not exist
		// file will be created by OpenFile later
		if err := os.Mkdir(filepath.Dir(filename), 0755); err != nil {
			return nil, err
		}
	}

	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return nil, err
	}

	return fh, nil
}

func sanitizeLogFileName(filename string) (string, error) {
	var err error
	// must do ~ replacement ourselves
	if filename[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		homedir := usr.HomeDir
		filename = strings.Replace(filename, "~", homedir, 1)
	}

	filename, err = filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// can't do it as part of decl, or there's a reference loop
func setupCommands() {
	commands["help"] = appCommand{runHelp, printMainHelp}
	commands["in"] = appCommand{runIn, printInHelp}
	commands["out"] = appCommand{runOut, printOutHelp}
	commands["summary"] = appCommand{runSummary, printSummaryHelp}
}

// cleanup necessary parts of the appConfig, such as closing file handles
func (config *appConfig) close() {
	if config.logfile != nil {
		config.logfile.Sync()
		config.logfile.Close()
	}
}
