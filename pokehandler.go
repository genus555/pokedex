package main

import (
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

func GetLocations(c *Config) error {
	res, err := http.Get(*c.next)
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
	c.pageNum += 1
	fmt.Printf("PAGE: %d\n", c.pageNum)
	return nil
}

func GetPreviousLocations(c *Config) error {
	if c.pageNum <= 1 {
		return fmt.Errorf("No previous page")
	}
	res, err := http.Get(*c.previous)
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
	return nil
}