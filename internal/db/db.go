package db

import (
    "database/sql"
    "errors"

    "github.com/George-Anagnostou/countries/internal/models"

    _ "github.com/mattn/go-sqlite3"
)

var db, err = sql.Open("sqlite3", "./data/countries.db")

var (
	ErrInvalidRegistration = errors.New("invalid registration")
)

func AddUser(username, password string) error {
    hashedPassword, err := models.HashPassword(password)
    if err != nil {
        return err
    }
    query := "INSERT INTO users (username, password) VALUES (?, ?)"
    result, err := db.Exec(query, username, hashedPassword)
    if err != nil {
        return ErrInvalidRegistration
    }
    rows, err := result.RowsAffected()
    if err != nil {
        return ErrInvalidRegistration
    }
    if rows != 1 {
        return ErrInvalidRegistration
    }
    return err
}

func AuthenticateUser(username, password string) (*models.User, error) {
    user, err := GetUserByUsername(username)
    if err != nil {
        return nil, models.ErrInvalidLogin
    }
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

func GetUserByID(id int) (*models.User, error) {
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
            WHERE id = ?
        ;`
    err := db.QueryRow(query, id).Scan(
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
