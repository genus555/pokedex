package main
import (
	"bufio"
	"os"
	"fmt"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

var commandRegistry map[string]cliCommand

func init() {
	commandRegistry = map[string]cliCommand{
		"exit": {
			name:			"exit",
			description:	"Exit the Pokedex",
			callback:		CommandExit,
		},
		"help": {
			name:			"help",
			description:	"Manual for the Pokedex",
			callback:		CommandHelp,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	running := true
	for running {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		inputList := CleanInput(input)
		for _, inp := range inputList {
			cmd, ok := commandRegistry[inp]
			if !ok {
				fmt.Println("Unknown command")
			} else if err := cmd.callback(); err != nil {
				fmt.Println(err)
			}
			if inp == "exit" {
				running = false
			}
		}
	}
}