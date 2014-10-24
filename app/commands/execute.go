package commands

import (
	"io/ioutil"
	"os"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/history_entries"
	"github.com/9seconds/ah/app/utils"
)

// ExecuteCommandNumber executes command by its number in history file.
func ExecuteCommandNumber(number int, env *environments.Environment) {
	if number < 0 {
		panic("Cannot find such command")
	}

	commands, err := history_entries.GetCommands(history_entries.GET_COMMANDS_PRECISE, nil, env, number)
	if err != nil {
		panic(err)
	}
	command, _ := commands.Result().(history_entries.HistoryEntry)
	cmd, _ := command.GetCommand()

	execute(cmd)
}

// ExecuteBookmark executes command by its bookmark name.
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
