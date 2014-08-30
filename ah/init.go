package utils

import (
	"time"
)


const (
	DEFAULT_TERMINAL_WIDTH = 80
	DEFAULT_HISTORY_SIZE = 5000
)

var (
	TerminalWidth = DEFAULT_TERMINAL_WIDTH
	ShellEnv = Shell{}
	HistoryFilePath string
	HistoryDateFormat string
	TimeNow time.Time
)

func init() {
	ShellEnv.Discover()

	width, err := GetTerminalWidth()
	if err == nil {
		TerminalWidth = width
	}

	// HistoryFilePath = ShellEnv.GetEnv("HISTFILE")
	HistoryFilePath = "/home/nineseconds/.zsh_hhh"
	// HistoryDateFormat = os.Getenv("HISTTIMEFORMAT")
	TimeNow = time.Now()
}
