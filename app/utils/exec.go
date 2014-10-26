package utils

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	pty "github.com/kr/pty"
)

var (
	commandLock    = new(sync.Mutex)
	currentCommand *exec.Cmd
)

// Exec runs a command with connected streams and according to the TTY usage.
func Exec(pseudoTTY bool, stdin io.Reader, stdout io.Writer, stderr io.Writer, cmd string, args ...string) *exec.ExitError {
	command := exec.Command(cmd, args...)
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
	if currentCommand != nil {
		Logger.Panic("Command already executing")
	}
	commandLock.Lock()
	defer commandLock.Unlock()
	if currentCommand != nil {
		Logger.Panic("Command already executing")
	}

	currentCommand = command

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

			if currentCommand != nil {
				currentCommand.Process.Signal(incomingSignal)
			}
		}
	}()
}
