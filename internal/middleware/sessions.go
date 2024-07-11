package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/George-Anagnostou/countries/internal/db"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func InitSessionStore() echo.MiddlewareFunc {
    // load .env file from root
    err := godotenv.Load()
    if err != nil {
        log.Println("no .env file found")
    }

    sessionKey := os.Getenv("SESSION_KEY")
    if sessionKey == "" {
        log.Fatal("session key not set")
    }

    keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
    if err != nil {
        log.Fatal(err)
    }

    store := sessions.NewCookieStore(keyBytes)
    return session.Middleware(store)
}

func GetSession(sessionName string, c echo.Context) (*sessions.Session, error) {
    sess, err := session.Get(sessionName, c)
    sess.Options = &sessions.Options{
        Path: "/",
        MaxAge: 3600,
    }
    if err != nil {
        return nil, err
    }
    return sess, nil
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        sess, err := session.Get("session", c)
        if err != nil {
            log.Printf("error getting session: %v", err)
        }
        userID, ok := sess.Values["userID"].(int64)
        if !ok {
            log.Printf("error getting user session: %v", err)
        }
        user, err := db.GetUserByID(userID)
        if err != nil {
            log.Printf("error getting user: %v", err)
        }
        c.Set("user", user)
        return next(c)
    }
}

func SetCookie(key string, value string) *http.Cookie {
    cookie := new(http.Cookie)
    cookie.Name = key
    cookie.Value = value
    cookie.Expires = time.Now().Add(5 * time.Minute)
    return cookie
}

func ResetCookie(key string) *http.Cookie {
    cookie := new(http.Cookie)
    cookie.Name = key
    cookie.Expires = time.Now().Add(-time.Hour)
    return cookie
}
