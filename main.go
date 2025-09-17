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
	callback	func(*Config, *pokecache.Cache, []string, map[string]Pokemon) error
}

type Config struct {
	next		*string
	previous	*string
	pageNum		int
}

var currentPage Config
var commandRegistry map[string]cliCommand
var UserPokedex map[string]Pokemon

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
		"catch": {
			name:			"catch",
			description:	"Add name of pokemon as an argument (catch pikachu) to attempt to catch the pokemon (pokemon catch rate based on base experience)",
			callback:		AttemptCapture,
		},
		"locate": {
			name:			"locate",
			description:	"Display list of locations the pokemon appears",
			callback:		LocatePokemon,
		},
		"inspect": {
			name:			"inspect",
			description:	"Gives stats and types of pokemon in possession",
			callback:		InspectPokemon,
		},
		"release": {
			name:			"release",
			description:	"Deletes data on pokemon from pokedex",
			callback:		ReleasePokemon,
		},
		"pokedex": {
			name:			"pokedex",
			description:	"List of pokemon user has caught",
			callback:		GetPokedex,
		},
	}
	nextURL := "https://pokeapi.co/api/v2/location-area"
	currentPage = Config {
		next:		&nextURL,
		previous:	nil,
		pageNum:	0,
	}
	UserPokedex = make(map[string]Pokemon)
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
		} else if err := cmd.callback(&currentPage, currentCache, inputList, UserPokedex); err != nil {
			fmt.Println(err)
		}
		if inputList[0] == "exit" {
			running = false
		}
	}
}