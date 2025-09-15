package main
import (
	"strings"
	"fmt"
)

func CleanInput(input string) []string {
	lowerInput := strings.ToLower(input)
	tmp := strings.Split(lowerInput, " ")
	var inputList []string
	for _, str := range tmp {
		if str != "" {
			inputList = append(inputList, str)
		}
	}
	return inputList
}

func CommandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func CommandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}