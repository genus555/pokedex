package main
import (
	"bufio"
	"os"
	"fmt"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*Config) error
}

type Config struct {
	next		*string
	previous	*string
	pageNum		int
}

var currentPage Config
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
		"map": {
			name:			"map",
			description:	"Display next page of locations to explore",
			callback:		GetLocations,
		},
		"mapb": {
			name:			"mapb",
			description:	"Display previous page of locations to explore",
			callback:		GetPreviousLocations,
		},
	}
	nextURL := "https://pokeapi.co/api/v2/location-area"
	currentPage = Config {
		next:		&nextURL,
		previous:	nil,
		pageNum:	0,
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
			} else if err := cmd.callback(&currentPage); err != nil {
				fmt.Println(err)
			}
			if inp == "exit" {
				running = false
			}
		}
	}
}