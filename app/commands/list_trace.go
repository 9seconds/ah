package commands

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strconv"

	"../environments"
	"../history_entries"
	"../utils"
)

func ListTrace(argument string, env *environments.Environment) {
	number, err := strconv.Atoi(argument)
	if err != nil || number < 0 {
		panic(fmt.Sprintf("Cannot convert argument to a command number: %s", argument))
	}

	commands, err := history_entries.GetCommands(history_entries.GET_COMMANDS_PRECISE, nil, env, number)
	if err != nil {
		panic(err)
	}
	command, _ := commands.Result().(history_entries.HistoryEntry)
	hashFilename := command.GetTraceName()
	filename := env.GetTraceFileName(hashFilename)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		panic(fmt.Sprintf("Output for %s is not exist", argument))
	}

	file := utils.Open(filename)
	defer file.Close()
	ungzippedFile, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(ungzippedFile)

	// scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
