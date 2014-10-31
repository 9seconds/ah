package commands

import (
	"io/ioutil"
	"os"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

// ExecuteCommandNumber executes command by its number in history file.
func ExecuteCommandNumber(number int, interactive bool, pseudoTTY bool, env *environments.Environment) {
	if number < 0 {
		utils.Logger.Panic("Cannot find such command")
	}

	commands, err := historyentries.GetCommands(historyentries.GetCommandsPrecise, nil, env, number)
	if err != nil {
		utils.Logger.Panic(err)
	}
	command := commands.Result().(historyentries.HistoryEntry)

	execute(command.GetCommand(), env.GetShell(), interactive, pseudoTTY)
}

// ExecuteBookmark executes command by its bookmark name.
func ExecuteBookmark(name string, interactive bool, pseudoTTY bool, env *environments.Environment) {
	content, err := ioutil.ReadFile(env.GetBookmarkFileName(name))
	if err != nil {
		utils.Logger.Panic("Unknown bookmark")
	}

	execute(string(content), env.GetShell(), interactive, pseudoTTY)
}

func execute(command string, shell string, interactive bool, pseudoTTY bool) {
	err := utils.Exec(command,
		shell, interactive, pseudoTTY,
		os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(utils.GetStatusCode(err))
	}
}
