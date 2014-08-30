package utils

import (
	"time"
	"strconv"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

var historyLineRegex = pcre.MustCompile(`^:\s*(\d*)\s*:(\d*)\s*;(.*)`, 0)

type historyScannerZsh struct {
	historyScanner
}


func (hsb *historyScannerZsh) GetCommands() ([]HistoryCommand, error) {
	for hsb.scanner.Scan() {
		matcher := historyLineRegex.MatcherString(hsb.scanner.Text(), 0)

		timestamp := matcher.GroupString(1)
		command := matcher.GroupString(3)
		if timestamp == "" || command == "" {
			continue
		}

		commandTime := TimeNow
		if parsed, err := strconv.Atoi(timestamp); err == nil {
			commandTime = time.Unix(int64(parsed), 0)
		}

		parsedCommand := HistoryCommand{command, commandTime}
		hsb.content = append(hsb.content, parsedCommand)
	}

	if err := hsb.scanner.Err(); err != nil {
		return nil, err
	}

	return hsb.content, nil
}
