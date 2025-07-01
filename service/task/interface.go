package task

import Model "ThreeLayerArch/models"

type TaskStore interface {
	AddTask(task string) (bool, error)
	ViewTask() ([]Model.Tasks, error)
	GetByID(id int) (Model.Tasks, error)
	UpdateTask(id int) (bool, error)
	DeleteTask(id int) (bool, error)
	CheckIfExists(i int) bool
}
