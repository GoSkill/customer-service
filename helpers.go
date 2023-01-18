package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// 1. помошники - чтение ресурсов
// чтение файла данных
func openAndReadCSV(file string) [][]string {
	var stringCSV [][]string
	fileCSV, err := os.Open(file)
	if err != nil {
		log.Fatalf("ошибка открытия файла: %s", err)
	}
	defer fileCSV.Close()
	reader := bufio.NewReader(fileCSV)
	for {
		line, err := reader.ReadString('\n')
		str := strings.TrimSpace(line) //удаляем символ "\n"
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println(err)
			}
		}
		strCSV := strings.Split((str), ";")
		stringCSV = append(stringCSV, strCSV)
	}
	return stringCSV
}

type Structures interface {
	[]MMSData
}

// чтение URL MMS
func parsingMMS(url string) []MMSData {
	var strURL []MMSData
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("получение ошибки, код состояния: %d \nbody: %s\n", resp.StatusCode, resp.Body)
	}
	if err = json.NewDecoder(resp.Body).Decode(&strURL); err != nil {
		log.Println("ERROR: " + err.Error())
	}
	return strURL
}

// чтение URL Support
func parsingSupport(url string) []SupportData {
	data := make([]SupportData, 0)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("получение ошибки, код состояния: %d \nbody: %s\n", resp.StatusCode, resp.Body)
		return data
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("ERROR: " + err.Error())
		return data
	}
	return data
}

// 2. валидатор (поиск соответствия)
func ValidData(dataSet []string, substr string) bool {
	for _, str := range dataSet {
		if str == substr {
			return true
		}
	}
	return false
}
func getCountriesList() []string {
	return []string{"RU", "US", "GB", "FR", "BL", "AT", "BG", "DK", "CA", "ES", "CH", "TR", "PE", "NZ", "MC"}
}

// 3. помошники - парсим ISO 3166-1 alpha-2
type ISOalpha2 struct { //Структура помогает разобрать полученные данные
	Code string `json:"code"`
	Name string `json:"name"`
}

// парсит ISO 3166-1 alpha-2
func ParseISOalpha2() []ISOalpha2 {
	var ResponseAlpha []ISOalpha2
	resp, err := http.Get("https://pkgstore.datahub.io/core/country-list/data_json/data/8c458f2d15d9f2119654b29ede6e45b8/data_json.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&ResponseAlpha); err != nil {
		log.Fatal(err)
	}
	return ResponseAlpha
}

// получает код страны из ISO 3166-1 alpha-2
func CodeISOalpha2() []string {
	var alfaCode []string
	for _, p := range ParseISOalpha2() {
		code := ISOalpha2(p)
		alfaCode = append(alfaCode, code.Code)
	}
	return alfaCode
}

// получает название страны из кода
func NameDecoding(list string) string {
	for _, p := range ParseISOalpha2() {
		code := ISOalpha2(p)
		if code.Code == list {
			list = code.Name
		}
	}
	return list
}
