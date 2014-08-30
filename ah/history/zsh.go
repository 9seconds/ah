package history

// --- Imports

import (
	"strconv"
	"time"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

// --- Vars

var historyLineRegex = pcre.MustCompile(`^:\s*(\d*)\s*:(\d*)\s*;(.*)`, 0)

// --- Structs

type historyScannerZsh struct {
	historyScanner
}

// --- Methods

func (hsz *historyScannerZsh) GetCommands() ([]HistoryCommand, error) {
	for hsz.scanner.Scan() {
		matcher := historyLineRegex.MatcherString(hsz.scanner.Text(), 0)

		timestamp := matcher.GroupString(1)
		command := matcher.GroupString(3)
		if timestamp == "" || command == "" {
			continue
		}

		commandTime := TimeNow
		if parsed, err := strconv.Atoi(timestamp); err == nil {
			commandTime = convertTimestamp(parsed)
		}

		hsz.addCommand(command, commandTime)
	}

	return hsz.finishScan()
}
