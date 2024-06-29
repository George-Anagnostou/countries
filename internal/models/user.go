package models

import (
    "errors"
    // "fmt"
    "time"

    // "github.com/George-Anagnostou/countries/internal/db"
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
    CreatedAt           time.Time
}

// func (u *User) String() string {
//     return fmt.Sprintf("username = %s\nid = %d\ncountryScore = %d\ncapitalScore = %d\ncreatedAt = %v",
//         u.Username, u.ID, u.CountryScore, u.CapitalScore, u.CreatedAt)
// }

// type UserStore interface {
//     AddUser(user *User) error
//     GetUserByUsername(username string) (*User, error)
// }

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
