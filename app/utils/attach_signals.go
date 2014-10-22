package utils

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

var (
	commandLock    *sync.Mutex = new(sync.Mutex)
	currentCommand *exec.Cmd   = nil
)

func AttachSignalsToProcess(command *exec.Cmd) {
	if currentCommand != nil {
		panic("Command already executing")
	}
	commandLock.Lock()
	defer commandLock.Unlock()

	currentCommand = command

	channel := make(chan os.Signal, 10)
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
		for {
			incomingSignal := <-channel

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
