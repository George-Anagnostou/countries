package models

import (
    "errors"
    "fmt"
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
    ID           int64
    Username     string
    Password     []byte
    CountryScore int64
    CapitalScore int64
    CreatedAt    time.Time
}

func (u *User) String() string {
    return fmt.Sprintf("%s\n\t%d\n\t%d\n\t%d\n\t%v",
        u.Username, u.ID, u.CountryScore, u.CapitalScore, u.CreatedAt)
}

// type UserStore interface {
//     AddUser(user *User) error
//     GetUserByUsername(username string) (*User, error)
// }

// func AddUser(username, password string) error {
//     hashedPassword, err := HashPassword(password)
//     if err != nil {
//         return err
//     }
//     err = db.InsertUser(username, hashedPassword)
//     if err != nil {
//         return err
//     }
//     return nil
// }

// func AuthenticateUser(username, password string) error {
//     storedPassword, err := db.GetHashedPassword(username)
//     if err != nil {
//         return ErrInvalidLogin
//     }
//     err = CheckPassword(storedPassword, password)
//     if err != nil {
//         return ErrInvalidLogin
//     }
//     return nil
// }

// func GetUser(username string) (*User, error) {
//     user, err := db.GetUser(username)
//     if err != nil {
//         return nil, err
//     }
//     return user, nil
// }

func (u *User) GetID() int64 {
    return u.ID
}

func (u *User) GetUsername() string {
    return u.Username
}

func (u *User) SetUsername(username string) {
    u.Username = username
}

func (u *User) GetCountryScore() int64 {
    return u.CountryScore
}

func (u *User) IncrementCountry(score int64) {
    u.CountryScore++
}

func (u *User) SetCountryScore(score int64) {
    u.CountryScore = score
}

func (u User) GetCapitalScore() int64 {
    return u.CapitalScore
}

func (u *User) IncrementCapital(score int64) {
    u.CapitalScore++
}

func (u *User) SetCapitalScore(score int64) {
    u.CapitalScore = score
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
