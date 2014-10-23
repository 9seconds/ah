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

	commands, err_ := history_entries.GetCommands(nil, env)
	if err_ != nil || number > len(commands) {
		panic("Sorry, cannot detect the trace of the command")
	}

	hashFilename := commands[number-1].GetTraceName()
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
