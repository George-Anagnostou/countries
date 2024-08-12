package services

import (
	"github.com/George-Anagnostou/countries/internal/models"
	"github.com/George-Anagnostou/countries/internal/repositories"
)

type UserService struct {
	UserRepo repositories.UserRepository
}

func (s *UserService) AddUser(username, password string) error {
	return s.UserRepo.AddUser(username, password)
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	return s.UserRepo.AuthenticateUser(username, password)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.UserRepo.GetUserByUsername(username)
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	return s.UserRepo.GetUserByID(id)
}

func (s *UserService) UpdateCountryScore(userID int64, correct bool) {
	s.UserRepo.UpdateCountryScore(userID, correct)
}

func (s *UserService) UpdateCapitalScore(userID int64, correct bool) {
	s.UserRepo.UpdateCapitalScore(userID, correct)
}

func (s *UserService) UpdateCurrentCountry(user *models.User) error {
	return s.UserRepo.UpdateCurrentCountry(user)
}

func (s *UserService) UpdateCurrentCapital(user *models.User) error {
	return s.UserRepo.UpdateCurrentCapital(user)
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.UserRepo.GetAllUsers()
}

// GetHashedPassword(username string) (string, error)
