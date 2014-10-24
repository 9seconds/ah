package history_entries

import (
	"bufio"
	"errors"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

type GetCommandsMode uint8

const (
	GET_COMMANDS_ALL GetCommandsMode = iota
	GET_COMMANDS_RANGE
	GET_COMMANDS_SINGLE
	GET_COMMANDS_PRECISE
)

func GetCommands(mode GetCommandsMode, filter *utils.Regexp, env *environments.Environment, varargs ...int) (Keeper, error) {
	if !env.OK() {
		return nil, errors.New("Environment is not prepared")
	}

	resultChan, consumeChan := processHistories(env)
	keeper := getKeeper(mode, varargs...)
	parser := getParser(env)

	histFile, _ := env.GetHistFile()
	file := utils.Open(histFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if commands, err := parser(keeper, scanner, filter, consumeChan); err == nil {
		<-resultChan
		return commands, nil
	} else {
		return nil, err
	}
}
