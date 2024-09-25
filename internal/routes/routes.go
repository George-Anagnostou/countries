package routes

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/George-Anagnostou/countries/internal/middleware"
	"github.com/George-Anagnostou/countries/internal/models"
	"github.com/George-Anagnostou/countries/internal/repositories"
	"github.com/George-Anagnostou/countries/internal/services"
	"github.com/George-Anagnostou/countries/internal/templates"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Start() {

	var db, err = sql.Open("sqlite3", "./data/countries.db")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := &repositories.SQLiteUserRepository{DB: db}
	userRepo.Initialize()
	userService := &services.UserService{UserRepo: userRepo}

	e := echo.New()
	e.Static("/static", "static")
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.InitSessionStore())
	e.Use(middleware.AuthMiddleware(*userService))
	e.Renderer = templates.NewTemplate()

	RegisterRoutes(e, *userService)
	e.Logger.Fatal(e.Start(":8000"))
}

func RegisterRoutes(e *echo.Echo, userService services.UserService) {
	e.GET("/", func(c echo.Context) error {
		return getHome(c)
	})

	e.GET("/register", func(c echo.Context) error {
		return getRegister(c)
	})
	e.POST("/register", func(c echo.Context) error {
		return postRegister(c, userService)
	})

	e.GET("/login", func(c echo.Context) error {
		return getLogin(c)
	})
	e.POST("/login", func(c echo.Context) error {
		return postLogin(c, userService)
	})

	e.GET("/users/:id", func(c echo.Context) error {
		return getUser(c, userService)
	})

	e.GET("/leaderboard", func(c echo.Context) error {
		return getLeaderboard(c, userService)
	})

	e.POST("/logout", func(c echo.Context) error {
		return postLogout(c)
	})

	e.GET("/search_continents", func(c echo.Context) error {
		return getContinents(c)
	})

	e.GET("/guess_countries", func(c echo.Context) error {
		return getGuessCountry(c)
	})
	e.POST("/guess_countries", func(c echo.Context) error {
		return postGuessCountry(c, userService)
	})

	e.GET("/guess_capitals", func(c echo.Context) error {
		return getGuessCapital(c)
	})
	e.POST("/guess_capitals", func(c echo.Context) error {
		return postGuessCapital(c, userService)
	})

	e.POST("/skip_country", func(c echo.Context) error {
		return postSkipCountry(c, userService)
	})
	e.POST("/skip_capital", func(c echo.Context) error {
		return postSkipCapital(c, userService)
	})
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

func postRegister(c echo.Context, userService services.UserService) error {
	basePayload := models.NewBasePayload(getUserFromContext(c))
	username := c.FormValue("username")
	password := c.FormValue("password")
	err := userService.AddUser(username, password)
	if err == repositories.ErrInvalidRegistration {
		basePayload.AddError(err)
		return c.Render(200, "register", basePayload)
	}
	if err != nil {
		basePayload.AddError(err)
		return c.Render(200, "register", basePayload)
	}
	return c.Redirect(301, "/login")
}

func getLogin(c echo.Context) error {
	basePayload := models.NewBasePayload(getUserFromContext(c))
	return c.Render(200, "login", basePayload)
}

func postLogin(c echo.Context, userService services.UserService) error {
	basePayload := models.NewBasePayload(getUserFromContext(c))
	username := c.FormValue("username")
	password := c.FormValue("password")
	user, err := userService.AuthenticateUser(username, password)
	if err == models.ErrInvalidLogin {
		basePayload.AddError(err)
		return c.Render(200, "login", basePayload)
	}
	if err != nil {
		basePayload.AddError(err)
		return c.Render(200, "login", basePayload)
	}
	sess, err := middleware.GetSession("session", c)
	if err != nil {
		basePayload.AddError(err)
		return c.Render(200, "login", basePayload)
	}
	sess.Values["userID"] = user.ID
	if err = sess.Save(c.Request(), c.Response()); err != nil {
		basePayload.AddError(err)
		return c.Render(200, "login", basePayload)
	}
	return c.Redirect(301, "/")
}

func getUser(c echo.Context, userService services.UserService) error {
	userIDString := c.Param("id")
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	contextUser := getUserFromContext(c)
	basePayload := models.NewBasePayload(contextUser)
	if err != nil {
		return c.Redirect(301, "/")
	}
	queryUser, err := userService.GetUserByID(userID)
	if err != nil {
		basePayload.AddError(err)
		return c.Render(200, "login", basePayload)
	}
	if queryUser == nil || contextUser == nil {
		basePayload.AddError(models.ErrUserNotFound)
		return c.Render(200, "login", basePayload)
	}
	if queryUser.ID != contextUser.ID {
		basePayload.AddError(echo.ErrUnauthorized)
		return c.Render(200, "login", basePayload)
	}
	return c.Render(200, "user", basePayload)
}

func getLeaderboard(c echo.Context, userService services.UserService) error {
	basePayload := models.NewBasePayload(getUserFromContext(c))
	users, err := userService.GetAllUsers()
	if err != nil {
		c.Redirect(301, "/")
	}
	usersPayload := models.NewUsersPayload(users)
	payload := models.CombinePayloads(usersPayload, *basePayload)
	return c.Render(200, "leaderboard", payload)
}

func postLogout(c echo.Context) error {
	sess, _ := middleware.GetSession("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
	return c.Redirect(301, "/")
}

func getContinents(c echo.Context) error {
	basePayload := models.NewBasePayload(getUserFromContext(c))
	continents := models.GetAllContinents()
	allCountries := models.GetAllCountries()
	filter := c.FormValue("continent")
	countries := models.FilterCountriesByContinent(allCountries, filter)
	sortMethod := c.FormValue("sort-method")
	models.SortCountries(countries, sortMethod)
	continentsPayload := models.NewContinentPayload(continents, countries)
	payload := models.CombinePayloads(continentsPayload, *basePayload)
	return c.Render(200, "search_continents", payload)
}

func getGuessCountry(c echo.Context) error {
	user := getUserFromContext(c)
	basePayload := models.NewBasePayload(user)
	var answerCountry models.CountryData
	if user != nil {
		answerCountry = *models.GetCountryByName(user.CurrentCountry)
	} else {
		answerCountry = models.GetRandomCountry()
		cookie := middleware.SetCookie("answerCountryName", answerCountry.Name.CommonName)
		c.SetCookie(cookie)
	}
	countries := models.Countries
	var passed bool = false
	countriesPayload := models.NewCountriesPayload(countries, &answerCountry, nil, passed)
	payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_countries", payload)
}

func postGuessCountry(c echo.Context, userService services.UserService) error {
	user := getUserFromContext(c)
	basePayload := models.NewBasePayload(user)
	var answerCountry models.CountryData
	if user != nil {
		answerCountry = *models.GetCountryByName(user.CurrentCountry)
	} else {
		answerCountryCookie, err := c.Cookie("answerCountryName")
		if err != nil {
			if err != http.ErrNoCookie {
				c.Redirect(301, "/guess_countries")
				return err
			}
			return err
		}
		answerCountry = *models.GetCountryByName(answerCountryCookie.Value)
	}
	countries := models.Countries
	guessCountryName := c.FormValue("country-guess")

	var passed = false
	if guessCountryName == answerCountry.Name.CommonName {
		passed = true
		if user != nil {
			userService.UpdateCurrentCountry(user)
		}
	}

	if user != nil {
		userService.UpdateCountryScore(user.ID, passed)
	}
	guessCountry := models.GetCountryByName(guessCountryName)
	countriesPayload := models.NewCountriesPayload(countries, &answerCountry, guessCountry, passed)
	payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_countries", payload)
}

func getGuessCapital(c echo.Context) error {
	user := getUserFromContext(c)
	basePayload := models.NewBasePayload(user)
	// don't use countries where capital == null
	var answerCountry models.CountryData
	if user != nil {
		answerCountry = *models.GetCountryByName(user.CurrentCapital)
	} else {
		for len(answerCountry.Capitals) < 1 {
			answerCountry = models.GetRandomCountry()
		}
		cookie := middleware.SetCookie("answerCapitalName", answerCountry.Name.CommonName)
		c.SetCookie(cookie)
	}
	countries := models.Countries
	var passed bool = false
	countriesPayload := models.NewCountriesPayload(countries, &answerCountry, nil, passed)
	payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_capitals", payload)
}

func postGuessCapital(c echo.Context, userService services.UserService) error {
	user := getUserFromContext(c)
	basePayload := models.NewBasePayload(user)
	var answerCountry models.CountryData
	if user != nil {
		answerCountry = *models.GetCountryByName(user.CurrentCapital)
	} else {
		answerCountryCookie, err := c.Cookie("answerCapitalName")
		if err != nil {
			c.Redirect(301, "/guess_capitals")
		}
		answerCountry = *models.GetCountryByName(answerCountryCookie.Value)
	}
	countries := models.Countries
	guessCapital := c.FormValue("guess-capital")
	guessCountries := models.GetCountryByCapital(guessCapital)
	var passed bool = false
	for range guessCountries {
		for _, capital := range answerCountry.Capitals {
			if capital == guessCapital {
				passed = true
				if user != nil {
					err := userService.UpdateCurrentCapital(user)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
	var guessCountry *models.CountryData
	if len(guessCountries) < 1 {
		guessCountry = &models.CountryData{}
	} else {
		guessCountry = guessCountries[0]
	}
	if user != nil {
		userService.UpdateCapitalScore(user.ID, passed)
	}
	countriesPayload := models.NewCountriesPayload(countries, &answerCountry, guessCountry, passed)
	payload := models.CombinePayloads(countriesPayload, *basePayload)
	return c.Render(200, "guess_capitals", payload)
}

func postSkipCountry(c echo.Context, userService services.UserService) error {
	user := getUserFromContext(c)
	if user != nil {
		userService.UpdateCountryScore(user.ID, false)
		userService.UpdateCurrentCountry(user)
	}

	return c.Redirect(http.StatusSeeOther, "/guess_countries")
}

func postSkipCapital(c echo.Context, userService services.UserService) error {
	user := getUserFromContext(c)
	if user != nil {
		userService.UpdateCapitalScore(user.ID, false)
		userService.UpdateCurrentCapital(user)
	}

	return c.Redirect(http.StatusSeeOther, "/guess_capitals")
}
