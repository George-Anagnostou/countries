package models

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sort"
)

type Payload interface {
    SetBasePayload(BasePayload)
}

type BasePayload struct {
	FlagEmoji string
    User *User
}

func NewBasePayload(user *User) *BasePayload {
    basePayload := &BasePayload{
        FlagEmoji: getFlagEmoji(),
        User: user,
    }
    return basePayload
}

type UserPayload struct {
    User *User
    BasePayload
}

func NewUserPayload(user *User) *UserPayload {
    return &UserPayload{
        User: user,
    }
}

func (p *UserPayload) SetBasePayload(base BasePayload) {
    p.BasePayload = base
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

type CountriesPayload struct {
    Countries []CountryData
    AnswerCountry *CountryData
    GuessCountry *CountryData
    Passed bool
    BasePayload
}

func NewCountriesPayload(countries []CountryData, answerCountry *CountryData, guessCountry *CountryData, passed bool) *CountriesPayload {
    return &CountriesPayload{
        Countries: countries,
        AnswerCountry: answerCountry,
        GuessCountry: guessCountry,
        Passed: passed,
    }
}

func (p *CountriesPayload) SetBasePayload(base BasePayload) {
    p.BasePayload = base
}

type ContinentPayload struct {
    Continents []string
    Countries []CountryData
    BasePayload
}

func NewContinentPayload(continents []string, countries []CountryData) *ContinentPayload {
    return &ContinentPayload{
        Continents: continents,
        Countries: countries,
    }
}

func (p *ContinentPayload) SetBasePayload(base BasePayload) {
    p.BasePayload = base
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

func getFlagEmoji() string {
	randFlag := Countries[rand.Intn(len(Countries))].FlagEmoji
	return randFlag
}

func GetRandomCountry() CountryData {
    return Countries[rand.Intn(len(Countries))]
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

func GetAllCountries() []CountryData {
    return Countries
}

func GetCountryByName(name string) *CountryData {
    for _, country := range Countries {
        if name == country.Name.CommonName {
            return &country
        }
    }
    return DefaultCountryData()
}

func DefaultCountryData() *CountryData {
    return &CountryData{
        FlagEmoji: "",
        Continents: []string{""},
        Population: 0,
        Capitals: []string{""},
        Name: struct {
            CommonName string `json:"common"`
        }{
            CommonName: "",
        },
    }
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

func GetAllContinents() []string {
	var continents []string
	seen := make(map[string]bool)
	// populate datalist
	for _, country := range Countries {
		for _, continent := range country.Continents {
			if !seen[continent] {
				continents = append(continents, continent)
				seen[continent] = true
			}
		}
	}
    return continents
}

func FilterCountriesByContinent(countries []CountryData, filter string) []CountryData {
	var filteredCountries []CountryData
	for _, country := range countries {
		for _, continent := range country.Continents {
			if filter == "All" || filter == "" {
				filteredCountries = countries
			} else {
				if continent == filter {
					filteredCountries = append(filteredCountries, country)
				}
			}
		}
	}
    return filteredCountries
}

func SortCountries(countries []CountryData, sortMethod string) {
	switch sortMethod {
	case "alpha":
		CountriesByName(countries)
	case "alpha-reverse":
		CountriesByNameReverse(countries)
	case "pop":
		CountriesByPop(countries)
	case "pop-reverse":
		CountriesByPopReverse(countries)
	}
}

func CombinePayloads(specificPayload Payload, base BasePayload) Payload {
    specificPayload.SetBasePayload(base)
    return specificPayload
}

