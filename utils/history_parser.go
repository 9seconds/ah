package utils


import (
	"io"
	"bufio"
	"time"
)

type HistoryCommand struct {
	Command string
	Timestamp time.Time
}

type HistoryScannerInterface interface {
	Init(io.Reader)
	GetCommands() ([]HistoryCommand, error)
}

type historyScanner struct {
	scanner *bufio.Scanner
	content []HistoryCommand
}


func (hs *historyScanner) Init(reader io.Reader) {
	hs.scanner = bufio.NewScanner(reader)
	hs.content = make([]HistoryCommand, 0, DEFAULT_HISTORY_SIZE)
}
