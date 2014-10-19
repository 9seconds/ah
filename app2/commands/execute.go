package commands

import (
	"io/ioutil"
	"os"

	"../environments"
	"../history_entries"
	"../utils"
)

func ExecuteCommandNumber(number int, env *environments.Environment) {
	if number < 0 {
		panic("Cannot find such command")
	}

	commands, err := history_entries.GetCommands(nil, env)
	if err != nil {
		panic(err)
	}

	if len(commands) <= number {
		panic("Cannot find such command")
	}
	command, _ := commands[number-1].GetCommand()

	execute(command)
}

func ExecuteBookmark(name string, env *environments.Environment) {
	content, err := ioutil.ReadFile(env.GetBookmarkFileName(name))
	if err != nil {
		panic("Unknown bookmark")
	}

	execute(string(content))
}

func execute(command string) {
	cmd, args := utils.SplitCommandToChunks(command)
	executeError := utils.Exec(cmd, args...)

	if executeError != nil {
		os.Exit(utils.GetStatusCode(executeError))
	}
}
