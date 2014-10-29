package commands

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

// Tee implements t (trace, tee) command.
func Tee(input []string, pseudoTTY bool, env *environments.Environment) {
	hashChan := getPreciseHashChan(env)

	output, err := ioutil.TempFile(os.TempDir(), "ah")
	if err != nil {
		utils.Logger.Panic("Cannot create temporary file")
	}

	bufferedOutput := gzip.NewWriter(output)
	combinedStdout := io.MultiWriter(os.Stdout, bufferedOutput)
	combinedStderr := io.MultiWriter(os.Stderr, bufferedOutput)

	commandError := utils.Exec(pseudoTTY,
		os.Stdin, combinedStdout, combinedStderr,
		input[0], input[1:]...)

	bufferedOutput.Close()
	output.Close()

	result := <-hashChan
	if hash, ok := result.(string); ok {
		err := os.Rename(output.Name(), env.GetTraceFileName(hash))
		if err != nil {
			utils.Logger.Errorf("Cannot save trace: %v", err)
		}
	} else {
		utils.Logger.Errorf("Error occured on fetching command number: %v", result)
	}

	if commandError != nil {
		os.Exit(utils.GetStatusCode(commandError))
	}
}

func getPreciseHashChan(env *environments.Environment) (hashChan chan interface{}) {
	hashChan = make(chan interface{}, 1)

	go func() {
		commands, err := historyentries.GetCommands(historyentries.GetCommandsAll, getPreciseFilter(), env)
		if err != nil {
			hashChan <- fmt.Errorf("Cannot fetch commands list: %v", err)
			return
		}
		commandList := commands.Result().([]historyentries.HistoryEntry)

		if len(commandList) == 0 {
			hashChan <- errors.New("Command list is empty")
			return
		}

		found := len(commandList) - 1
		for idx := len(commandList) - 2; idx >= 0; idx-- {
			if commandList[idx].GetTimestamp() < environments.CreatedAt {
				break
			}
			found = idx
		}
		hashChan <- commandList[found].GetTraceName()
	}()

	return hashChan
}

func getPreciseFilter() *utils.Regexp {
	buffer := new(bytes.Buffer)
	for idx := 0; idx < len(os.Args)-1; idx++ {
		buffer.WriteString(regexp.QuoteMeta(os.Args[idx]))
		buffer.WriteString(`\s+`)
	}
	buffer.WriteString(regexp.QuoteMeta(os.Args[len(os.Args)-1]))

	return utils.CreateRegexp(buffer.String())
}
