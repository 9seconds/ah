package commands

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strconv"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

// ListTrace implements l command (list trace).
func ListTrace(argument string, env *environments.Environment) {
	number, err := strconv.Atoi(argument)
	if err != nil || number < 0 {
		panic(fmt.Sprintf("Cannot convert argument to a command number: %s", argument))
	}

	commands, err := historyentries.GetCommands(historyentries.GetCommandsPrecise, nil, env, number)
	if err != nil {
		panic(err)
	}
	command, _ := commands.Result().(historyentries.HistoryEntry)
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
	for scanner.Scan() {
		os.Stdout.WriteString(scanner.Text())
		os.Stdout.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
