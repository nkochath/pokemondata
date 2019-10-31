package main

import (
	"bufio"
	"os"
	"strconv"
)

type Pokemon struct {
	number                     int
	name                       string
	first_type                 string
	second_type                string
	first_ability              string
	first_ability_description  string
	second_ability             string
	second_ability_description string
	stats                      [7]int
}

var err error

func main() {

	pokemonMap := make(map[int]Pokemon)

	pokedex := "http://pokedream.com/pokedex/pokemon?display=gen1"

	page := makeRequest(pokedex)

	writeFile("Pokedex", page)

	separateFile("Pokedex")

	for i := 1; i <= 151; i++ {

		pokemonFile := "./PokemonNumber/" + strconv.Itoa(i)

		file, err := os.Open(pokemonFile)
		handleError(err)

		scanner := bufio.NewScanner(file)

		var pokemonData [9]string
		lineCount := 0

		for scanner.Scan() {
			pokemonData[lineCount] = scanner.Text()
			lineCount++
		}

		pokemonMap[i] = createStruct(pokemonData)
	}

	writeJSON(pokemonMap)
	getImages(pokemonMap)

}
