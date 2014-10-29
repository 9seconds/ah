package historyentries

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

const (
	markHasHistory   = '*'
	markHasNoHistory = ' '
)

// HistoryEntry stores a command with its context.
type HistoryEntry struct {
	number     uint
	command    string
	timestamp  int64
	hasHistory bool
}

// GetNumber returns a history number (may be executed with ! later).
func (he HistoryEntry) GetNumber() uint {
	return he.number
}

// GetCommand returns a command line which was executed.
func (he HistoryEntry) GetCommand() string {
	return he.command
}

// GetTimestamp returns a timestamp of the history entry.
func (he HistoryEntry) GetTimestamp() int64 {
	return he.timestamp
}

// GetTime returns a time structure of the history entry.
func (he HistoryEntry) GetTime() *time.Time {
	return utils.ConvertTimestamp(he.timestamp)
}

// GetFormattedTime returns formatted time stamp of the history entry.
func (he HistoryEntry) GetFormattedTime(env *environments.Environment) (string, error) {
	return env.FormatTimeStamp(he.timestamp)
}

// HasHistory tells if history entry has a trace stored.
func (he HistoryEntry) HasHistory() bool {
	return he.hasHistory
}

// String makes a string representation of the structure
func (he HistoryEntry) String() string {
	timestamp := utils.ConvertTimestamp(he.timestamp).Format(time.RFC3339)
	return fmt.Sprintf("HistoryEntry{number=%d command=\"%s\" timestamp=\"%s\" hasHistory=%t}",
		he.number, he.command, timestamp, he.hasHistory)
}

// ToString converts history entry to the string representation according to the environment setting.
func (he HistoryEntry) ToString(env *environments.Environment) string {
	timestamp := ""
	if formattedTimestamp, err := env.FormatTimeStamp(he.timestamp); err == nil {
		timestamp = "  (" + formattedTimestamp + ")"
	}

	history := markHasNoHistory
	if he.hasHistory {
		history = markHasHistory
	}

	return fmt.Sprintf("!%-5d%s %c  %s", he.number, timestamp, history, he.command)
}

// GetTraceName returns a trace name of the history entry.
func (he HistoryEntry) GetTraceName() string {
	digest := md5.New()
	binary.Write(digest, binary.LittleEndian, int64(he.timestamp))
	io.WriteString(digest, he.command)

	return fmt.Sprintf("%x", digest.Sum(nil))
}
