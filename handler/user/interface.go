package user

import (
	"ThreeLayerArch/models"
	"gofr.dev/pkg/gofr"
)

type UserService interface {
	CreateUser(ctx *gofr.Context, name string) (*models.User, error)
	GetUserByID(ctx *gofr.Context, id int) (*models.User, error)
	View_Users(ctx *gofr.Context) ([]models.User, error)
}
