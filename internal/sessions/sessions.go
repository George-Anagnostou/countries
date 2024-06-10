package sessions

import (
    "net/http"
    "time"
)

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
