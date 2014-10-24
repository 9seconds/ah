package historyentries

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

const (
	markHasHistory    = '*'
	markHasNoHistory = ' '
)

// HistoryEntry stores a command with its context.
type HistoryEntry struct {
	number     uint
	command    string
	timestamp  int
	hasHistory bool
}

// GetNumber returns a history number (may be executed with ! later).
func (he HistoryEntry) GetNumber() (uint, error) {
	if he.number == 0 {
		return 0, errors.New("Number is not set yet")
	}
	return he.number, nil
}

// GetCommand returns a command line which was executed.
func (he HistoryEntry) GetCommand() (string, error) {
	if he.command == "" {
		return "", errors.New("Command is not set yet")
	}
	return he.command, nil
}

// GetTimestamp returns a timestamp of the history entry.
func (he HistoryEntry) GetTimestamp() (int, error) {
	if he.timestamp == 0 {
		return 0, errors.New("Timestamp is not set yet")
	}
	return he.timestamp, nil
}

// GetTime returns a time structure of the history entry.
func (he HistoryEntry) GetTime() (*time.Time, error) {
	timestamp, err := he.GetTimestamp()
	if err != nil {
		return nil, err
	}
	return utils.ConvertTimestamp(timestamp), nil
}

// GetFormattedTime returns formatted time stamp of the history entry.
func (he HistoryEntry) GetFormattedTime(env *environments.Environment) (string, error) {
	timestamp, err := he.GetTimestamp()
	if err != nil {
		return "", err
	}
	return env.FormatTimeStamp(timestamp)
}

// HasHistory tells if history entry has a trace stored.
func (he HistoryEntry) HasHistory() bool {
	return he.hasHistory
}

// ToString converts history entry to the string representation according to the environment setting.
func (he HistoryEntry) ToString(env *environments.Environment) string {
	command := he.command
	number := he.number

	timestamp := ""
	if formattedTimestamp, err := env.FormatTimeStamp(he.timestamp); err == nil {
		timestamp = "  (" + formattedTimestamp + ")"
	}

	history := markHasNoHistory
	if he.hasHistory {
		history = markHasHistory
	}

	return fmt.Sprintf("!%-5d%s %c  %s", number, timestamp, history, command)
}

// GetTraceName returns a trace name of the history entry.
func (he HistoryEntry) GetTraceName() string {
	digest := md5.New()
	binary.Write(digest, binary.LittleEndian, int64(he.timestamp))
	io.WriteString(digest, he.command)

	return fmt.Sprintf("%x", digest.Sum(nil))
}
