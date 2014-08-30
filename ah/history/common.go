package history

// --- Imports

import (
	"bufio"
	"io"
	"time"
)

// --- Constants

const TimeNow = time.Now()

// --- Structs

type HistoryCommand struct {
	Command   string
	Timestamp time.Time
}

type historyScanner struct {
	scanner *bufio.Scanner
	content []HistoryCommand
}

// --- Interfaces

type HistoryScannerInterface interface {
	Init(io.Reader)
	GetCommands() ([]HistoryCommand, error)
}

// --- Methods

func (hs *historyScanner) Init(reader io.Reader, historySize int) {
	hs.scanner = bufio.NewScanner(reader)
	hs.content = make([]HistoryCommand, 0, historySize)
}

func (hs *historyScanner) addCommand(command string, timestamp time.Time) {
	commandToAdd := HistoryCommand{command, timestamp}
	hs.content = append(hs.content, commandToAdd)
}

func (hs *historyScanner) finishScan() ([]HistoryCommand, error) {
	if err := hs.scanner.Err(); err != nil {
		return nil, err
	}
	return hs.content, nil
}

// --- Functions

func convertTimestamp(timestamp int) time.Time {
	return time.Unix(int64(timestamp), 0)
}
