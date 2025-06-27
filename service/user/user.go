package usersvc

import (
	"ThreeLayerArch/models"
	"ThreeLayerArch/store/user"
	"errors"
)

type UserService struct {
	Store *userstore.UserStore
}

func (s *UserService) CreateUser(name string) (*models.User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	id, err := s.Store.Create(name)
	if err != nil {
		return nil, err
	}
	return &models.User{UserID: id, Name: name}, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Store.GetByID(id)
}
func (s *UserService) View_Users() ([]models.User, error) {

	return s.Store.ViewUsers()
}
