package historyentries

import (
	"bufio"
	"errors"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

// GetCommandsMode defines the mode GetCommands has to work.
type GetCommandsMode uint8

// Defines a list of possible GetCommandsMode variants.
const (
	GetCommandsAll GetCommandsMode = iota
	GetCommandsRange
	GetCommandsSingle
	GetCommandsPrecise
)

// GetCommands returns a keeper for the commands based on given mode and regular expression.
// varargs is the auxiliary list of numbers which makes sense in the context of GetCommandsMode setting
// only.
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

	commands, err := parser(keeper, scanner, filter, consumeChan)
	if err == nil {
		<-resultChan
		return commands, nil
	}
	return nil, err
}
