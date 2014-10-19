package commands

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"../environments"
	"../history_entries"
	"../utils"
)

func Tee(input []string, env *environments.Environment) {
	output, err := ioutil.TempFile(os.TempDir(), "ah")
	if err != nil {
		panic("Cannot create temporary file")
	}
	bufferedOutput := bufio.NewWriter(output)

	command := exec.Command(input[0], input[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = io.MultiWriter(os.Stdout, bufferedOutput)
	command.Stderr = io.MultiWriter(os.Stdout, bufferedOutput)
	err = command.Run()

	bufferedOutput.Flush()
	output.Close()

	commands, err_ := history_entries.GetCommands(nil, env)
	if err_ != nil {
		panic("Sorry, cannot detect the number of the command")
	}
	os.Rename(output.Name(), env.GetTraceFileName(len(commands)))

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(utils.GetStatusCode(exitError))
		} else {
			panic(err.Error())
		}
	}
}
