package app

import (
	"io/ioutil"
	"os"
	"os/exec"

	"../utils"
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
	cmd, args := utils.SplitCommandToChunks(command.Command)

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

	cmd, args := utils.SplitCommandToChunks(string(content))

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
