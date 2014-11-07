package commands

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"

	"github.com/9seconds/ah/app/environments"
)

type autoCommand struct {
	Interactive bool
	PseudoTTY   bool
	Command     string
}

func (ac *autoCommand) String() string {
	interactive := 0
	pseudoTty := 0

	if ac.Interactive {
		interactive = 1
	}
	if ac.PseudoTTY {
		pseudoTty = 1
	}

	return fmt.Sprintf("%d %d %s", interactive, pseudoTty, ac.Command)
}

func AutoTeeList(env *environments.Environment) {
	autoCommands := getAutoCommands(env)

	keys := make([]string, 0, len(autoCommands))
	for cmd, _ := range autoCommands {
		keys = append(keys, cmd)
	}
	sort.Strings(keys)

	for idx := 0; idx < len(keys); idx++ {
		os.Stdout.WriteString(autoCommands[keys[idx]].String())
		os.Stdout.WriteString("\n")
	}
}

func AutoTeeAdd(commands []string, tty bool, interactive bool, env *environments.Environment) {
	commandsAlreadyHave := getAutoCommands(env)
	for _, cmd := range commands {
		if strct, ok := commandsAlreadyHave[cmd]; ok {
			strct.Interactive = interactive
			strct.PseudoTTY = tty
		} else {
			auto := autoCommand{Interactive: interactive, PseudoTTY: tty, Command: cmd}
			commandsAlreadyHave[cmd] = &auto
		}
	}

	file, err := os.Create(env.GetAutoCommandFileName())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(commandsAlreadyHave)
}

func getAutoCommands(env *environments.Environment) (commands map[string]*autoCommand) {
	file, err := os.Open(env.GetAutoCommandFileName())
	if err != nil {
		commands = make(map[string]*autoCommand)
	} else {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&commands)
		if err != nil {
			commands = nil
		}
		file.Close()
	}

	return
}
