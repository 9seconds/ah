package app

import (
	"fmt"
	"os"
)

func CommandBookmark(commandNumber int, bookmarkAs string, env *Environment) {
	if commandNumber < 0 {
		panic("Command number should be >= 0")
	}

	commands, err := getCommands(nil, env)
	if err != nil {
		panic(err)
	}
	if len(commands) < commandNumber-1 {
		panic("Command number does not exist")
	}
	command := commands[commandNumber-1]
	filename := env.GetBookmarkFileName(bookmarkAs)

	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot create bookmark %s: %v", filename, err))
	}
	defer file.Close()

	file.WriteString(command.Command)
}
