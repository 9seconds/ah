package commands

import (
	"io/ioutil"
	"os"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

// ExecuteCommandNumber executes command by its number in history file.
func ExecuteCommandNumber(pseudoTTY bool, number int, env *environments.Environment) {
	if number < 0 {
		utils.Logger.Panic("Cannot find such command")
	}

	commands, err := historyentries.GetCommands(historyentries.GetCommandsPrecise, nil, env, number)
	if err != nil {
		utils.Logger.Panic(err)
	}
	command := commands.Result().(historyentries.HistoryEntry)

	execute(pseudoTTY, command.GetCommand())
}

// ExecuteBookmark executes command by its bookmark name.
func ExecuteBookmark(pseudoTTY bool, name string, env *environments.Environment) {
	content, err := ioutil.ReadFile(env.GetBookmarkFileName(name))
	if err != nil {
		utils.Logger.Panic("Unknown bookmark")
	}

	execute(pseudoTTY, string(content))
}

func execute(pseudoTTY bool, command string) {
	cmd, args := utils.SplitCommandToChunks(command)

	executeError := utils.Exec(pseudoTTY,
		os.Stdin, os.Stdout, os.Stderr,
		cmd, args...)
	if executeError != nil {
		os.Exit(utils.GetStatusCode(executeError))
	}
}
