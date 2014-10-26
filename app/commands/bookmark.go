package commands

import (
	"os"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

// Bookmark implements "b" (bookmark) command.
func Bookmark(commandNumber int, bookmarkAs string, env *environments.Environment) {
	if commandNumber < 0 {
		utils.Logger.Panic("Command number should be >= 0")
	}

	commandsKeeper, err := historyentries.GetCommands(historyentries.GetCommandsAll, nil, env)
	if err != nil {
		utils.Logger.Panic(err)
	}
	commands := commandsKeeper.Result().([]historyentries.HistoryEntry)
	if len(commands) <= commandNumber {
		utils.Logger.Panic("Command number does not exist")
	}
	command := commands[commandNumber-1]

	filename := env.GetBookmarkFileName(bookmarkAs)
	file, err := os.Create(filename)
	if err != nil {
		utils.Logger.Panicf("Cannot create bookmark %s: %v", filename, err)
	}
	defer file.Close()

	file.WriteString(command.GetCommand())
}
