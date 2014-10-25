package utils

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	logrus "github.com/Sirupsen/logrus"
)

// SplitCommandToChunks splits a command to the command name and
// its arguments.
func SplitCommandToChunks(cmd string) (string, []string) {
	chunks := strings.Split(cmd, " ")
	for idx := 0; idx < len(chunks); idx++ {
		chunks[idx] = strings.Trim(chunks[idx], " \t")
	}

	return chunks[0], chunks[1:]
}

// ConvertTimestamp converts timestamp to time structure
func ConvertTimestamp(timestamp int) *time.Time {
	converted := time.Unix(int64(timestamp), 0)
	return &converted
}

// Open just a small wrapper on os.Open which panics if something goes wrong.
func Open(filename string) *os.File {
	handler, err := os.Open(filename)
	if err != nil {
		Logger.Panic(err)
	}
	return handler
}

// GetStatusCode returns an exit code from exec.ExitError
func GetStatusCode(err *exec.ExitError) int {
	if err == nil {
		return 0
	}

	waitStatus, ok := err.Sys().(syscall.WaitStatus)
	if !ok {
		Logger.Panic("It seems you have an unsupported OS")
	}
	return waitStatus.ExitStatus()
}

// RemoveWithLogging does the same as os.Remove does but logs.
func RemoveWithLogging(fileName string) error {
	err := os.Remove(fileName)

	if err == nil {
		Logger.WithFields(logrus.Fields{
			"filename": fileName,
		}).Info("File was deleted")
	} else {
		Logger.WithFields(logrus.Fields{
			"filename": fileName,
			"error":    err,
		}).Warn("File was not deleted")
	}

	return nil
}
