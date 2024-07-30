package models

import (
    "errors"
    "time"

    "golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
	ErrUsernameTaken = errors.New("username taken")
)

type User struct {
    ID                  int64
    Username            string
    Password            []byte
    LongestCountryScore int64
    CurrentCountryScore int64
    LongestCapitalScore int64
    CurrentCapitalScore int64
    CurrentCountry      string
    CurrentCapital      string
    CreatedAt           time.Time
}

// hashes plaintext password using bcrypt
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// checks if plaintext password matches hashed password
func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
