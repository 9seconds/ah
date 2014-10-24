package commands

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	pty "github.com/kr/pty"

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
	command := exec.Command(input[0], input[1:]...)
	utils.AttachSignalsToProcess(command)

	if pseudoTTY {
		Stdin, Stdout, Stderr, ptyError := runPTYCommand(command)
		if ptyError != nil {
			output.Close()
			panic(ptyError)
		}
		go io.Copy(Stdin, os.Stdin)
		go io.Copy(combinedStdout, Stdout)
		go io.Copy(combinedStderr, Stderr)
	} else {
		command.Stdin = os.Stdin
		command.Stdout = combinedStdout
		command.Stderr = combinedStderr
		if startError := command.Start(); startError != nil {
			output.Close()
			panic(startError)
		}
	}

	commandError := command.Wait()

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
		if exitError, ok := commandError.(*exec.ExitError); ok {
			os.Exit(utils.GetStatusCode(exitError))
		} else {
			panic(commandError.Error())
		}
	}
}

func runPTYCommand(cmd *exec.Cmd) (*os.File, *os.File, *os.File, error) {
	inPTY, inTTY, inErr := pty.Open()
	if inErr != nil {
		return nil, nil, nil, inErr
	}
	defer inTTY.Close()

	outPTY, outTTY, outErr := pty.Open()
	if outErr != nil {
		return nil, nil, nil, outErr
	}
	defer outTTY.Close()

	errPTY, errTTY, errErr := pty.Open()
	if errErr != nil {
		return nil, nil, nil, errErr
	}
	defer errTTY.Close()

	cmd.Stdin = inTTY
	cmd.Stdout = outTTY
	cmd.Stderr = errTTY
	cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	err := cmd.Start()
	if err != nil {
		inPTY.Close()
		outPTY.Close()
		errPTY.Close()
		return nil, nil, nil, err
	}

	return inPTY, outPTY, errPTY, nil
}
