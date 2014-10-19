package history_entries

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../environments"
)

type Parser func(env *environments.Environment, scanner *bufio.Scanner, filter *regexp.Regexp) ([]HistoryEntry, error)

var (
	bashTimestampRegexp = regexp.MustCompile(`^#'s*\d+$`)
	zshLineRegexp       = regexp.MustCompile(`^:\s*(\d*)\s*:(\d*)\s*;(.*)`)

	historyEventsCapacity = 5000
)

func init() {
	histFileSize := os.Getenv("HISTFILESIZE")
	if histFileSize != "" {
		if converted, err := strconv.Atoi(histFileSize); err == nil && converted > 0 {
			historyEventsCapacity = converted
		}
	}
}

func getParser(env *environments.Environment) Parser {
	shell, err := env.GetShell()
	if err != nil {
		panic(err)
	}

	if shell == environments.SHELL_ZSH {
		return parseZsh
	} else {
		return parseBash
	}
}

func parseBash(env *environments.Environment, scanner *bufio.Scanner, filter *regexp.Regexp) ([]HistoryEntry, error) {
	var currentNumber uint = 0
	currentTime := 0
	events := prepareHistoryEntries()
	currentEvent := HistoryEntry{}
	continueToConsume := false

	for scanner.Scan() {
		text := scanner.Text()

		if continueToConsume {
			currentEvent.command += "\n" + text
			if strings.HasSuffix(text, "\\") {
				continue
			}
			continueToConsume = false
			events = append(events, currentEvent)
			currentEvent = HistoryEntry{}
		}

		if bashTimestampRegexp.MatchString(text) {
			if converted, err := strconv.Atoi(text[1:]); err == nil {
				currentTime = converted
			}
		} else {
			if filter.MatchString(text) {
				currentEvent.command += text
				currentEvent.number = currentNumber
				currentEvent.timestamp = currentTime

				continueToConsume = strings.HasSuffix(text, "\\")
				if !continueToConsume {
					events = append(events, currentEvent)
					currentEvent = HistoryEntry{}
				}
			}
			currentNumber++
		}
	}

	return scanEnd(scanner, events)
}

func parseZsh(env *environments.Environment, scanner *bufio.Scanner, filter *regexp.Regexp) ([]HistoryEntry, error) {
	var currentNumber uint = 0
	events := prepareHistoryEntries()
	currentEvent := HistoryEntry{}
	continueToConsume := false
	logger, _ := env.GetLogger()

	for scanner.Scan() {
		text := scanner.Text()
		logger.Info("Got text line ", text)

		if continueToConsume {
			currentEvent.command += "\n" + text
			if strings.HasSuffix(text, "\\") {
				continue
			}
			continueToConsume = false
			events = append(events, currentEvent)
			currentEvent = HistoryEntry{}
		}

		matcher := zshLineRegexp.FindStringSubmatch(text)
		if matcher == nil {
			continue
		}

		timestamp := matcher[1]
		command := matcher[3]
		if timestamp == "" || command == "" {
			continue
		}
		currentNumber++

		if filter != nil && !filter.MatchString(command) {
			continue
		}

		converted, _ := strconv.Atoi(timestamp)
		currentEvent.command += command
		currentEvent.number = currentNumber
		currentEvent.timestamp = converted

		continueToConsume = strings.HasSuffix(text, "\\")
		if !continueToConsume {
			events = append(events, currentEvent)
			currentEvent = HistoryEntry{}
		}
	}

	return scanEnd(scanner, events)
}

func scanEnd(scanner *bufio.Scanner, events []HistoryEntry) ([]HistoryEntry, error) {
	if err := scanner.Err(); err != nil {
		return nil, err
	} else {
		return events, nil
	}
}

func prepareHistoryEntries() []HistoryEntry {
	return make([]HistoryEntry, 0, historyEventsCapacity)
}
