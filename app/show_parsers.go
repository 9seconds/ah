package app

import (
	"bufio"
	"regexp"
	"strconv"
	"time"
)

var (
	bashTimestampRegexp = regexp.MustCompile(`^#'s*\d+$`)
	zshLineRegexp       = regexp.MustCompile(`^:\s*(\d*)\s*:(\d*)\s*;(.*)`)

	defaultEventsCapacity = 5000
)

func parseBash(scanner *bufio.Scanner, filter *regexp.Regexp) ([]HistoryEntry, error) {
	currentTime := time.Now()
	currentNumber := 1
	events := make([]HistoryEntry, 0, defaultEventsCapacity)

	for scanner.Scan() {
		text := scanner.Text()
		if bashTimestampRegexp.MatchString(text) {
			if converted, err := strconv.Atoi(text[1:]); err == nil {
				currentTime = convertTimestamp(converted)
			}
		} else if filter.MatchString(text) {
			newEvent := HistoryEntry{
				Number:     currentNumber,
				Command:    text,
				Timestamp:  currentTime,
				HasHistory: false} // HasHistory = false temporarily
			events = append(events, newEvent)
			currentNumber += 1
		}
	}

	return scanEnd(scanner, events)
}

func parseZsh(scanner *bufio.Scanner, filter *regexp.Regexp) ([]HistoryEntry, error) {
	currentNumber := 1
	events := make([]HistoryEntry, 0, defaultEventsCapacity)

	for scanner.Scan() {
		matcher := zshLineRegexp.FindStringSubmatch(scanner.Text())
		if matcher == nil {
			continue
		}
		timestamp := matcher[1]
		command := matcher[3]

		if timestamp == "" || command == "" {
			continue
		}
		if filter != nil && filter.MatchString(command) {
			continue
		}

		converted, _ := strconv.Atoi(timestamp)
		newEvent := HistoryEntry{
			Number:     currentNumber,
			Command:    command,
			Timestamp:  convertTimestamp(converted),
			HasHistory: false} // HasHistory = false temporarily
		events = append(events, newEvent)
		currentNumber += 1
	}

	return scanEnd(scanner, events)
}

func convertTimestamp(timestamp int) time.Time {
	return time.Unix(int64(timestamp), 0)
}

func scanEnd(scanner *bufio.Scanner, events []HistoryEntry) ([]HistoryEntry, error) {
	if err := scanner.Err(); err != nil {
		return nil, err
	} else {
		return events, nil
	}
}
