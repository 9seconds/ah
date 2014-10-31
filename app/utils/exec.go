package utils

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	pty "github.com/kr/pty"
)

// Exec runs a command with connected streams and according to the TTY usage.
func Exec(cmd string, shell string, interactive bool, pseudoTTY bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) *exec.ExitError {
	command := getCommand(cmd, interactive, shell)
	attachSignalsToProcess(command)

	if pseudoTTY {
		pseudoStdin, pseudoStdout, pseudoStderr, pseudoPTYError := runPTYCommand(command)
		if pseudoPTYError != nil {
			Logger.Panic(pseudoPTYError)
		}
		go io.Copy(pseudoStdin, stdin)
		go io.Copy(stdout, pseudoStdout)
		go io.Copy(stderr, pseudoStderr)
	} else {
		command.Stdin = stdin
		command.Stdout = stdout
		command.Stderr = stderr
		if startError := command.Start(); startError != nil {
			Logger.Panic(startError)
		}
	}

	err := command.Wait()
	if err == nil {
		return nil
	} else if convertedError, ok := err.(*exec.ExitError); ok {
		return convertedError
	}
	Logger.Panic(err.Error())

	return nil // dammit, go!!!
}

func runPTYCommand(cmd *exec.Cmd) (inPTY *os.File, outPTY *os.File, errPTY *os.File, err error) {
	inPTY, inTTY, err := pty.Open()
	if err != nil {
		return
	}
	defer inTTY.Close()

	outPTY, outTTY, err := pty.Open()
	if err != nil {
		return
	}
	defer outTTY.Close()

	errPTY, errTTY, err := pty.Open()
	if err != nil {
		return
	}
	defer errTTY.Close()

	cmd.Stdin = inTTY
	cmd.Stdout = outTTY
	cmd.Stderr = errTTY
	cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	err = cmd.Start()
	if err != nil {
		inPTY.Close()
		outPTY.Close()
		errPTY.Close()
	}

	return
}

func attachSignalsToProcess(command *exec.Cmd) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTSTP,
		syscall.SIGCONT,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
		syscall.SIGBUS,
		syscall.SIGSEGV)

	go func() {
		for incomingSignal := range channel {
			switch incomingSignal {
			case syscall.SIGBUS:
				fallthrough
			case syscall.SIGSYS:
				fallthrough
			case syscall.SIGSEGV:
				incomingSignal = syscall.SIGKILL
			}

			if command != nil {
				command.Process.Signal(incomingSignal)
			}
		}
	}()
}

func getCommand(cmd string, interactive bool, shell string) (command *exec.Cmd) {
	if interactive {
		command = exec.Command(shell, "-i", "-c", cmd)
	} else {
		cmd, args := SplitCommandToChunks(cmd)
		command = exec.Command(cmd, args...)
	}

	return
}
