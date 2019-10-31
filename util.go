package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getImages(pokemonMap map[int]Pokemon) {

	for _, value := range pokemonMap {
		url := "http://pokedream.com/pokedex/images/sugimori/"
		number := value.number
		numberString := strconv.Itoa(number)

		if len(numberString) == 1 {

			numberString = "00" + numberString

		} else if len(numberString) == 2 {

			numberString = "0" + numberString

		}

		url = url + numberString + ".jpg"
		response, err := http.Get(url)
		handleError(err)

		file, err := os.Create("./images/" + value.name + ".jpg")
		_, err = io.Copy(file, response.Body)
		handleError(err)
		file.Close()
		response.Body.Close()

	}
}

func makeRequest(url string) string {

	resp, err := http.Get(url)
	handleError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	page := string(body)
	return page

}

func writeFile(fileName string, data string) {
	file, err := os.Create(fileName)
	handleError(err)
	defer file.Close()

	file.WriteString(data)

}

func separateFile(fileName string) {

	file, err := os.Open(fileName)
	handleError(err)
	defer file.Close()

	addToString := false
	writeString := ""
	count := 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "UILinkedTableRow") {

			addToString = true
			count++

		} else if addToString == true && strings.Contains(scanner.Text(), "</tr>") {

			addToString = false

			pokemonNumber := strconv.Itoa(count)

			fileName := "./PokemonNumber/" + pokemonNumber

			writeFile(fileName, writeString)

			writeString = ""

		}

		if addToString == true {
			writeString += scanner.Text() + "\n"
		}
	}
}

func writeJSON(pokemonMap map[int]Pokemon) {

	file, err := os.OpenFile("pokemon.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	defer file.Close()

	file.WriteString("{\n")

	for i := 1; i < len(pokemonMap); i++ {

		numberString := strconv.Itoa(pokemonMap[i].number)

		file.WriteString("   \"" + numberString + "\": {\n")
		file.WriteString("   \"number\": \"" + numberString + "\",\n")
		file.WriteString("   \"name\": \"" + pokemonMap[i].name + "\",\n")
		file.WriteString("   \"first_type\": \"" + pokemonMap[i].first_type + "\",\n")
		file.WriteString("   \"second_type\": \"" + pokemonMap[i].second_type + "\",\n")
		file.WriteString("   \"first_ability\": \"" + pokemonMap[i].first_ability + "\",\n")
		file.WriteString("   \"first_ability_description\": \"" + pokemonMap[i].first_ability_description + "\",\n")
		file.WriteString("   \"second_ability\": \"" + pokemonMap[i].second_ability + "\",\n")
		file.WriteString("   \"second_ability_description\": \"" + pokemonMap[i].second_ability_description + "\",\n")

		stats := pokemonMap[i].stats

		hp := strconv.Itoa(stats[0])
		attack := strconv.Itoa(stats[1])
		special_attack := strconv.Itoa(stats[2])
		defense := strconv.Itoa(stats[3])
		special_defense := strconv.Itoa(stats[4])
		speed := strconv.Itoa(stats[5])
		total := strconv.Itoa(stats[6])

		file.WriteString("   \"hp\": \"" + hp + "\",\n")
		file.WriteString("   \"atk\": \"" + attack + "\",\n")
		file.WriteString("   \"spAtk\": \"" + special_attack + "\",\n")
		file.WriteString("   \"df\": \"" + defense + "\",\n")
		file.WriteString("   \"spDf\": \"" + special_defense + "\",\n")
		file.WriteString("   \"spd\": \"" + speed + "\",\n")
		file.WriteString("   \"total\": \"" + total + "\"\n")

		if i+1 == len(pokemonMap) {
			file.WriteString("   }\n")
		} else {
			file.WriteString("   },\n")
		}

	}

	file.WriteString("}")
}
func createStruct(pokemonMap [9]string) Pokemon {

	number := getNumber(pokemonMap[4])
	name := getName(pokemonMap[1])
	first_type := getType(pokemonMap[6])
	second_type := getType(pokemonMap[7])

	var stats [7]int
	getStats(pokemonMap[8], &stats)

	first_ability, first_ability_description,
		second_ability, second_ability_description := getAbilities(name)

	pokemon := Pokemon{number, name, first_type, second_type,
		first_ability, first_ability_description,
		second_ability, second_ability_description, stats}
	return pokemon
}

func getAbilities(name string) (string, string, string, string) {

	var first_ability, first_ability_description,
		second_ability, second_ability_description string

	url := "http://pokedream.com/pokedex/pokemon/" + name

	page := makeRequest(url)

	writeFile("./PokemonAbilities/"+name, page)

	file, err := os.Open("./PokemonAbilities/" + name)
	handleError(err)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		currentLine := scanner.Text()

		if strings.Contains(currentLine, "/pokedex/abilities/") {

			start := strings.Index(currentLine, `name`) + 6
			end := strings.Index(currentLine, "</td>")

			first_ability = currentLine[start:end]

			start = strings.Index(currentLine, `text`) + 6
			end = strings.Index(currentLine, "</td></tr>")

			first_ability_description = currentLine[start:end]

			currentLine := currentLine[end:]

			if currentLine == "</td></tr>" {
				break
			}

			currentLine = strings.Replace(currentLine, "</td></tr>", "", 1)
			start = strings.Index(currentLine, `name`) + 6
			end = strings.Index(currentLine, "</td><td")

			second_ability = currentLine[start:end]

			start = strings.Index(currentLine, `text`) + 6
			end = strings.Index(currentLine, "</td></tr>")

			second_ability_description = currentLine[start:end]

		}

	}
	return first_ability, first_ability_description, second_ability, second_ability_description

}
