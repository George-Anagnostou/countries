package models

import (
)

type User struct {
    ID           int64
    Username     string
    CountryScore int64
    CapitalScore int64
}

func NewUser(username string) *User {
    return &User{
        Username: username,
        CountryScore: 0,
        CapitalScore: 0,
    }
}

func (u *User) GetUserID() int64 {
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
