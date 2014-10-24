package environments

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	logrus "github.com/Sirupsen/logrus"
	strftime "github.com/weidewang/go-strftime"

	"github.com/9seconds/ah/app/utils"
)

const (
	TRACES_DIR    = "traces"
	BOOKMARKS_DIR = "bookmarks"

	SHELL_ZSH  = "zsh"
	SHELL_BASH = "bash"

	DEFAULT_ZSH_HISTFILE  = ".zsh_history"
	DEFAULT_BASH_HISTFILE = ".bash_history"
	DEFAULT_APP_DIR       = ".ah"
)

var (
	CURRENT_USER *user.User
)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		os.Stderr.WriteString("Impossible to detect current user\n")
		os.Exit(1)
	}
	CURRENT_USER = currentUser
}

type Environment struct {
	appDir         string
	histFile       string
	histTimeFormat string
	shell          string
	log            *logrus.Logger
}

func (e *Environment) OK() bool {
	return e.log != nil && e.appDir != "" && e.histFile != "" && e.shell != ""
}

func (e *Environment) GetTracesDir() string {
	return filepath.Join(e.appDir, TRACES_DIR)
}

func (e *Environment) GetBookmarksDir() string {
	return filepath.Join(e.appDir, BOOKMARKS_DIR)
}

func (e *Environment) GetTraceFileName(hash string) string {
	return filepath.Join(e.GetTracesDir(), hash)
}

func (e *Environment) GetBookmarkFileName(name string) string {
	return filepath.Join(e.GetBookmarksDir(), name)
}

func (e *Environment) DiscoverAppDir() error {
	return e.SetAppDir(filepath.Join(CURRENT_USER.HomeDir, DEFAULT_APP_DIR))
}

func (e *Environment) GetAppDir() (string, error) {
	if e.appDir == "" {
		return "", errors.New("AppDir is not set yet")
	}
	return e.appDir, nil
}

func (e *Environment) SetAppDir(path string) error {
	e.appDir = path
	return nil
}

func (e *Environment) DiscoverShell() error {
	return e.SetShell(os.Getenv("SHELL"))
}

func (e *Environment) GetShell() (string, error) {
	if e.shell == "" {
		return "", errors.New("Shell is not set yet")
	}
	return e.shell, nil
}

func (e *Environment) SetShell(shell string) error {
	baseShell := path.Base(shell)

	if baseShell == SHELL_ZSH || baseShell == SHELL_BASH {
		e.shell = baseShell
		return nil
	}
	return fmt.Errorf("Shell %s is not supported", shell)
}

func (e *Environment) DiscoverHistFile() error {
	shell, err := e.GetShell()
	if err != nil {
		return err
	}

	if shell == SHELL_ZSH {
		return e.SetHistFile(filepath.Join(CURRENT_USER.HomeDir, DEFAULT_ZSH_HISTFILE))
	}
	return e.SetHistFile(filepath.Join(CURRENT_USER.HomeDir, DEFAULT_BASH_HISTFILE))
}

func (e *Environment) GetHistFile() (string, error) {
	if e.histFile == "" {
		return "", errors.New("HistFile is not set yet")
	}
	return e.histFile, nil
}

func (e *Environment) SetHistFile(path string) error {
	e.histFile = path
	return nil
}

func (e *Environment) GetHistTimeFormat() (string, error) {
	if e.histTimeFormat == "" {
		return "", errors.New("HistTimeFormat is not set yet")
	}
	return e.histTimeFormat, nil
}

func (e *Environment) SetHistTimeFormat(histTimeFormat string) error {
	e.histTimeFormat = histTimeFormat
	return nil
}

func (e *Environment) FormatTimeStamp(timestamp int) (string, error) {
	return e.FormatTime(utils.ConvertTimestamp(timestamp))
}

func (e *Environment) FormatTime(timestamp *time.Time) (string, error) {
	format, err := e.GetHistTimeFormat()
	if err != nil {
		return "", err
	}
	return strftime.Strftime(timestamp, format), nil
}

func (e *Environment) GetLogger() (*logrus.Logger, error) {
	if e.log == nil {
		return nil, errors.New("Logger is not set yet")
	}
	return e.log, nil
}

func (e *Environment) EnableDebugLog() {
	e.log = logrus.New()
	e.log.Out = os.Stderr
	e.log.Level = logrus.DebugLevel
}

func (e *Environment) DisableDebugLog() {
	e.log = logrus.New()
	e.log.Out = os.Stderr
	e.log.Level = logrus.ErrorLevel
}

func (e *Environment) GetTraceFilenames() ([]os.FileInfo, error) {
	fileInfos := make([]os.FileInfo, 0, 16)

	files, err := ioutil.ReadDir(e.GetTracesDir())
	if err != nil {
		return fileInfos, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileInfos = append(fileInfos, file)
	}

	return fileInfos, nil
}
