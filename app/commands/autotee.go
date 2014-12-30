package commands

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"strings"

	logrus "github.com/Sirupsen/logrus"

	"github.com/9seconds/ah/app/environments"
	"github.com/9seconds/ah/app/utils"
)

type autoCommand struct {
	Interactive bool
	PseudoTTY   bool
	Command     string
}

func (ac *autoCommand) String() string {
	return fmt.Sprintf("%-20s [interactive=%-5t, pseudoTty=%-5t]",
		ac.Command, ac.Interactive, ac.PseudoTTY)
}

func (ac *autoCommand) Args(piped bool) string {
	buffer := bytes.NewBufferString(" ")

	if ac.PseudoTTY {
		buffer.WriteString("-y ")
	}
	if ac.Interactive || piped {
		buffer.WriteString("-x ")
	}

	return buffer.String()
}

// AutoTeeCreate creates a command to execute regarding to the auto tee
// information.
func AutoTeeCreate(command string, env *environments.Environment) {
	defer os.Stdout.WriteString("\n")

	command = strings.TrimSpace(command)
	key := strings.SplitN(command, " ", 2)[0]
	autoCommands := getAutoCommands(env)

	if auto, ok := autoCommands[key]; !ok || strings.Contains(command, ";") {
		os.Stdout.WriteString(command)
	} else {
		piped := strings.Contains(command, "|")
		fmt.Printf(`%s t%s-- "%s"`, os.Args[0], auto.Args(piped), command)
	}
}

// AutoTeeList returns a formatted list of the commands which should be
// executed with tee automatically. Basically output looks like
//
// ls                   [interactive=false, pseudoTty=false]
// python               [interactive=false, pseudoTty=false]
// ssh                  [interactive=false, pseudoTty=false]/
//
// where first column is interactive mode (-x) and second is pseudoTty (-y)
// 1 means true, 0 means false.
func AutoTeeList(env *environments.Environment) {
	autoCommands := getAutoCommands(env)

	keys := make([]string, 0, len(autoCommands))
	for cmd := range autoCommands {
		keys = append(keys, cmd)
	}
	sort.Strings(keys)

	for idx := 0; idx < len(keys); idx++ {
		os.Stdout.WriteString(autoCommands[keys[idx]].String())
		os.Stdout.WriteString("\n")
	}
}

// AutoTeeAdd adds a commands to the list of commands which should be executed
// automatically by tee.
func AutoTeeAdd(commands []string, tty bool, interactive bool, env *environments.Environment) {
	autoCommands := getAutoCommands(env)

	for _, cmd := range commands {
		if strct, ok := autoCommands[cmd]; ok {
			utils.Logger.WithFields(logrus.Fields{
				"autoCommand": strct.String(),
				"interactive": interactive,
				"pseudoTty":   tty,
			}).Info("Change command parameters")

			strct.Interactive = interactive
			strct.PseudoTTY = tty
		} else {
			auto := autoCommand{Interactive: interactive, PseudoTTY: tty, Command: cmd}
			autoCommands[cmd] = &auto

			utils.Logger.WithField("autoCommand", (&auto).String()).Info("Add new command")
		}
	}

	saveAutoTee(autoCommands, env)
}

// AutoTeeRemove removes a commands from the list of commands which should be executed
// automatically by tee.
func AutoTeeRemove(commands []string, env *environments.Environment) {
	autoCommands := getAutoCommands(env)

	for _, cmd := range commands {
		utils.Logger.WithField("command", cmd).Info("Remove command from the list")

		delete(autoCommands, cmd)
	}

	saveAutoTee(autoCommands, env)
}

func saveAutoTee(commands map[string]*autoCommand, env *environments.Environment) {
	file, err := os.Create(env.AutoCommandsFileName)
	if err != nil {
		utils.Logger.Panic(err)
	}
	defer file.Close()

	gob.NewEncoder(file).Encode(commands)
}

func getAutoCommands(env *environments.Environment) (commands map[string]*autoCommand) {
	file, err := os.Open(env.AutoCommandsFileName)
	if err != nil {
		utils.Logger.WithField("error", err).Warn("Cannot open auto tee commands file")
		commands = make(map[string]*autoCommand)
	} else {
		err = gob.NewDecoder(file).Decode(&commands)
		if err != nil {
			utils.Logger.WithField("error", err).Warn("Cannot decode GOB correctly")
			commands = nil
		}
		file.Close()
	}

	return
}
