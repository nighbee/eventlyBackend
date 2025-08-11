package auth

import (
	"errors"

	model "github.com/nighbee/evently/internal/model"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

// Register creates a new user with the provided credentials.
func (s *Service) Register(name, email, password string) error {
	_, err := s.repo.GetUserByEmail(email)
	if err == nil {
		return errors.New("email already in use")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
	}
	return s.repo.CreateUser(user)
}

// Login validates credentials and returns the authenticated user when valid.
func (s *Service) Login(email, password string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
