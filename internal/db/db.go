package db

import (
    "database/sql"
    "errors"
    "log"

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
    stmt := "INSERT INTO users (username, password) VALUES (?, ?)"
    result, err := db.Exec(stmt, username, hashedPassword)
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
            longest_country_score,
            current_country_score,
            longest_capital_score,
            current_capital_score,
            created_at
        FROM users
            WHERE username = ?
        ;`
    err := db.QueryRow(query, username).Scan(
            &user.ID,
            &user.Username,
            &user.Password,
            &user.LongestCountryScore,
            &user.CurrentCountryScore,
            &user.LongestCapitalScore,
            &user.CurrentCapitalScore,
            &user.CreatedAt,
        )
    if err != nil {
        return nil, err
    }
    return user, nil
}

func GetUserByID(id int64) (*models.User, error) {
    user := new(models.User)
    query := `
        SELECT
            id,
            username,
            password,
            longest_country_score,
            current_country_score,
            longest_capital_score,
            current_capital_score,
            created_at
    FROM users
            WHERE id = ?
        ;`
    err := db.QueryRow(query, id).Scan(
            &user.ID,
            &user.Username,
            &user.Password,
            &user.LongestCountryScore,
            &user.CurrentCountryScore,
            &user.LongestCapitalScore,
            &user.CurrentCapitalScore,
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

func UpdateCountryScore(userID int64, correct bool) {
    var currentStreak, longestStreak int64
    query := `
        SELECT
            current_country_score,
            longest_country_score
        FROM
            users
        WHERE
            id = ?
        ;`
    row := db.QueryRow(query, userID)
    err := row.Scan(&currentStreak, &longestStreak)
    if err != nil {
        log.Println("error fetching user streaks:", err)
        return
    }
    if correct {
        currentStreak++
        if currentStreak > longestStreak {
            longestStreak = currentStreak
        }
    } else {
        currentStreak = 0
    }
    stmt := `
            UPDATE
                users
            SET
                current_country_score = ?,
                longest_country_score = ?
            WHERE id = ?
        ;`
    _, err = db.Exec(stmt, currentStreak, longestStreak, userID)
    if err != nil {
        log.Println("error logging answer:", err)
    }
}

func UpdateCapitalScore(userID int64, correct bool) {
    var currentStreak, longestStreak int64
    query := `
        SELECT
            current_capital_score,
            longest_capital_score
        FROM
            users
        WHERE
            id = ?
        ;`
    row := db.QueryRow(query, userID)
    err := row.Scan(&currentStreak, &longestStreak)
    if err != nil {
        log.Println("error fetching user streaks:", err)
        return
    }
    if correct {
        currentStreak++
        if currentStreak > longestStreak {
            longestStreak = currentStreak
        }
    } else {
        currentStreak = 0
    }
    stmt := `
            UPDATE
                users
            SET
                current_capital_score = ?,
                longest_capital_score = ?
            WHERE id = ?
        ;`
    _, err = db.Exec(stmt, currentStreak, longestStreak, userID)
    if err != nil {
        log.Println("error logging answer:", err)
    }
}
