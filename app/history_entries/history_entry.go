package history_entries

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"../environments"
	"../utils"
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

func (he *HistoryEntry) GetNumber() (uint, error) {
	if he.number == 0 {
		return 0, errors.New("Number is not set yet")
	}
	return he.number, nil
}

func (he *HistoryEntry) GetCommand() (string, error) {
	if he.command == "" {
		return "", errors.New("Command is not set yet")
	}
	return he.command, nil
}

func (he *HistoryEntry) GetTimestamp() (int, error) {
	if he.timestamp == 0 {
		return 0, errors.New("Timestamp is not set yet")
	}
	return he.timestamp, nil
}

func (he *HistoryEntry) GetTime() (*time.Time, error) {
	timestamp, err := he.GetTimestamp()
	if err != nil {
		return nil, err
	}
	return utils.ConvertTimestamp(timestamp), nil
}

func (he *HistoryEntry) GetFormattedTime(env *environments.Environment) (string, error) {
	timestamp, err := he.GetTimestamp()
	if err != nil {
		return "", err
	}
	return env.FormatTimeStamp(timestamp)
}

func (he *HistoryEntry) HasHistory() bool {
	return he.hasHistory
}

func (he *HistoryEntry) ToString(env *environments.Environment) string {
	command, _ := he.GetCommand()
	number, _ := he.GetNumber()

	timestamp := ""
	if formattedTimestamp, err := he.GetFormattedTime(env); err == nil {
		timestamp = "\t" + formattedTimestamp
	}

	history := MARK_HAS_NO_HISTORY
	if he.HasHistory() {
		history = MARK_HAS_HISTORY
	}

	return fmt.Sprintf("!%d %c%s\t%s", number, history, timestamp, command)
}

func (he *HistoryEntry) GetTraceName() string {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, he.timestamp)
	io.WriteString(buffer, he.command)

	return fmt.Sprintf("%x", sha1.Sum(buffer.Bytes()))
}
