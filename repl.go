package main
import "strings"

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