package routes

import (
    "net/http"
    "strconv"
    "time"

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

	e.GET("/guess_countries", getGuessCountry)
	e.POST("/guess_countries", postGuessCountry)

	e.GET("/guess_capitals", getGuessCapital)
	e.POST("/guess_capitals", postGuessCapital)

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
	sortMethod := c.FormValue("sort-method")
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
	switch sortMethod {
	case "alpha":
		models.CountriesByName(countries)
	case "alpha-reverse":
		models.CountriesByNameReverse(countries)
	case "pop":
		models.CountriesByPop(countries)
	case "pop-reverse":
		models.CountriesByPopReverse(countries)
	}
	pageData := models.PageData{
		FlagEmoji: flagEmoji,
		Payload: struct {
			Continents []string
			Countries  []models.CountryData
		}{
			Continents: continents,
			Countries:  countries,
		},
	}
	return c.Render(200, "search_continents", pageData)
}

func getGuessCountry(c echo.Context) error {
	flagEmoji := models.GetFlagEmoji()
	countries := models.Countries
	answerCountry := models.GetRandomCountry()
    // get cookie as the country to guess
    cookie := new(http.Cookie)
    cookie.Name = "answerCountryName"
    cookie.Value = answerCountry.Name.CommonName
    cookie.Expires = time.Now().Add(5 * time.Minute)
    c.SetCookie(cookie)
    var passed bool = false
	pageData := models.PageData{
		FlagEmoji: flagEmoji,
		Payload: struct {
			Countries     []models.CountryData
			AnswerCountry *models.CountryData
			GuessCountry *models.CountryData
            Passed bool
		}{
			Countries:     countries,
			AnswerCountry: answerCountry,
			GuessCountry: nil,
            Passed: passed,
		},
	}
	return c.Render(200, "guess_countries", pageData)
}

func postGuessCountry(c echo.Context) error {
	flagEmoji := models.GetFlagEmoji()
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
	guessCountry := models.GetCountryByName(guessCountryName)
    var passed bool = false
    if answerCountry.Name.CommonName == guessCountry.Name.CommonName {
        passed = true
        cookie := new(http.Cookie)
        cookie.Name = "answerCountryName"
        cookie.Expires = time.Now().Add(-time.Hour)
        c.SetCookie(cookie)
    }
	pageData := models.PageData{
		FlagEmoji: flagEmoji,
		Payload: struct {
			Countries     []models.CountryData
			AnswerCountry *models.CountryData
			GuessCountry  *models.CountryData
            Passed bool
		}{
			Countries:     countries,
			AnswerCountry: answerCountry,
			GuessCountry:  guessCountry,
            Passed: passed,
		},
	}
	return c.Render(200, "guess_countries", pageData)
}

func getGuessCapital(c echo.Context) error {
	flagEmoji := models.GetFlagEmoji()
	countries := models.Countries
    // don't use countries where capital == null
    var answerCountry = models.GetRandomCountry()
    for len(answerCountry.Capitals) < 1 {
        answerCountry = models.GetRandomCountry()
    }
    // get cookie as the country to guess
    cookie := new(http.Cookie)
    cookie.Name = "answerCountryCapital"
    cookie.Value = answerCountry.Name.CommonName
    cookie.Expires = time.Now().Add(5 * time.Minute)
    c.SetCookie(cookie)
    var passed bool = false
	pageData := models.PageData{
		FlagEmoji: flagEmoji,
		Payload: struct {
			Countries     []models.CountryData
			AnswerCountry *models.CountryData
			GuessCountry *models.CountryData
            Passed bool
		}{
			Countries:     countries,
			AnswerCountry: answerCountry,
			GuessCountry: nil,
            Passed: passed,
		},
	}
	return c.Render(200, "guess_capitals", pageData)
}

func postGuessCapital(c echo.Context) error {
	flagEmoji := models.GetFlagEmoji()
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
                cookie := new(http.Cookie)
                cookie.Name = "answerCountryCapital"
                cookie.Expires = time.Now().Add(-time.Hour)
                c.SetCookie(cookie)
            }
        }
    }
    guessCountry := guessCountries[0]
	pageData := models.PageData{
		FlagEmoji: flagEmoji,
		Payload: struct {
			Countries       []models.CountryData
			AnswerCountry   *models.CountryData
			GuessCountry    *models.CountryData
            Passed bool
		}{
			Countries:     countries,
			AnswerCountry: answerCountry,
			GuessCountry:  guessCountry,
            Passed: passed,
		},
	}
	return c.Render(200, "guess_capitals", pageData)
}
