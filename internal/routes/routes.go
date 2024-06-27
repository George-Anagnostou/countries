package routes

import (
    "log"
    "net/http"

    "github.com/George-Anagnostou/countries/internal/db"
	"github.com/George-Anagnostou/countries/internal/models"
	// "github.com/George-Anagnostou/countries/internal/templates"
	"github.com/George-Anagnostou/countries/internal/middleware"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", getHome)

    e.GET("/register", getRegister)
    e.POST("/register", postRegister)

    e.GET("/login", getLogin)
    e.POST("/login", postLogin)

    e.POST("/logout", postLogout)

	e.GET("/search_continents", getContinents)

	e.GET("/guess_countries", getGuessCountry)
	e.POST("/guess_countries", postGuessCountry)

	e.GET("/guess_capitals", getGuessCapital)
	e.POST("/guess_capitals", postGuessCapital)
}

func agetHome(c echo.Context) error {
    sess, err := middleware.GetSession(c, "user-session")
    if err != nil {
        log.Print("error getting session")
    }
    untyped := sess.Values["username"]
    username, ok := untyped.(string)
    if !ok {
        log.Print("error getting username")
    }
    user, err := db.GetUserByUsername(username)
    if err != nil {
        log.Print("error getting user from username")
    }
    payload := models.NewUserPayload(user)
    pageData := models.NewPageData(payload)
	return c.Render(200, "home", pageData)
}

func getHome(c echo.Context) error {
    // var user *models.User
    // userID := c.Get("userID")
    // log.Printf("\n\n userID = %v", userID)
    // if userID == nil {
    //     user = nil
    // } else {
    //     var err error
    //     userID := userID.(int)
    //     user, err = db.GetUserByID(userID)
    //     if err != nil {
    //         return err
    //     }
    // }
    // payload := models.NewUserPayload(user)
    pageData := models.NewPageData()
    return c.Render(200, "home", pageData)
}

func getRegister(c echo.Context) error {
    pageData := models.NewPageData()
	return c.Render(200, "register", pageData)
}

func postRegister(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("password")
    err := db.AddUser(username, password)
    pageData := models.NewPageData()
    if err == db.ErrInvalidRegistration {
        return c.Render(200, "internalServerError", pageData)
    }
    if err != nil {
        return c.Render(500, "internalServerError", pageData)
    }
    return c.Redirect(301, "/login")
}

func getLogin(c echo.Context) error {
    pageData := models.NewPageData()
	return c.Render(200, "login", pageData)
}

func postLogin(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("password")

    pageData := models.NewPageData()
    user, err := db.AuthenticateUser(username, password)
    if err == models.ErrInvalidLogin {
        // not the right way to handle this error...
        // return c.Render(401, "unauthorized", pageData)
        return c.JSON(http.StatusUnauthorized, echo.Map{
            "error": "invalid username or password",
        })
    }
    if err != nil {
        // return c.Render(500, "internalServerError", pageData)
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "error": "internal service error",
        })
    }
    sess, err := middleware.GetSession(c, "user-session")
    if err != nil {
        return c.Render(500, "internalServerError", pageData)
    }
    sess.Values["username"] = user.Username
    sess.Save(c.Request(), c.Response())

    return c.Redirect(301, "/")
}

func postLogout(c echo.Context) error {
    sess, _ := middleware.GetSession(c, "user-session")
    sess.Options.MaxAge = -1
    sess.Save(c.Request(), c.Response())
    return c.Redirect(301, "/")
}

func getContinents(c echo.Context) error {
    continents := models.GetAllContinents()
    allCountries := models.GetAllCountries()
    filter := c.FormValue("continent")
    countries := models.FilterCountriesByContinent(allCountries, filter)
	sortMethod := c.FormValue("sort-method")
    models.SortCountries(countries, sortMethod)
    payload := models.NewContinentPayload(continents, countries)
    untyped := c.Get("userID")
    userID, _ := untyped.(int)
    pageData := models.NewPageData(payload)
    user, err := db.GetUserByID(userID)
    if err == nil {
        pageData.User = nil
    }
    pageData.User = user
	return c.Render(200, "search_continents", pageData)
}

func getGuessCountry(c echo.Context) error {
	countries := models.Countries
	answerCountry := models.GetRandomCountry()
    var passed bool = false
    // set cookie as the country to guess
    cookie := middleware.SetCookie("answerCountryName", answerCountry.Name.CommonName)
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
        cookie := middleware.ResetCookie("answerCountryName")
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
    cookie := middleware.SetCookie("answerCountryCapital", answerCountry.Name.CommonName)
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
                cookie := middleware.ResetCookie("answerCountryCapital")
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
