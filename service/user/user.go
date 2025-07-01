package user

import (
	"ThreeLayerArch/models"
	"errors"
)

type UserStore interface {
	Create(string) (int, error)
	GetByID(int) (*models.User, error)
	ViewUsers() ([]models.User, error)
}

type Service struct {
	Store UserStore
}

func (s *Service) CreateUser(name string) (*models.User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	id, err := s.Store.Create(name)
	if err != nil {
		return nil, err
	}
	return &models.User{UserID: id, Name: name}, nil
}

func (s *Service) GetUserByID(id int) (*models.User, error) {
	return s.Store.GetByID(id)
}
func (s *Service) View_Users() ([]models.User, error) {

	return s.Store.ViewUsers()
}
