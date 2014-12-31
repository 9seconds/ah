package environments

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	strftime "github.com/jehiah/go-strftime"
	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"

	"github.com/9seconds/ah/app/utils"
)

const (
	defaultAppDirName       = ".ah"
	defaultTracesDirName    = "traces"
	defaultBookmarksDirName = "bookmarks"

	defaultConfigFileName       = "config.yaml"
	defaultZshHistFileName      = ".zsh_history"
	defaultBashHistFileName     = ".bash_history"
	defaultAutoCommandsFileName = "autocommands.gob"

	// ShellBash defines code name of the Bash shell
	ShellBash = "bash"
	// ShellZsh defines code name of the Z Shell
	ShellZsh  = "zsh"
)

var (
	homeDir string

	defaultTmpDir = os.TempDir()

	// CreatedAt defines the timestamp when the launched
	CreatedAt = time.Now().Unix()
)

// Environment defines common structure which carries all information
// about environment where ah is executed.
type Environment struct {
	Shell          string `yaml:"shell"`
	HistFile       string `yaml:"histfile"`
	HistTimeFormat string `yaml:"histtimeformat"`

	HomeDir      string `yaml:"homedir"`
	AppDir       string `yaml:"appdir"`
	TmpDir       string `yaml:"tmpdir"`
	TracesDir    string `yaml:"tracesdir"`
	BookmarksDir string `yaml:"bookmarksdir"`

	AutoCommandsFileName string `yaml:"autocommands"`
	ConfigFileName       string `yaml:"config"`
}

func init() {
	currentHomeDir, err := homedir.Dir()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Cannot fetch your home directory: %v", err))
		os.Exit(1)
	}
	homeDir = currentHomeDir
}

// GetTraceFileName returns filename of the trace based on the given hash.
func (e *Environment) GetTraceFileName(hash string) string {
	return filepath.Join(e.TracesDir, hash)
}

// GetBookmarkFileName returns filename of the bookmark based on the given name.
func (e *Environment) GetBookmarkFileName(name string) string {
	return filepath.Join(e.BookmarksDir, name)
}

// GetHistFileName returns filename of the history file or error if something goes wrong (e.g unsupported shell).
func (e *Environment) GetHistFileName() (fileName string, err error) {
	fileName = e.HistFile

	if fileName == "" {
		switch e.Shell {
		case ShellBash:
			fileName = filepath.Join(e.HomeDir, defaultBashHistFileName)
		case ShellZsh:
			fileName = filepath.Join(e.HomeDir, defaultZshHistFileName)
		default:
			err = fmt.Errorf("Shell %s is not supported", e.Shell)
		}
	}

	return
}

// FormatTimeStamp is just a small wrapper around FormatTime method.
func (e *Environment) FormatTimeStamp(timestamp int64) string {
	return e.FormatTime(utils.ConvertTimestamp(timestamp))
}

// FormatTime formats a given time structure accroding to the settings.
func (e *Environment) FormatTime(timestamp *time.Time) (formatted string) {
	if e.HistTimeFormat != "" {
		formatted = strftime.Format(e.HistTimeFormat, *timestamp)
	}

	return
}

// GetTracesFileInfos returns file metadata structures on all traces.
func (e *Environment) GetTracesFileInfos() ([]os.FileInfo, error) {
	return e.getFileNames(e.TracesDir)
}

// GetBookmarksFileInfos returns file metadata structures on all bookmarks.
func (e *Environment) GetBookmarksFileInfos() ([]os.FileInfo, error) {
	return e.getFileNames(e.BookmarksDir)
}

func (e *Environment) getFileNames(directory string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileInfos = append(fileInfos, file)
	}

	return fileInfos, nil
}

// ReadFromConfig reads environment from config file.
func (e *Environment) ReadFromConfig() (configEnv *Environment, err error) {
	configEnv = new(Environment)

	content, err := ioutil.ReadFile(e.ConfigFileName)
	if err == nil {
		err = yaml.Unmarshal(content, configEnv)
	}

	return
}

func (e *Environment) String() string {
	return fmt.Sprintf("<Environment(shell='%s', histFile='%s', histTimeFormat='%s', homeDir='%s', appDir='%s', tracesDir='%s', bookmarksDir='%s', tmpDir='%s', configFileName='%s', autoCommandsFileName='%s')>",
		e.Shell,
		e.HistFile,
		e.HistTimeFormat,
		e.HomeDir,
		e.AppDir,
		e.TracesDir,
		e.BookmarksDir,
		e.TmpDir,
		e.ConfigFileName,
		e.AutoCommandsFileName)
}

// MakeDefaultEnvironment creates environment with default settings.
func MakeDefaultEnvironment() (env *Environment) {
	env = new(Environment)

	env.Shell = path.Base(os.Getenv("SHELL"))
	env.HistFile = os.Getenv("HISTFILE")
	env.HistTimeFormat = os.Getenv("HISTTIMEFORMAT")

	env.HomeDir = homeDir
	env.AppDir = filepath.Join(homeDir, defaultAppDirName)
	env.TracesDir = filepath.Join(env.AppDir, defaultTracesDirName)
	env.BookmarksDir = filepath.Join(env.AppDir, defaultBookmarksDirName)
	env.TmpDir = defaultTmpDir

	env.ConfigFileName = filepath.Join(env.AppDir, defaultConfigFileName)
	env.AutoCommandsFileName = filepath.Join(env.AppDir, defaultAutoCommandsFileName)

	return
}

// MergeEnvironments combines several environments into one.
func MergeEnvironments(envs ...*Environment) (result *Environment) {
	result = new(Environment)

	for _, value := range envs {
		result.Shell = getNotEmpty(result.Shell, value.Shell)
		result.HistFile = getNotEmpty(result.HistFile, value.HistFile)
		result.HistTimeFormat = getNotEmpty(result.HistTimeFormat, value.HistTimeFormat)
		result.HomeDir = getNotEmpty(result.HomeDir, value.HomeDir)
		result.AppDir = getNotEmpty(result.AppDir, value.AppDir)
		result.TracesDir = getNotEmpty(result.TracesDir, value.TracesDir)
		result.BookmarksDir = getNotEmpty(result.BookmarksDir, value.BookmarksDir)
		result.TmpDir = getNotEmpty(result.TmpDir, value.TmpDir)
		result.ConfigFileName = getNotEmpty(result.ConfigFileName, value.ConfigFileName)
		result.AutoCommandsFileName = getNotEmpty(result.AutoCommandsFileName, value.AutoCommandsFileName)
	}

	return
}

func getNotEmpty(first string, second string) string {
	if second != "" {
		return second
	}
	return first
}
