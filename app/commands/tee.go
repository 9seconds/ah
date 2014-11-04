package commands

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

// Tee implements t (trace, tee) command.
func Tee(input string, interactive bool, pseudoTTY bool, env *environments.Environment) {
	output, err := ioutil.TempFile(env.GetTmpDir(), "ah")
	if err != nil {
		utils.Logger.Panic("Cannot create temporary file")
	}

	bufferedOutput := bufio.NewWriter(output)
	gzippedWrapper := utils.NewSynchronizedWriter(gzip.NewWriter(bufferedOutput))
	combinedStdout := io.MultiWriter(os.Stdout, gzippedWrapper)
	combinedStderr := io.MultiWriter(os.Stderr, gzippedWrapper)

	var commandError *exec.ExitError
	defer func() {
		// defer here because command may cause a panic but we do not want to lose any output
		gzippedWrapper.Close()
		bufferedOutput.Flush()
		output.Close()

		if hash, err := getPreciseHash(input, env); err == nil {
			err = os.Rename(output.Name(), env.GetTraceFileName(hash))
			if err != nil {
				utils.Logger.Errorf("Cannot save trace: %v. Get it here: %s", err, output.Name())
			} else {
				os.Remove(output.Name())
			}
		} else {
			utils.Logger.Errorf("Error occured on fetching command number: %v", err)
		}

		if commandError != nil {
			os.Exit(utils.GetStatusCode(commandError))
		}
	}()

	commandError = utils.Exec(input,
		string(env.GetShell()), interactive, pseudoTTY,
		os.Stdin, combinedStdout, combinedStderr)
}

func getPreciseHash(cmd string, env *environments.Environment) (hash string, err error) {
	commands, err := historyentries.GetCommands(historyentries.GetCommandsAll, nil, env)
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
