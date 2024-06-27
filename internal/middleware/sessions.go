package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/George-Anagnostou/countries/internal/db"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func InitSessionStore() echo.MiddlewareFunc {
    err := godotenv.Load()
    if err != nil {
        log.Println("no .env file found")
    }

    sessionKey := os.Getenv("SESSION_KEY")
    if sessionKey == "" {
        log.Fatal("SESSION_KEY not set")
    }

    keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
    if err != nil {
        log.Fatal(err)
    }

    store := sessions.NewCookieStore(keyBytes)
    return session.Middleware(store)
}

func GetSession(c echo.Context, sessionName string) (*sessions.Session, error) {
    return session.Get(sessionName, c)
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        sess, _ := session.Get("session", c)
        userID := sess.Values["user_id"]

        if userID == nil {
            c.Set("userID", nil)
        } else {
            c.Set("userID", userID)
        }

        return next(c)
    }
}

func SetCookie(key string, name string) *http.Cookie {
    cookie := new(http.Cookie)
    cookie.Name = key
    cookie.Value = name
    cookie.Expires = time.Now().Add(5 * time.Minute)
    return cookie
}

func ResetCookie(key string) *http.Cookie {
    cookie := new(http.Cookie)
    cookie.Name = key
    cookie.Expires = time.Now().Add(-time.Hour)
    return cookie
}

// func IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
//     return func(c echo.Context) error {
//         session := sessions.
//         userID := session.Get("user_id").(int64)
//         user, _ := db.GetUserByID(userID)
//         c.Set("user", user)
//         return next(c)
//     }
// }
