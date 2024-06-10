package routes

import (
    "net/http"

	"github.com/George-Anagnostou/countries/internal/models"
	"github.com/George-Anagnostou/countries/internal/templates"
	"github.com/George-Anagnostou/countries/internal/sessions"

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

	e.GET("/guess_countries", getGuessCountry)
	e.POST("/guess_countries", postGuessCountry)

	e.GET("/guess_capitals", getGuessCapital)
	e.POST("/guess_capitals", postGuessCapital)

	e.Logger.Fatal(e.Start(":8000"))
}

func getHome(c echo.Context) error {
    pageData := models.NewPageData()
	return c.Render(200, "home", pageData)
}

func getContinents(c echo.Context) error {
    continents := models.GetAllContinents()
    allCountries := models.GetAllCountries()
    filter := c.FormValue("continent")
    countries := models.FilterCountriesByContinent(allCountries, filter)
	sortMethod := c.FormValue("sort-method")
    models.SortCountries(countries, sortMethod)
    payload := models.NewContinentPayload(continents, countries)
    pageData := models.NewPageData(payload)
	return c.Render(200, "search_continents", pageData)
}

func getGuessCountry(c echo.Context) error {
	countries := models.Countries
	answerCountry := models.GetRandomCountry()
    var passed bool = false
    // set cookie as the country to guess
    cookie := sessions.SetCookie("answerCountryName", answerCountry.Name.CommonName)
    c.SetCookie(cookie)
    payload := models.NewCountriesPayload(countries, answerCountry, nil, passed)
    pageData := models.NewPageData(payload)
	return c.Render(200, "guess_countries", pageData)
}

func postGuessCountry(c echo.Context) error {
	countries := models.Countries
    // get the target country from cookie set with getGuessCountry
    answerCountryCookie, err := c.Cookie("answerCountryName")
    if err != nil {
        if err == http.ErrNoCookie {
            c.Redirect(301, "/guess_countries")
            return err
        }
        return err
    }
    answerCountry := models.GetCountryByName(answerCountryCookie.Value)
	guessCountryName := c.FormValue("country-guess")
    var passed bool = false
    if guessCountryName == answerCountry.Name.CommonName  {
        passed = true
        cookie := sessions.ResetCookie("answerCountryName")
        c.SetCookie(cookie)
    }
    guessCountry := models.GetCountryByName(guessCountryName)
    payload := models.NewCountriesPayload(countries, answerCountry, guessCountry, passed)
    pageData := models.NewPageData(payload)
	return c.Render(200, "guess_countries", pageData)
}

func getGuessCapital(c echo.Context) error {
	countries := models.Countries
    // don't use countries where capital == null
    var answerCountry = models.GetRandomCountry()
    for len(answerCountry.Capitals) < 1 {
        answerCountry = models.GetRandomCountry()
    }
    // get cookie as the country to guess
    // use CommonName as in getGuessCountry, since the Capitals
    // are less specific / determinite for finding countries
    cookie := sessions.SetCookie("answerCountryCapital", answerCountry.Name.CommonName)
    c.SetCookie(cookie)
    var passed bool = false
    payload := models.NewCountriesPayload(countries, answerCountry, nil, passed)
    pageData := models.NewPageData(payload)
	return c.Render(200, "guess_capitals", pageData)
}

func postGuessCapital(c echo.Context) error {
	countries := models.Countries
    // get the target country from cookie set with getGuessCapital
    answerCookie, err := c.Cookie("answerCountryCapital")
    if err != nil {
        if err == http.ErrNoCookie {
            c.Redirect(301, "/guess_capitals")
            return err
        }
        return err
    }
    answerCountry := models.GetCountryByName(answerCookie.Value)
	guessCapital := c.FormValue("guess-capital")
    guessCountries := models.GetCountryByCapital(guessCapital)
    var passed bool = false
    for range guessCountries {
        for _, capital := range answerCountry.Capitals {
            if capital == guessCapital {
                passed = true
                cookie := sessions.ResetCookie("answerCountryCapital")
                c.SetCookie(cookie)
            }
        }
    }
    var guessCountry *models.CountryData
    if len(guessCountries) < 1 {
        guessCountry = &models.CountryData{}
    } else {
        guessCountry = guessCountries[0]
    }
    payload := models.NewCountriesPayload(countries, answerCountry, guessCountry, passed)
    pageData := models.NewPageData(payload)
	return c.Render(200, "guess_capitals", pageData)
}
