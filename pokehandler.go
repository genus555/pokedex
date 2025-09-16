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

func GetLocations(c *Config, cache *pokecache.Cache) error {
	url := *c.next
	info, exists := cache.Get(url)
	if exists {
		var locs Locations
		if err := json.Unmarshal(info, &locs); err != nil {
			return err
		}

		for _, location := range locs.Results {
			fmt.Println(location.Name)
		}
		c.previous = c.next
		c.next = locs.Next
		c.pageNum++
		fmt.Printf("PAGE: %d\n", c.pageNum)

	} else {
		res, err := http.Get(url)
		if err != nil {return err}
		defer res.Body.Close()
	
		data, err := io.ReadAll(res.Body)
		if err != nil {return err}

		var locs Locations
		if err := json.Unmarshal(data, &locs); err != nil {
			return err
		}

		for _, location := range locs.Results {
			fmt.Println(location.Name)
		}
		c.previous = c.next
		c.next = locs.Next
		c.pageNum++
		fmt.Printf("PAGE: %d\n", c.pageNum)
		
		locsJson, err := json.Marshal(locs)
		if err != nil {return err}
		cache.Add(url, locsJson)
	}
	return nil
}

func GetPreviousLocations(c *Config, cache *pokecache.Cache) error {
	if c.pageNum <= 1 {
		return fmt.Errorf("No previous page")
	}
	url := *c.previous
	info, exists := cache.Get(url)
	if exists {
		var locs Locations
		if err := json.Unmarshal(info, &locs); err != nil {
			return err
		}

		for _, location := range locs.Results {
			fmt.Println(location.Name)
		}

		c.next = c.previous
		c.previous = locs.Previous
		c.pageNum -= 1
		fmt.Printf("PAGE: %d\n", c.pageNum)
	} else {
		res, err := http.Get(url)
		if err != nil {return err}
		defer res.Body.Close()
	
		data, err := io.ReadAll(res.Body)
		if err != nil {return err}

		var locs Locations
		if err := json.Unmarshal(data, &locs); err != nil {
			return err
		}

		for _, location := range locs.Results {
			fmt.Println(location.Name)
		}
		c.next = c.previous
		c.previous = locs.Previous
		c.pageNum -= 1
		fmt.Printf("PAGE: %d\n", c.pageNum)

		locsJson, err := json.Marshal(locs)
		if err != nil {return err}
		cache.Add(url, locsJson)
	}
	return nil
}