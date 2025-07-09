package user

import (
	"ThreeLayerArch/models"
	"errors"
	"gofr.dev/pkg/gofr"
)

type Service struct {
	Store UserStore
}

func (s *Service) CreateUser(ctx *gofr.Context, name string) (*models.User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	id, err := s.Store.Create(ctx, name)
	if err != nil {
		return nil, err
	}
	return &models.User{UserID: id, Name: name}, nil
}

func (s *Service) GetUserByID(ctx *gofr.Context, id int) (*models.User, error) {
	return s.Store.GetByID(ctx, id)
}
func (s *Service) View_Users(ctx *gofr.Context) ([]models.User, error) {

	return s.Store.ViewUsers(ctx)
}
