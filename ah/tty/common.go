package tty

// --- Imports

import (
	"runtime"
	"syscall"
	"unsafe"
)

// --- Consts

const (
	TIOCGWINSZ_NIX = 0x5413
	TIOCGWINSZ_OSX = 1074295912

	DEFAULT_TERMINAL_WIDTH = 80
)

// --- Vars

var TerminalWidth = DEFAULT_TERMINAL_WIDTH

// --- Structs

type tioResponse struct {
	Row    uint16
	Col    uint16
	XPixel uint16
	YPixel uint16
}

// --- Init

func init() {
	if width, err := GetTerminalWidth(); err == nil {
		TerminalWidth = width
	}
}

// --- Funcs

func getTerminalWidnowResponse() (*tioResponse, error) {
	response := new(tioResponse)

	tio := TIOCGWINSZ_NIX
	if runtime.GOOS == "darwin" {
		tio = TIOCGWINSZ_OSX
	}

	result, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tio),
		uintptr(unsafe.Pointer(response)))

	if int(result) == -1 {
		return nil, err
	}

	return response, nil
}

func GetTerminalWidth() (int, error) {
	if response, err := getTerminalWidnowResponse(); err == nil {
		return int(response.Col), nil
	} else {
		return 0, err
	}
}
