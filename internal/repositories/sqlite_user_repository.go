package repositories

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/George-Anagnostou/countries/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserRepository struct {
    DB *sql.DB
}

func NewSQLiteUserRepository() *SQLiteUserRepository {
    databaseFile := os.Getenv("SQLITE_FILE")
    db, err := sql.Open("sqlite3", databaseFile)
    if err != nil {
        panic(err)
    }
    return &SQLiteUserRepository{
        DB: db,
    }
}

func (r *SQLiteUserRepository) AddUser(username, password string) error {
    hashedPassword, err := models.HashPassword(password)
    if err != nil {
        return err
    }
    countryJson, err := json.Marshal(models.GetRandomCountry())
    if err != nil {
        return err
    }
    capitalJson, err := json.Marshal(models.GetRandomCountry())
    if err != nil {
        return err
    }
    stmt := "INSERT INTO users (username, password, current_country, current_capital) VALUES (?, ?, ?, ?)"
    result, err := r.DB.Exec(
        stmt,
        username,
        hashedPassword,
        string(countryJson),
        string(capitalJson),
    )
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

func (r *SQLiteUserRepository) AuthenticateUser(username, password string) (*models.User, error) {
    user, err := r.GetUserByUsername(username)
    if err != nil {
        return nil, models.ErrInvalidLogin
    }
    err = models.CheckPassword(string(user.Password), password)
    if err != nil {
        return nil, models.ErrInvalidLogin
    }
    return user, nil
}

func (r *SQLiteUserRepository) GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    var countryJSON string
    var capitalJSON string
    query := `
        SELECT
            id,
            username,
            password,
            longest_country_score,
            current_country_score,
            longest_capital_score,
            current_capital_score,
            current_country,
            current_capital,
            created_at
        FROM users
            WHERE username = ?
        ;`
    err := r.DB.QueryRow(query, username).Scan(
        &user.ID,
        &user.Username,
        &user.Password,
        &user.LongestCountryScore,
        &user.CurrentCountryScore,
        &user.LongestCapitalScore,
        &user.CurrentCapitalScore,
        &countryJSON,
        &capitalJSON,
        &user.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal([]byte(countryJSON), &user.CurrentCountry)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal([]byte(capitalJSON), &user.CurrentCapital)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *SQLiteUserRepository) GetUserByID(id int64) (*models.User, error) {
    var user models.User
    var countryJSON string
    var capitalJSON string
    query := `
        SELECT
            id,
            username,
            longest_country_score,
            current_country_score,
            longest_capital_score,
            current_capital_score,
            current_country,
            current_capital,
            created_at
        FROM users
            WHERE id = ?
        ;`
    err := r.DB.QueryRow(query, id).Scan(
        &user.ID,
        &user.Username,
        &user.LongestCountryScore,
        &user.CurrentCountryScore,
        &user.LongestCapitalScore,
        &user.CurrentCapitalScore,
        &countryJSON,
        &capitalJSON,
        &user.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal([]byte(countryJSON), &user.CurrentCountry)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal([]byte(capitalJSON), &user.CurrentCapital)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *SQLiteUserRepository) GetAllUsers() ([]*models.User, error) {
    query := `
        SELECT
            id,
            username,
            longest_country_score,
            current_country_score,
            longest_capital_score,
            current_capital_score,
            current_country,
            current_capital,
            created_at
        FROM users
        ;`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var users []*models.User
    for rows.Next() {
        var user models.User
        var countryJSON string
        var capitalJSON string
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.LongestCountryScore,
            &user.CurrentCountryScore,
            &user.LongestCapitalScore,
            &user.CurrentCapitalScore,
            &countryJSON,
            &capitalJSON,
            &user.CreatedAt,
        )
        if err != nil {
            return users, err
        }
        err = json.Unmarshal([]byte(countryJSON), &user.CurrentCountry)
        if err != nil {
            return nil, err
        }
        err = json.Unmarshal([]byte(capitalJSON), &user.CurrentCapital)
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    return users, nil
}

func (r *SQLiteUserRepository) GetHashedPassword(username string) (string, error) {
    var hashedPassword string
    err := r.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
    if err != nil {
        return "", err
    }
    return hashedPassword, nil
}

func (r *SQLiteUserRepository) UpdateCountryScore(userID int64, correct bool) error {
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
    row := r.DB.QueryRow(query, userID)
    err := row.Scan(&currentStreak, &longestStreak)
    if err != nil {
        return err
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
    _, err = r.DB.Exec(stmt, currentStreak, longestStreak, userID)
    if err != nil {
        return err
    }
    return nil
}

func (r *SQLiteUserRepository) UpdateCapitalScore(userID int64, correct bool) error {
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
    row := r.DB.QueryRow(query, userID)
    err := row.Scan(&currentStreak, &longestStreak)
    if err != nil {
        return err
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
    _, err = r.DB.Exec(stmt, currentStreak, longestStreak, userID)
    if err != nil {
        return err
    }
    return nil
}

func (r *SQLiteUserRepository) UpdateCurrentCountry(user *models.User) error {
    nextCountryJSON, err := json.Marshal(models.GetRandomCountry())
    if err != nil {
        return err
    }
    stmt := `
            UPDATE
                users
            SET
                current_country = ?
            WHERE id = ?
        ;`
    _, err = r.DB.Exec(stmt, string(nextCountryJSON), user.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *SQLiteUserRepository) UpdateCurrentCapital(user *models.User) error {
    nextCountryJSON, err := json.Marshal(models.GetRandomCountry())
    stmt := `
            UPDATE
                users
            SET
                current_capital  = ?
            WHERE id = ?
        ;`
    _, err = r.DB.Exec(stmt, string(nextCountryJSON), user.ID)
    if err != nil {
        return err
    }
    return nil
}
