package history_entries

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
	MARK_HAS_HISTORY    = '*'
	MARK_HAS_NO_HISTORY = ' '
)

type HistoryEntry struct {
	number     uint
	command    string
	timestamp  int
	hasHistory bool
}

func (he HistoryEntry) GetNumber() (uint, error) {
	if he.number == 0 {
		return 0, errors.New("Number is not set yet")
	}
	return he.number, nil
}

func (he HistoryEntry) GetCommand() (string, error) {
	if he.command == "" {
		return "", errors.New("Command is not set yet")
	}
	return he.command, nil
}

func (he HistoryEntry) GetTimestamp() (int, error) {
	if he.timestamp == 0 {
		return 0, errors.New("Timestamp is not set yet")
	}
	return he.timestamp, nil
}

func (he HistoryEntry) GetTime() (*time.Time, error) {
	timestamp, err := he.GetTimestamp()
	if err != nil {
		return nil, err
	}
	return utils.ConvertTimestamp(timestamp), nil
}

func (he HistoryEntry) GetFormattedTime(env *environments.Environment) (string, error) {
	timestamp, err := he.GetTimestamp()
	if err != nil {
		return "", err
	}
	return env.FormatTimeStamp(timestamp)
}

func (he HistoryEntry) HasHistory() bool {
	return he.hasHistory
}

func (he HistoryEntry) ToString(env *environments.Environment) string {
	command := he.command
	number := he.number

	timestamp := ""
	if formattedTimestamp, err := he.GetFormattedTime(env); err == nil {
		timestamp = "\t" + formattedTimestamp
	}

	history := MARK_HAS_NO_HISTORY
	if he.hasHistory {
		history = MARK_HAS_HISTORY
	}

	return fmt.Sprintf("!%d %c%s\t%s", number, history, timestamp, command)
}

func (he HistoryEntry) GetTraceName() string {
	digest := md5.New()
	binary.Write(digest, binary.LittleEndian, int64(he.timestamp))
	io.WriteString(digest, he.command)

	return fmt.Sprintf("%x", digest.Sum(nil))
}
