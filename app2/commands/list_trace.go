package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"../environments"
	"../utils"
)

func ListTrace(argument string, env *environments.Environment) {
	number, err := strconv.Atoi(argument)
	if err != nil {
		panic(fmt.Sprintf("Cannot convert argument to a number: %s", argument))
	}

	filename := env.GetTraceFileName(number)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		panic(fmt.Sprintf("Output for %s is not exist", argument))
	}

	file := utils.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
