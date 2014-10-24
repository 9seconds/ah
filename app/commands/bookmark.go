package commands

import (
	"fmt"
	"os"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
)

// Bookmark implements "b" (bookmark) command.
func Bookmark(commandNumber int, bookmarkAs string, env *environments.Environment) {
	if commandNumber < 0 {
		panic("Command number should be >= 0")
	}

	commandsKeeper, err := historyentries.GetCommands(historyentries.GetCommandsAll, nil, env)
	if err != nil {
		panic(err)
	}
	commands := commandsKeeper.Result().([]historyentries.HistoryEntry)
	if len(commands) <= commandNumber {
		panic("Command number does not exist")
	}

	command := commands[commandNumber-1]
	cmd, _ := command.GetCommand()
	filename := env.GetBookmarkFileName(bookmarkAs)

	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot create bookmark %s: %v", filename, err))
	}
	defer file.Close()

	file.WriteString(cmd)
}
