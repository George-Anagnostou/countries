package models

import (
	"encoding/json"
    "log"
    "math/rand"
    "os"
    "sort"
)

type PageData struct {
	FlagEmoji string
	Payload   interface{}
}

type CountryData struct {
    FlagEmoji   string      `json:"flag"`
    Continents  []string    `json:"continents"`
    Population  int         `json:"population"`
    Capitals    []string    `json:"capital"`
    Name struct {
        CommonName string   `json:"common"`
    }   `json:"name"`
}

// `countries` will be used for the entire life of the server and
// never change, so might as well read it once on start up
// var countries []map[string]CountryData = readCountries()
var Countries []CountryData = ReadCountries()

func ReadCountries() []CountryData {
	content, err := os.ReadFile("data/countries.json")
	if err != nil {
		log.Fatal("error reading countries data: ", err)
	}
	var payload []CountryData
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("error during unmarshal(): ", err)
	}
	return payload
}

func GetFlagEmoji() string {
	randFlag := Countries[rand.Intn(len(Countries))].FlagEmoji
	return randFlag
}

func GetRandomCountry() *CountryData {
    return &Countries[rand.Intn(len(Countries))]
}

func CountriesByName(slice []CountryData) {
    sort.Slice(slice, func(i, j int) bool {
        return slice[i].Name.CommonName < slice[j].Name.CommonName
    })
}

func CountriesByNameReverse(slice []CountryData) {
    sort.Slice(slice, func(i, j int) bool {
        return slice[i].Name.CommonName > slice[j].Name.CommonName
    })
}

func CountriesByPop(slice []CountryData) {
    sort.Slice(slice, func(i, j int) bool {
        return slice[i].Population < slice[j].Population
    })
}

func CountriesByPopReverse(slice []CountryData) {
    sort.Slice(slice, func(i, j int) bool {
        return slice[i].Population > slice[j].Population
    })
}

func GetCountryByName(name string) *CountryData {
    for _, country := range Countries {
        if name == country.Name.CommonName {
            return &country
        }
    }
    return &CountryData{}
}

func GetCountryByCapital(name string) []*CountryData {
    var countries []*CountryData
    for _, country := range Countries {
        for _, capital := range country.Capitals {
            if name == capital {
                countries = append(countries, &country)
            }
        }
    }
    return countries
}
