package commands

import (
	"fmt"
	"regexp"

	"../environments"
	"../history_entries"
	"../slices"
)

func Show(slice *slices.Slice, filter *regexp.Regexp, env *environments.Environment) {
	commands, err := history_entries.GetCommands(filter, env)
	sliceStart := slices.GetSliceIndex(slice.Start, len(commands))
	sliceFinish := slices.GetSliceIndex(slice.Finish, len(commands))

	if err != nil || sliceStart < 0 || sliceFinish < 0 || sliceFinish <= sliceStart {
		return
	}

	for _, value := range commands[sliceStart:sliceFinish] {
		fmt.Println(value.ToString(env))
	}
}
