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
	"github.com/9seconds/ah/app/history_entries"
	"github.com/9seconds/ah/app/utils"
)

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

	preciseCommand, err_ := history_entries.GetCommands(history_entries.GET_COMMANDS_SINGLE, nil, env)
	if err_ != nil {
		panic("Sorry, cannot detect the number of the command")
	}
	cmd := preciseCommand.Result().(history_entries.HistoryEntry)

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
	in_pty, in_tty, in_err := pty.Open()
	if in_err != nil {
		return nil, nil, nil, in_err
	}
	defer in_tty.Close()

	out_pty, out_tty, out_err := pty.Open()
	if out_err != nil {
		return nil, nil, nil, out_err
	}
	defer out_tty.Close()

	err_pty, err_tty, err_err := pty.Open()
	if err_err != nil {
		return nil, nil, nil, err_err
	}
	defer err_tty.Close()

	cmd.Stdin = in_tty
	cmd.Stdout = out_tty
	cmd.Stderr = err_tty
	cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	err := cmd.Start()
	if err != nil {
		in_pty.Close()
		out_pty.Close()
		err_pty.Close()
		return nil, nil, nil, err
	}

	return in_pty, out_pty, err_pty, nil
}
