package homestorage

// --- Imports

import (
	"bufio"
	"os"
	"os/user"
	"path"
	"strings"
)

// --- Consts

const (
	StorageDir = ".ah"
	ConfigFile = "ah.conf"
)

// --- Vars

var HomeDir string

// --- Structs

type HomeStorage struct {
	content        map[string]string
	storageDir     string
	configFilePath string
	read           bool
}

// --- Init

func init() {
	currentUser, currentUserError := user.Current()
	if currentUserError != nil {
		panic("Impossible to detect current user")
	}
	HomeDir = currentUser.HomeDir
}

// --- Methods

func (hs *HomeStorage) Init() {
	hs.storageDir = path.Join(HomeDir, StorageDir)
	hs.configFilePath = path.Join(hs.storageDir, ConfigFile)
	hs.content = make(map[string]string)
	hs.read = false

	os.MkdirAll(hs.storageDir, os.ModeDir)
}

func (hs *HomeStorage) ReadConfig() {
	file, err := os.Open(hs.configFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(text, "=", 2)
		hs.content[parts[0]] = parts[1]
	}
}

func (hs *HomeStorage) GetKey(key string) (string, bool) {
	if value, ok := hs.content[key]; ok {
		return value, true
	}

	hs.ReadConfig()

	value, ok := hs.content[key]
	return value, ok
}

func (hs *HomeStorage) SetKey(key string, value string) {
	hs.ReadConfig()
	hs.content[key] = value

	file, err := os.Create(hs.configFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for key, value := range hs.content {
		_, _ = writer.WriteString(key + "=" + value + "\n")
	}
}
