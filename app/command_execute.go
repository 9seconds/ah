package app

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func CommandExecuteCommandNumber(number int, env *Environment) {
	if number < 0 {
		panic("Cannot find such command")
	}

	commands, err := getCommands(nil, env)
	if err != nil {
		panic(err)
	}

	if len(commands) <= number {
		panic("Cannot find such command!")
	}
	command := commands[number-1]
	cmd, args := splitCommandToChunks(command.Command)

	toExecute := exec.Command(cmd, args...)
	toExecute.Stdout = os.Stdout
	toExecute.Stderr = os.Stderr
	toExecute.Stdin = os.Stdin
	err = toExecute.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(getStatusCode(exitError))
		} else {
			panic(err.Error())
		}
	}
}

func CommandExecuteBookMark(name string, env *Environment) {
	content, err := ioutil.ReadFile(env.GetBookmarkFileName(name))
	if err != nil {
		panic("Unknown bookmark")
	}

	cmd, args := splitCommandToChunks(string(content))

	toExecute := exec.Command(cmd, args...)
	toExecute.Stdout = os.Stdout
	toExecute.Stderr = os.Stderr
	toExecute.Stdin = os.Stdin
	err = toExecute.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(getStatusCode(exitError))
		} else {
			panic(err.Error())
		}
	}
}

func splitCommandToChunks(cmd string) (string, []string) {
	chunks := strings.Split(cmd, " ")
	for idx := 0; idx < len(chunks); idx++ {
		chunks[idx] = strings.Trim(chunks[idx], " \t")
	}

	return chunks[0], chunks[1:]
}
