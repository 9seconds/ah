package history_entries

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	logrus "github.com/Sirupsen/logrus"

	"../environments"
)

type Parser func(env *environments.Environment, scanner *bufio.Scanner, filter *regexp.Regexp) ([]*HistoryEntry, error)

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

func parseBash(env *environments.Environment, scanner *bufio.Scanner, filter *regexp.Regexp) ([]*HistoryEntry, error) {
	var currentNumber uint = 0
	currentTime := 0
	events := prepareHistoryEntries()
	currentEvent := new(HistoryEntry)
	continueToConsume := false
	logger, _ := env.GetLogger()

	logger.Debug("Parse Bash")

	for scanner.Scan() {
		text := scanner.Text()
		logger.WithFields(logrus.Fields{
			"currentEvent":      currentEvent,
			"continueToConsume": continueToConsume,
			"currentNumber":     currentNumber,
			"text":              text,
		}).Info("Parse new line")

		if continueToConsume {
			currentEvent.command += "\n" + text
			if strings.HasSuffix(text, "\\") {
				continue
			}
			continueToConsume = false

			logger.WithFields(logrus.Fields{
				"currentEvent":      currentEvent,
				"continueToConsume": continueToConsume,
				"currentNumber":     currentNumber,
				"text":              text,
			}).Info("Stop consuming")

			events = append(events, currentEvent)
			currentEvent = new(HistoryEntry)
		}

		if bashTimestampRegexp.MatchString(text) {
			logger.WithFields(logrus.Fields{
				"currentEvent":      currentEvent,
				"continueToConsume": continueToConsume,
				"currentNumber":     currentNumber,
				"text":              text,
			}).Info("Matched timestamp")
			if converted, err := strconv.Atoi(text[1:]); err == nil {
				currentTime = converted
				logger.WithFields(logrus.Fields{
					"currentEvent":      currentEvent,
					"continueToConsume": continueToConsume,
					"currentNumber":     currentNumber,
					"text":              text,
					"currentTime":       converted,
				}).Info("Set new current time")
			} else {
				logger.WithFields(logrus.Fields{
					"currentEvent":      currentEvent,
					"continueToConsume": continueToConsume,
					"currentNumber":     currentNumber,
					"text":              text,
					"error":             err,
				}).Warn("Cannot convert timestamp")
			}
		} else {
			if filter.MatchString(text) {
				currentEvent.command += text
				currentEvent.number = currentNumber
				currentEvent.timestamp = currentTime

				continueToConsume = strings.HasSuffix(text, "\\")
				if !continueToConsume {
					events = append(events, currentEvent)
					currentEvent = new(HistoryEntry)
				}
			} else {
				logger.WithFields(logrus.Fields{
					"currentEvent":      currentEvent,
					"continueToConsume": continueToConsume,
					"currentNumber":     currentNumber,
					"text":              text,
					"filter":            filter,
				}).Warn("Ignored by filter")
			}
			currentNumber++
		}
	}

	logger.WithFields(logrus.Fields{
		"error": scanner.Err(),
	}).Info("Finish scan")

	return scanEnd(scanner, events)
}

func parseZsh(env *environments.Environment, scanner *bufio.Scanner, filter *regexp.Regexp) ([]*HistoryEntry, error) {
	var currentNumber uint = 0
	events := prepareHistoryEntries()
	currentEvent := new(HistoryEntry)
	continueToConsume := false
	logger, _ := env.GetLogger()

	logger.Debug("Parse ZSH")

	for scanner.Scan() {
		text := scanner.Text()
		logger.WithFields(logrus.Fields{
			"currentEvent":      currentEvent,
			"continueToConsume": continueToConsume,
			"currentNumber":     currentNumber,
			"text":              text,
		}).Info("Parse new line")

		if continueToConsume {
			currentEvent.command += "\n" + text
			if strings.HasSuffix(text, "\\") {
				continue
			}
			continueToConsume = false

			logger.WithFields(logrus.Fields{
				"currentEvent":      currentEvent,
				"continueToConsume": continueToConsume,
				"currentNumber":     currentNumber,
				"text":              text,
			}).Info("Stop consuming")

			events = append(events, currentEvent)
			currentEvent = new(HistoryEntry)
		}

		matcher := zshLineRegexp.FindStringSubmatch(text)
		if matcher == nil {
			logger.WithFields(logrus.Fields{
				"currentEvent":      currentEvent,
				"continueToConsume": continueToConsume,
				"currentNumber":     currentNumber,
				"text":              text,
				"regexp":            zshLineRegexp,
			}).Info("Text does not match regexp")
			continue
		}

		timestamp := matcher[1]
		command := matcher[3]
		if timestamp == "" || command == "" {
			logger.WithFields(logrus.Fields{
				"currentEvent":      currentEvent,
				"continueToConsume": continueToConsume,
				"currentNumber":     currentNumber,
				"text":              text,
				"regexp":            zshLineRegexp,
				"matcher":           matcher,
			}).Info("Timestamp and command are empty so skip.")
			continue
		}
		currentNumber++

		if filter != nil && !filter.MatchString(command) {
			logger.WithFields(logrus.Fields{
				"currentEvent":      currentEvent,
				"continueToConsume": continueToConsume,
				"currentNumber":     currentNumber,
				"text":              text,
				"filter":            filter,
			}).Info("Ignored by filter.")
			continue
		}

		converted, _ := strconv.Atoi(timestamp)
		currentEvent.command += command
		currentEvent.number = currentNumber
		currentEvent.timestamp = converted

		continueToConsume = strings.HasSuffix(text, "\\")
		if !continueToConsume {
			events = append(events, currentEvent)
			currentEvent = new(HistoryEntry)
		}
	}

	logger.WithFields(logrus.Fields{
		"error": scanner.Err(),
	}).Info("Finish scan")

	return scanEnd(scanner, events)
}

func scanEnd(scanner *bufio.Scanner, events []*HistoryEntry) ([]*HistoryEntry, error) {
	if err := scanner.Err(); err != nil {
		return nil, err
	} else {
		return events, nil
	}
}

func prepareHistoryEntries() []*HistoryEntry {
	return make([]*HistoryEntry, 0, historyEventsCapacity)
}
