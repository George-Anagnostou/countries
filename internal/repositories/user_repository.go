package repositories

import (
	"errors"

	"github.com/George-Anagnostou/countries/internal/models"
)

var (
	ErrInvalidRegistration = errors.New("invalid registration")
)

type UserRepository interface {
	Initialize() error
	AddUser(username, password string) error
	AuthenticateUser(username, password string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetHashedPassword(username string) (string, error)
	UpdateCountryScore(userID int64, correct bool) error
	UpdateCapitalScore(userID int64, correct bool) error
	UpdateCurrentCountry(user *models.User) error
	UpdateCurrentCapital(user *models.User) error
}
