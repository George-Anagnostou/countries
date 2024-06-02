package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
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

var count Count = Count{Count: 0}

// `countries` will be used for the entire life of the server and
// never change, so might as well read it once on start up
// var countries []map[string]CountryData = readCountries()
var countries []CountryData = readCountries()

func main() {
	e := echo.New()
	e.Static("/static", "static")
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	e.GET("/", getHome)
	e.GET("/search_continents", getContinents)

	e.GET("/count", getCount)
	e.POST("/count", postCount)

	e.Logger.Fatal(e.Start(":8000"))
}

func _readCountries() []map[string]CountryData {
	content, err := os.ReadFile("data/countries.json")
	if err != nil {
		log.Fatal("error reading countries data: ", err)
	}
	var payload []map[string]CountryData
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("error during unmarshal(): ", err)
	}
	return payload
}

func readCountries() []CountryData {
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
	randFlag := countries[rand.Intn(len(countries))].FlagEmoji
	return randFlag
}

func getHome(c echo.Context) error {
	flagEmoji := getFlagEmoji()
	pageData := PageData{
		FlagEmoji: flagEmoji,
	}
	return c.Render(200, "home", pageData)
}

func getContinents(c echo.Context) error {
    flagEmoji := getFlagEmoji()
    var continents []string
    seen := make(map[string]bool)

    for _, country := range countries {
        for _, continent := range country.Continents {
            if !seen[continent] {
                continents = append(continents, continent)
                seen[continent] = true
            }
        }
    }

    pageData := PageData{
        FlagEmoji: flagEmoji,
        Payload: struct {
            Continents []string
            Countries []CountryData
        }{
            Continents: continents,
            Countries: countries,
        },
    }
    return c.Render(200, "search_continents", pageData)
}

func getCount(c echo.Context) error {
	return c.Render(200, "count", count)
}

func postCount(c echo.Context) error {
	count.Count++
	return c.String(200, strconv.Itoa(count.Count))
}
