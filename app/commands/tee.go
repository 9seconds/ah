package commands

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/historyentries"
	"github.com/9seconds/ah/app/utils"
)

// Tee implements t (trace, tee) command.
func Tee(input []string, pseudoTTY bool, env *environments.Environment) {
	output, err := ioutil.TempFile(os.TempDir(), "ah")
	if err != nil {
		panic("Cannot create temporary file")
	}

	bufferedOutput := gzip.NewWriter(output)
	combinedStdout := io.MultiWriter(os.Stdout, bufferedOutput)
	combinedStderr := io.MultiWriter(os.Stderr, bufferedOutput)

	commandError := utils.Exec(pseudoTTY,
		os.Stdin, combinedStdout, combinedStderr,
		input[0], input[1:]...)

	bufferedOutput.Close()
	output.Close()

	preciseCommand, cmdErr := historyentries.GetCommands(historyentries.GetCommandsSingle, nil, env)
	if cmdErr != nil {
		panic("Sorry, cannot detect the number of the command")
	}
	cmd := preciseCommand.Result().(historyentries.HistoryEntry)

	commandHash := cmd.GetTraceName()
	os.Rename(output.Name(), env.GetTraceFileName(commandHash))

	if commandError != nil {
		os.Exit(utils.GetStatusCode(commandError))
	}
}
