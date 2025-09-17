package main
import (
	"github.com/genus555/pokedex/internal/pokecache"
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

func CommandExit(c *Config, cache *pokecache.Cache, inputs []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return nil
}

func CommandHelp(c *Config, cache *pokecache.Cache, inputs []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println("\nWill read first command. For example:\n\"asd exit\" will not work but \"exit asd\" will exit.")
	return nil
}