package task

import "ThreeLayerArch/models"

type TaskService interface {
	Add_Task(task string) (bool, error)
	View_Task() ([]models.Tasks, error)
	Get_By_ID(i int) (models.Tasks, error)
	Update_Task(i int) (bool, error)
	Delete_Task(i int) (bool, error)
}
