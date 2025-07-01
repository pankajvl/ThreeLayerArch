package user

import "ThreeLayerArch/models"

type UserService interface {
	CreateUser(name string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	View_Users() ([]models.User, error)
}
