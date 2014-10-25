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
		panic(err)
	}
	return handler
}

// Exec executes a command with arguments. Also it attaches ah incoming signal pipeline
// to external process for a great good.
// Return nil if ok. If not ok, returns exec.ExitError
func Exec(cmd string, args ...string) *exec.ExitError {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	AttachSignalsToProcess(command)

	err := command.Run()
	if err == nil {
		return nil
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError
	}
	panic(err.Error())
}

// GetStatusCode returns an exit code from exec.ExitError
func GetStatusCode(err *exec.ExitError) int {
	if err == nil {
		return 0
	}

	waitStatus, ok := err.Sys().(syscall.WaitStatus)
	if !ok {
		panic("It seems you have an unsupported OS")
	}
	return waitStatus.ExitStatus()
}

// RemoveWithLogging does the same as os.Remove does but logs.
func RemoveWithLogging(logger *logrus.Logger, fileName string) error {
	err := os.Remove(fileName)

	if err == nil {
		logger.WithFields(logrus.Fields{
			"filename": fileName,
		}).Info("File was deleted")
	} else {
		logger.WithFields(logrus.Fields{
			"filename": fileName,
			"error":    err,
		}).Warn("File was not deleted")
	}

	return nil
}
