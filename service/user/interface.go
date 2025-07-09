package user

import (
	"ThreeLayerArch/models"
	"gofr.dev/pkg/gofr"
)

type UserStore interface {
	Create(ctx *gofr.Context, name string) (int, error)
	GetByID(ctx *gofr.Context, id int) (*models.User, error)
	ViewUsers(ctx *gofr.Context) ([]models.User, error)
}
