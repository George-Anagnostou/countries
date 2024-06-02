package routes

import (
    "log"
	"strconv"

    "github.com/George-Anagnostou/countries/internal/models"
    "github.com/George-Anagnostou/countries/internal/templates"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	e := echo.New()
	e.Static("/static", "static")
	e.Use(middleware.Logger())

	e.Renderer = templates.NewTemplate()

	e.GET("/", getHome)
	e.GET("/search_continents", getContinents)

	e.GET("/count", getCount)
	e.POST("/count", postCount)

	e.Logger.Fatal(e.Start(":8000"))
}

func getHome(c echo.Context) error {
	flagEmoji := models.GetFlagEmoji()
	pageData := models.PageData{
		FlagEmoji: flagEmoji,
	}
	return c.Render(200, "home", pageData)
}

func getContinents(c echo.Context) error {
    flagEmoji := models.GetFlagEmoji()
    var continents []string
    seen := make(map[string]bool)

    // populate datalist
    for _, country := range models.Countries {
        for _, continent := range country.Continents {
            if !seen[continent] {
                continents = append(continents, continent)
                seen[continent] = true
            }
        }
    }

    // populate countries
    var countries []models.CountryData
    filter := "All"
    filter = c.FormValue("continent")
    for _, country := range models.Countries {
        for _, continent := range country.Continents {
            if filter == "All" || filter == "" {
                countries = models.Countries
            } else {
                if continent == filter {
                    countries = append(countries, country)
                }
            }
        }
    }

    pageData := models.PageData{
        FlagEmoji: flagEmoji,
        Payload: struct {
            Continents []string
            Countries []models.CountryData
        }{
            Continents: continents,
            Countries: countries,
        },
    }
    return c.Render(200, "search_continents", pageData)
}

func getCount(c echo.Context) error {
    count := models.ReadCount()
	return c.Render(200, "count", count)
}

func postCount(c echo.Context) error {
    count := models.IncrementCount()
    strCount := strconv.Itoa(count.Count)
	return c.String(200, strCount)
}