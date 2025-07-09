package task

import (
	"ThreeLayerArch/models"
	"gofr.dev/pkg/gofr"
)

type TaskService interface {
	Add_Task(ctx *gofr.Context, task string) (bool, error)
	View_Task(ctx *gofr.Context) ([]models.Tasks, error)
	Get_By_ID(ctx *gofr.Context, i int) (models.Tasks, error)
	Update_Task(ctx *gofr.Context, i int) (bool, error)
	Delete_Task(ctx *gofr.Context, i int) (bool, error)
}
