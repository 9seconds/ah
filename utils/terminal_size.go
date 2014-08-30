package utils

import (
	"runtime"
	"syscall"
	"unsafe"
)

const (
	TIOCGWINSZ_NIX = 0x5413
	TIOCGWINSZ_OSX = 1074295912
)

type TIOCGWINSZResponse struct {
	Row    uint16
	Col    uint16
	XPixel uint16
	YPixel uint16
}

func getTerminalWindowResponse() (*TIOCGWINSZResponse, error) {
	response := new(TIOCGWINSZResponse)

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
	response, err := getTerminalWindowResponse()
	if err != nil {
		return 0, err
	}

	return int(response.Col), nil
}
