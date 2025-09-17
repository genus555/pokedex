package main
import (
	"github.com/genus555/pokedex/internal/pokecache"
	"bufio"
	"os"
	"fmt"
	"time"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*Config, *pokecache.Cache, []string) error
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
		"explore": {
			name:			"explore",
			description:	"Add name of location as an arguement (explore city_A) to check area for pokemon",
			callback:		ExploreLocation,
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
	currentCache := pokecache.NewCache(15 *time.Second)
	currentCache.Add("test", nil)
	scanner := bufio.NewScanner(os.Stdin)
	
	running := true
	for running {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		inputList := CleanInput(input)
		cmd, ok := commandRegistry[inputList[0]]
		if !ok {
			fmt.Println("Unknown command")
		} else if err := cmd.callback(&currentPage, currentCache, inputList); err != nil {
			fmt.Println(err)
		}
		if inputList[0] == "exit" {
			running = false
		}
	}
}