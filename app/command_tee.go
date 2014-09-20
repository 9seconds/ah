package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func watchPipe(source io.Reader, sink io.Writer, copyTo chan []byte) {
	defer close(copyTo)

	sourceSink := bufio.NewReader(io.TeeReader(source, sink))
	for {
		content, err := sourceSink.ReadBytes('\n')
		copyTo <- content
		if err != nil {
			return
		}
	}
}

func connectAllPipes(stdout io.Reader, stderr io.Reader, toCopy io.Writer, done chan bool) {
	defer close(done)

	storage := bufio.NewWriter(toCopy)
	defer storage.Flush()

	stdoutChan := make(chan []byte)
	stderrChan := make(chan []byte)

	go watchPipe(stdout, os.Stdout, stdoutChan)
	go watchPipe(stderr, os.Stderr, stderrChan)

	for {
		select {
		case msg, ok := <-stdoutChan:
			storage.Write(msg)
			if !ok {
				stdoutChan = nil
			}
		case msg, ok := <-stderrChan:
			storage.Write(msg)
			if !ok {
				stderrChan = nil
			}
		}

		if stdoutChan == nil && stderrChan == nil {
			return
		}
	}
}

func getStatusCode(err *exec.ExitError) int {
	if err == nil {
		return 0
	}

	waitStatus, ok := err.Sys().(syscall.WaitStatus)
	if !ok {
		panic("It seems you have an unsupported OS")
	}
	return waitStatus.ExitStatus()
}

func getCommandArgs(input []string) (string, []string) {
	if len(input) < 1 {
		panic("What command do you want to execute?")
	}

	command := input[0]
	args := make([]string, 0, len(input))

	for _, value := range input[1:] {
		arg := strings.Replace(value, `\`, `\\`, -1)
		arg = strings.Replace(arg, `"`, `\"`, -1)
		args = append(args, arg)
	}

	return command, args
}

func CommandTee(input []string) {
	command, args := getCommandArgs(input)
	fmt.Printf("\"%s %s\"\n", command, args)

	toExecute := exec.Command(command, args...)
	toExecuteStderr, _ := toExecute.StderrPipe()
	toExecuteStdout, _ := toExecute.StdoutPipe()

	output, _ := os.Create("output.log")
	defer output.Close()

	doneChan := make(chan bool, 1)
	go connectAllPipes(toExecuteStdout, toExecuteStderr, output, doneChan)

	err := toExecute.Run()
	<-doneChan
	if err != nil {
		os.Exit(getStatusCode(err.(*exec.ExitError)))
	}
}
