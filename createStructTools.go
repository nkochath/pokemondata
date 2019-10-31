package main

import (
	"strconv"
	"strings"
)

func getNumber(data string) int {

	number := data
	number = strings.ReplaceAll(number, "<td>", "")
	number = strings.ReplaceAll(number, "</td>", "")
	pokemonNumber, err := strconv.Atoi(number)
	handleError(err)

	return pokemonNumber
}

func getName(data string) string {

	name := data

	if strings.Contains(name, "<td class=\"") {

		name = strings.ReplaceAll(name, "<td class=\"", "")
		name = strings.ReplaceAll(name, "\">", "")

	} else {

		name = strings.ReplaceAll(name, "<td class='", "")
		name = strings.ReplaceAll(name, "'>", "")

	}

	return name
}

func getType(data string) string {

	pokemonType := data

	if strings.Contains(pokemonType, `"`) {

		pokemonType = strings.ReplaceAll(pokemonType, "<td class=\"", "")
		slashIndex := strings.Index(pokemonType, "\"")
		pokemonType = pokemonType[:slashIndex]

	} else {

		pokemonType = ""

	}

	return pokemonType

}

func getStats(data string, stats *[7]int) {

	if stats[6] != 0 {
		return
	}

	data = strings.Replace(data, "<td>", "", 1)
	end := strings.Index(data, "</td>")
	collect := data[:end]
	data = data[end+5:]

	count := 0
	for {
		if stats[count] == 0 {
			intValue, err := strconv.Atoi(collect)
			handleError(err)
			stats[count] = intValue
			break
		}
		count++
	}

	getStats(data, stats)
}
