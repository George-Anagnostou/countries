package models

import (
	"encoding/json"
    "math/rand"
    "log"
    "os"
)

var count Count = Count{Count: 0}

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

type Count struct {
	Count int
}

func ReadCount() *Count {
    return &count
}

func IncrementCount() *Count {
    count.Count++
    return &count
}

type PageData struct {
	FlagEmoji string
	Payload   interface{}
}

type CountryData struct {
    FlagEmoji string    `json:"flag"`
    Continents []string `json:"continents"`
    Population int      `json:"population"`
}

