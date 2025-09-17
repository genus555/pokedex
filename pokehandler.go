package main

import (
	"github.com/genus555/pokedex/internal/pokecache"
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"math/rand"
)

type Locations struct {
	Count			int 	`json:"count"`
	Next			*string 	`json:"next"`
	Previous		*string 	`json:"previous"`
	Results			[]struct {
		Name		string	`json:"name"`
		URL			string	`json:"url"`
	} `json:"results"`
}

type LocationInformation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`

	GameIndex int `json:"game_index"`
	ID        int `json:"id"`

	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`

	Name  string `json:"name"`

	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`

	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int `json:"base_experience"`
	Height    int `json:"height"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Name          string `json:"name"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type EncounterArea []struct {
	LocationArea struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location_area"`
}

func GetLocations (c *Config, cache *pokecache.Cache, inputs []string, UserPokedex map[string]Pokemon) error {
	url := *c.next
	info, exists := cache.Get(url)
	var locs Locations
	if exists {
		if err := json.Unmarshal(info, &locs); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {return err}
		defer res.Body.Close()
	
		data, err := io.ReadAll(res.Body)
		if err != nil {return err}

		cache.Add(url, data)

		if err := json.Unmarshal(data, &locs); err != nil {
			return err
		}
	}

	for _, location := range locs.Results {
			fmt.Println(location.Name)
		}
		c.previous = c.next
		c.next = locs.Next
		c.pageNum++
		fmt.Printf("PAGE: %d\n", c.pageNum)

	return nil
}

func GetPreviousLocations (c *Config, cache *pokecache.Cache, inputs []string, UserPokedex map[string]Pokemon) error {
	if c.pageNum <= 1 {
		return fmt.Errorf("No previous page")
	}
	url := *c.previous
	info, exists := cache.Get(url)
	var locs Locations

	if exists {
		if err := json.Unmarshal(info, &locs); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {return err}
		defer res.Body.Close()
	
		data, err := io.ReadAll(res.Body)
		if err != nil {return err}

		cache.Add(url, data)

		if err := json.Unmarshal(data, &locs); err != nil {
			return err
		}
	}

	for _, location := range locs.Results {
			fmt.Println(location.Name)
		}
		c.next = c.previous
		c.previous = locs.Previous
		c.pageNum -= 1
		fmt.Printf("PAGE: %d\n", c.pageNum)

	return nil
}

func ExploreLocation (c *Config, cache *pokecache.Cache, input []string, UserPokedex map[string]Pokemon) error {
	if len(input) < 2 {
		return fmt.Errorf("No location provided")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", input[1])

	info, exists := cache.Get(url)
	var locationInfo LocationInformation

	if exists {
		if err := json.Unmarshal(info, &locationInfo); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {return err}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {return err}

		cache.Add(url, data)

		if err := json.Unmarshal(data, &locationInfo); err != nil {
			return err
		}
	}

	for _, pokemon := range locationInfo.PokemonEncounters {
			fmt.Println(pokemon.Pokemon.Name)
		}

	return nil
}

func getPokemon (pokemon string, cache *pokecache.Cache) (Pokemon, error) {

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)

	info, exists := cache.Get(url)
	var currentPokemon Pokemon

	if exists {
		if err := json.Unmarshal(info, &currentPokemon); err != nil {
			return Pokemon{}, err
		}
	} else {	
		res, err := http.Get(url)
		if err != nil {return Pokemon{}, err}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {return Pokemon{}, err}

		cache.Add(url, data)

		if err := json.Unmarshal(data, &currentPokemon); err != nil {
			return Pokemon{}, err
		}
	}

	return currentPokemon, nil
}

func AttemptCapture (c *Config, cache *pokecache.Cache, input []string, UserPokedex map[string]Pokemon) error {
	if len(input) < 2 {
		return fmt.Errorf("No pokemon inputted")
	}

	_, ok := UserPokedex[input[1]]
	if ok {
		return fmt.Errorf("You've already caught that pokemon")
	}

	currentPokemon, err := getPokemon(input[1], cache)
	if err != nil {return err}

	fmt.Printf("Throwing a Pokeball at %s...\n", currentPokemon.Name)

	catchChance := rand.Intn(currentPokemon.BaseExperience)
	if catchChance <= 50 {
		UserPokedex[currentPokemon.Name] = currentPokemon
		fmt.Printf("%s was caught!\n", currentPokemon.Name)
	} else {
		fmt.Printf("%s escaped!\n", currentPokemon.Name)
	}

	return nil
}

func ReleasePokemon (c *Config, cache *pokecache.Cache, input []string, UserPokedex map[string]Pokemon) error {
	if len(input) < 2 {
		return fmt.Errorf("No pokemon to release")
	}

	if _, ok := UserPokedex[input[1]]; ok {
		delete(UserPokedex, input[1])
		fmt.Printf("%s was released back into the wild...\n", input[1])
	} else {
		return fmt.Errorf("you have not caught that pokemon")
	}

	return nil
}

func LocatePokemon (c *Config, cache *pokecache.Cache, input []string, UserPokedex map[string]Pokemon) error {
	if len(input) < 2 {
		return fmt.Errorf("No pokemon to locate")
	}
	currentPokemon, err := getPokemon(input[1], cache)
	if err != nil {return err}

	url := currentPokemon.LocationAreaEncounters
	var currentPokemonLocations EncounterArea

	info, exists := cache.Get(url)
	if exists {
		if err := json.Unmarshal(info, &currentPokemonLocations); err != nil {
			return err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {return err}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {return err}

		cache.Add(url, data)

		if err := json.Unmarshal(data, &currentPokemonLocations); err != nil {
			return err
		}
	}

	for _, place := range currentPokemonLocations {
		fmt.Println(place.LocationArea.Name)
	}

	return nil
}

func InspectPokemon (c *Config, cache *pokecache.Cache, input []string, UserPokedex map[string]Pokemon) error {
	if len(input) < 2 {
		return fmt.Errorf("No pokmon to inspect")
	}

	currentPokemon, err := getPokemon(input[1], cache)
	if err != nil {return err}
	
	if _, ok := UserPokedex[currentPokemon.Name]; ok {
		fmt.Println("Name: ", currentPokemon.Name)
		fmt.Println("Height: ", currentPokemon.Height)
		fmt.Println("Weight: ", currentPokemon.Weight)
		fmt.Println("Stats:")
		
		for _, stat := range currentPokemon.Stats {
			fmt.Printf("   -%s: %d\n", stat.Stat.Name, stat.BaseStat)
			}

		fmt.Println("Types:")
		for _, pokeType := range currentPokemon.Types {
			fmt.Printf("   - %s\n", pokeType.Type.Name)
		}

		fmt.Println("Base Experience: ", currentPokemon.BaseExperience)

	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func GetPokedex (c *Config, cache *pokecache.Cache, input []string, UserPokedex map[string]Pokemon) error {
	if len(UserPokedex) == 0 {
		return fmt.Errorf("No pokemon has been caught yet")
	}

	fmt.Println("Your Pokedex:")
	for _, mon := range UserPokedex {
		fmt.Println("   - ", mon.Name)
	}
	
	return nil
}