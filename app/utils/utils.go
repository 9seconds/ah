package utils

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func SplitCommandToChunks(cmd string) (string, []string) {
	chunks := strings.Split(cmd, " ")
	for idx := 0; idx < len(chunks); idx++ {
		chunks[idx] = strings.Trim(chunks[idx], " \t")
	}

	return chunks[0], chunks[1:]
}

func ConvertTimestamp(timestamp int) *time.Time {
	converted := time.Unix(int64(timestamp), 0)
	return &converted
}

func Open(filename string) *os.File {
	handler, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return handler
}

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
