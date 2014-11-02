package utils

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	term "github.com/docker/docker/pkg/term"
	pty "github.com/kr/pty"
)

// Exec runs a command with connected streams and according to the TTY usage.
func Exec(cmd string, shell string, interactive bool, pseudoTTY bool, stdin io.Reader, stdout io.Writer, stderr io.Writer) *exec.ExitError {
	command := getCommand(cmd, interactive, shell)
	attachSignalsToProcess(command)

	var err error
	if pseudoTTY {
		err = runTtyCommand(command, stdin, stdout, stderr)
	} else {
		err = runStdCommand(command, stdin, stdout, stderr)
	}

	if err == nil {
		return nil
	} else if convertedError, ok := err.(*exec.ExitError); ok {
		return convertedError
	}
	Logger.Panic(err.Error())

	return nil // dammit, go!!!
}

func runStdCommand(command *exec.Cmd, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr

	if err := command.Start(); err != nil {
		Logger.Panic(err)
	}

	return command.Wait()
}

func runTtyCommand(command *exec.Cmd, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	pty, err := pty.Start(command)
	if err != nil {
		return err
	}
	defer pty.Close()

	hostFd := os.Stdin.Fd()
	oldTerminalState, err := term.SetRawTerminal(hostFd)
	if err != nil {
		return err
	}
	defer term.RestoreTerminal(hostFd, oldTerminalState)

	monitorTtyResize(hostFd, pty.Fd())

	go io.Copy(pty, stdin)
	go io.Copy(stdout, pty)

	return command.Wait()
}

func monitorTtyResize(hostFd uintptr, guestFd uintptr) {
	resizeTty(hostFd, guestFd)

	winchChan := make(chan os.Signal, 1)
	signal.Notify(winchChan, syscall.SIGWINCH)

	go func() {
		for _ = range winchChan {
			resizeTty(hostFd, guestFd)
		}
	}()
}

func resizeTty(hostFd uintptr, guestFd uintptr) {
	winsize, err := term.GetWinsize(hostFd)
	if err != nil {
		return
	}
	term.SetWinsize(guestFd, winsize)
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
