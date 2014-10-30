package commands

import (
	"bufio"
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
	output, err := ioutil.TempFile(os.TempDir(), "ah")
	if err != nil {
		utils.Logger.Panic("Cannot create temporary file")
	}

	bufferedOutput := bufio.NewWriter(output)
	gzippedWrapper := utils.NewSynchronizedWriter(gzip.NewWriter(bufferedOutput))
	combinedStdout := io.MultiWriter(os.Stdout, gzippedWrapper)
	combinedStderr := io.MultiWriter(os.Stderr, gzippedWrapper)

	commandError := utils.Exec(pseudoTTY,
		os.Stdin, combinedStdout, combinedStderr,
		input[0], input[1:]...)

	gzippedWrapper.Close()
	bufferedOutput.Flush()
	output.Close()

	if hash, err := getPreciseHash(input, env); err == nil {
		err = os.Rename(output.Name(), env.GetTraceFileName(hash))
		if err != nil {
			utils.Logger.Errorf("Cannot save trace: %v", err)
		}
	} else {
		utils.Logger.Errorf("Error occured on fetching command number: %v", err)
	}

	if commandError != nil {
		os.Exit(utils.GetStatusCode(commandError))
	}
}

func getPreciseHash(cmd []string, env *environments.Environment) (hash string, err error) {
	commands, err := historyentries.GetCommands(historyentries.GetCommandsAll, getPreciseFilter(cmd), env)
	if err != nil {
		err = fmt.Errorf("Cannot fetch commands list: %v", err)
		return
	}
	commandList := commands.Result().([]historyentries.HistoryEntry)

	if len(commandList) == 0 {
		err = errors.New("Command list is empty")
		return
	}

	found := len(commandList) - 1
	for idx := len(commandList) - 2; idx >= 0; idx-- {
		if commandList[idx].GetTimestamp() < environments.CreatedAt {
			break
		}
		found = idx
	}
	hash = commandList[found].GetTraceName()

	return
}

func getPreciseFilter(cmd []string) *utils.Regexp {
	buffer := new(bytes.Buffer)
	for idx := 0; idx < len(cmd)-1; idx++ {
		buffer.WriteString(regexp.QuoteMeta(cmd[idx]))
		buffer.WriteString(`\s+`)
	}
	buffer.WriteString(regexp.QuoteMeta(cmd[len(cmd)-1]))

	return utils.CreateRegexp(buffer.String())
}
