package shell

// --- Imports

import (
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/codeskyblue/go-sh"

	"github.com/9seconds/ah/history"
)

// --- Consts

type ShellType uint8

const (
	SHELL_BASH ShellType = iota
	SHELL_ZSH  ShellType = iota
)

// --- Structs

type Shell struct {
	Type     ShellType
	Path     string
	BaseName string
	RC       string

	shellSession *sh.Session
	scanner      HistoryScannerInterface
}

// --- Methods

func (s *Shell) Discover() {
	s.Path = os.Getenv("SHELL")
	s.BaseName = s.getBaseName(s.Path)

	currentUser := s.getCurrentUser()

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
	msg, err := s.Run("echo -n $" + env)
	if err == nil {
		return msg
	} else {
		return ""
	}
}

func (s *Shell) GetHistoryScanner(reader io.Reader) HistoryScannerInterface {
	s.scanner.Init(reader)
	return s.scanner
}

func (s *Shell) getBaseName(shellPath string) string {
	absPath, absError := filepath.Abs(shellPath)

	if absError != nil {
		panic("Cannot understand what shell you are using")
	}
	if evalPath, evalErr := filepath.EvalSymlinks(absPath); evalErr == nil {
		absPath = evalPath
	}

	return path.Base(absPath)
}

func (s *Shell) getCurrentUser() user.User {
	currentUser, currentUserError := user.Current()

	if currentUserError != nil {
		panic("Cannot get a current user")
	}

	return currentUser
}
