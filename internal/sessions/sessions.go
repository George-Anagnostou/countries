package sessions

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func InitSessionStore() *sessions.CookieStore {
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

    return sessions.NewCookieStore(keyBytes)
}

func GetSession(c echo.Context, sessionName string) (*sessions.Session, error) {
    return session.Get(sessionName, c)
}

func Middleware(store *sessions.CookieStore) echo.MiddlewareFunc {
    return session.Middleware(store)
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


