package utils


import (
	"os"
	"io"
	"os/user"
	"path"
	"strings"
	"path/filepath"

	"github.com/codeskyblue/go-sh"
)


type ShellType uint8
const (
	SHELL_BASH ShellType = iota
	SHELL_ZSH  ShellType = iota
)


type Shell struct {
	Type ShellType
	Path string
	BaseName string
	RC string

	shellSession *sh.Session
	scanner HistoryScannerInterface
}


func (s *Shell) Discover () {
	s.Path = os.Getenv("SHELL")

	absPath, absError := filepath.Abs(s.Path)
	if absError != nil {
		panic("Cannot understand what shell you are using")
	}
	if evalPath, evalError := filepath.EvalSymlinks(absPath); evalError == nil {
		absPath = evalPath
	}
	s.BaseName = path.Base(absPath)

	currentUser, currentUserError := user.Current()
	if currentUserError != nil {
		panic("Cannot get a current user")
	}

	if strings.Contains(s.BaseName, "zsh") {
		s.Type = SHELL_ZSH
		s.RC = path.Join(currentUser.HomeDir, ".zshrc")
		s.scanner = new(historyScannerZsh)
	} else if strings.Contains(s.BaseName, "bash") {
		s.Type = SHELL_BASH
		s.RC = path.Join(currentUser.HomeDir, ".bashrc")
		s.scanner = new(historyScannerBash)
	} else {
		panic("Unknown shell type. ah supports only bash and zsh")
	}

	s.shellSession = sh.NewSession()
}

func (s *Shell) Run(command string) (string, error) {
	msg, err := s.shellSession.Command(s.Path, "-i", "-c", command).Output()
	output := string(msg)
	return output, err
}

func (s *Shell) GetEnv(env string) string {
	fromEnv := os.Getenv(env)
	if fromEnv != "" {
		return fromEnv
	}

	msg, err := s.Run("echo -n $" + env)
	if err == nil {
		return msg
	} else {
		return ""
	}
}

func (s *Shell) GetHistoryScanner (reader io.Reader) HistoryScannerInterface {
	s.scanner.Init(reader)
	return s.scanner
}

