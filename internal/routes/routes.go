package routes

import (
	"net/http"
	"strconv"

	"github.com/George-Anagnostou/countries/internal/db"
	"github.com/George-Anagnostou/countries/internal/middleware"
	"github.com/George-Anagnostou/countries/internal/models"
	"github.com/George-Anagnostou/countries/internal/templates"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Start() {
	e := echo.New()
    e.Static("/static", "static")
	e.Use(echoMiddleware.Logger())
    e.Use(echoMiddleware.Recover())
    e.Use(middleware.InitSessionStore())
    e.Use(middleware.AuthMiddleware)
	e.Renderer = templates.NewTemplate()

    RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8000"))
}

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", getHome)

    e.GET("/register", getRegister)
    e.POST("/register", postRegister)

    e.GET("/login", getLogin)
    e.POST("/login", postLogin)

    e.GET("/users/:id", getUser)

    e.POST("/logout", postLogout)

	e.GET("/search_continents", getContinents)

	e.GET("/guess_countries", getGuessCountry)
	e.POST("/guess_countries", postGuessCountry)

	e.GET("/guess_capitals", getGuessCapital)
	e.POST("/guess_capitals", postGuessCapital)
}

func getUserFromContext(c echo.Context) *models.User {
    user, ok := c.Get("user").(*models.User)
    if !ok {
        return nil
    }
    return user
}

func getHome(c echo.Context) error {
    basePayload := models.NewBasePayload(getUserFromContext(c))
    return c.Render(200, "home", basePayload)
}

func getRegister(c echo.Context) error {
    basePayload := models.NewBasePayload(getUserFromContext(c))
	return c.Render(200, "register", basePayload)
}

func postRegister(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("password")
    err := db.AddUser(username, password)
    // basePayload := models.NewBasePayload()
    if err == db.ErrInvalidRegistration {
        return c.Redirect(301, "/register")
    }
    if err != nil {
        // return c.Render(500, "internalServerError", basePayload)
        return c.Redirect(301, "/register")
    }
    return c.Redirect(301, "/login")
}

func getLogin(c echo.Context) error {
    basePayload := models.NewBasePayload(getUserFromContext(c))
	return c.Render(200, "login", basePayload)
}

func postLogin(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("password")

    // basePayload := models.NewBasePayload(getUserFromContext(c))
    user, err := db.AuthenticateUser(username, password)
    if err == models.ErrInvalidLogin {
        // not the right way to handle this error...
        return c.JSON(http.StatusUnauthorized, echo.Map{
            "error": "invalid username or password",
        })
    }
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "error": "internal service error",
        })
    }
    sess, err := middleware.GetSession("session", c)
    if err != nil {
        // return c.Render(500, "internalServerError", basePayload)
        return c.Redirect(301, "/login")
    }
    // log.Printf("from routes: user.ID = %d", user.ID)
    sess.Values["userID"] = user.ID
    // log.Printf("from routes: sess.Values.userID = %v", sess.Values["userID"])
    if err = sess.Save(c.Request(), c.Response()); err != nil {
        // return c.Render(500, "internalServerError", basePayload)
        return c.Redirect(301, "/login")
    }

    // log.Printf("from routes: user logged in as %s", user.Username)
    return c.Redirect(301, "/")
}

func getUser(c echo.Context) error {
    userIDString := c.Param("id")
    userID, err := strconv.ParseInt(userIDString, 10, 64)
    if err != nil {
        return c.Redirect(301, "/")
    }
    queryUser, err := db.GetUserByID(userID)
    if err != nil {
        return echo.ErrNotFound
    }
    // check if the user page requested is the same
    // as the user requesting the page
    contextUser := getUserFromContext(c)
    if queryUser.ID != contextUser.ID {
        return echo.ErrUnauthorized
    }
    basePayload := models.NewBasePayload(contextUser)
    return c.Render(200, "user", basePayload)
}

func postLogout(c echo.Context) error {
    sess, _ := middleware.GetSession("session", c)
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
    continentsPayload := models.NewContinentPayload(continents, countries)
    basePayload := models.NewBasePayload(getUserFromContext(c))
    payload := models.CombinePayloads(continentsPayload, *basePayload)
	return c.Render(200, "search_continents", payload)
}

func getGuessCountry(c echo.Context) error {
	countries := models.Countries
	answerCountry := models.GetRandomCountry()
    var passed bool = false
    // set cookie as the country to guess
    answerCookie := middleware.SetCookie("answerCountryName", answerCountry.Name.CommonName)
    c.SetCookie(answerCookie)
    countriesPayload := models.NewCountriesPayload(countries, answerCountry, nil, passed)
    basePayload := models.NewBasePayload(getUserFromContext(c))
    payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_countries", payload)
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
    user := getUserFromContext(c)
    db.UpdateCountryScore(user.ID, passed)
    guessCountry := models.GetCountryByName(guessCountryName)
    countriesPayload := models.NewCountriesPayload(countries, answerCountry, guessCountry, passed)
    basePayload := models.NewBasePayload(user)
    payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_countries", payload)
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
    // are less specific for finding countries
    cookie := middleware.SetCookie("answerCountryCapital", answerCountry.Name.CommonName)
    c.SetCookie(cookie)
    var passed bool = false
    countriesPayload := models.NewCountriesPayload(countries, answerCountry, nil, passed)
    basePayload := models.NewBasePayload(getUserFromContext(c))
    payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_capitals", payload)
}

func postGuessCapital(c echo.Context) error {
	countries := models.Countries
    // get the target country from cookie
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
    user := getUserFromContext(c)
    db.UpdateCapitalScore(user.ID, passed)
    countriesPayload := models.NewCountriesPayload(countries, answerCountry, guessCountry, passed)
    basePayload := models.NewBasePayload(user)
    payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_capitals", payload)
}
