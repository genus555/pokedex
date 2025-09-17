package main

import (
	"github.com/genus555/pokedex/internal/pokecache"
	"fmt"
	"net/http"
	"io"
	"encoding/json"
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

/*type Pokemon struct {
	pokemon struct {
		Name string
		URL  string
	}
}*/

func GetLocations (c *Config, cache *pokecache.Cache, inputs []string) error {
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

func GetPreviousLocations (c *Config, cache *pokecache.Cache, inputs []string) error {
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

func ExploreLocation (c *Config, cache *pokecache.Cache, input []string) error {
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

func CatchPokemon (c *Config, cache *pokecache.Cache, input []string) error {
	if len(input) < 2 {
		fmt.Println("No pokemon inputted")
	}
	fmt.Println("Working")
}