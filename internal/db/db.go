package db

import (
    "database/sql"

    "github.com/George-Anagnostou/countries/internal/models"

    _ "github.com/mattn/go-sqlite3"
)

type UserDB struct{}

func (udb *UserDB) GetUserByUsername(username string) (*models.User, error) {
    user := new(models.User)
    query := `
        SELECT
            id,
            username,
            password,
            country_score,
            capital_score,
    created_at
    FROM users
            WHERE username = ?
        ;`
    err := db.QueryRow(query, username).Scan(
            &user.ID,
            &user.Username,
            &user.Password,
            &user.CountryScore,
            &user.CapitalScore,
            &user.CreatedAt,
        )
    if err != nil {
        return nil, err
    }
    return user, nil
}

// DELETE
func (udb *UserDB) AddUser(user *models.User) error {
    query := "INSERT INTO users (username, password) VALUES (?, ?)"
    _, err := db.Exec(query, user.Username, user.Password)
    return err
}

var db, err = sql.Open("sqlite3", "./data/countries.db")

func AddUser(username, password string) error {
    hashedPassword, err := models.HashPassword(password)
    if err != nil {
        return err
    }
    query := "INSERT INTO users (username, password) VALUES (?, ?)"
    _, err = db.Exec(query, username, hashedPassword)
    return err
}

func AuthenticateUser(username, password string) (*models.User, error) {
    user, err := GetUserByUsername(username)
    if err != nil {
        return nil, models.ErrInvalidLogin
    }
    // storedPassword, err := GetHashedPassword(username)
    // if err != nil {
    //     return nil, models.ErrInvalidLogin
    // }
    err = models.CheckPassword(string(user.Password), password)
    if err != nil {
        return nil, models.ErrInvalidLogin
    }
    return user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
    user := new(models.User)
    query := `
        SELECT
            id,
            username,
            password,
            country_score,
            capital_score,
            created_at
    FROM users
            WHERE username = ?
        ;`
    err := db.QueryRow(query, username).Scan(
            &user.ID,
            &user.Username,
            &user.Password,
            &user.CountryScore,
            &user.CapitalScore,
            &user.CreatedAt,
        )
    if err != nil {
        return nil, err
    }
    return user, nil
}

func GetHashedPassword(username string) (string, error) {
    var hashedPassword string
    err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
    if err != nil {
        return "", err
    }
    return hashedPassword, nil
}
