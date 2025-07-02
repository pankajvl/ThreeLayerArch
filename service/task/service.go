package task

import (
	Model "ThreeLayerArch/models"
	"errors"
	"gofr.dev/pkg/gofr"
	"log"
)

type Service struct {
	store TaskStore
}

func New(store TaskStore) *Service {
	return &Service{store: store}
}

func (s *Service) Add_Task(ctx *gofr.Context, task string) (bool, error) {
	if task == "" {
		return false, errors.New("task cannot be empty")
	}
	return s.store.AddTask(ctx, task)
}

func (s *Service) View_Task(ctx *gofr.Context) ([]Model.Tasks, error) {

	return s.store.ViewTask(ctx)
}

func (s *Service) Get_By_ID(ctx *gofr.Context, i int) (Model.Tasks, error) {

	if s.store.CheckIfExists(ctx, i) {
		ans, err := s.store.GetByID(ctx, i)
		if err != nil {
			log.Printf("Error in SERVICES.GetByID: %v", err)
			return Model.Tasks{}, err
		}
		return ans, nil
	}
	return Model.Tasks{}, errors.New("ID not found")
}

func (s *Service) Update_Task(ctx *gofr.Context, i int) (bool, error) {
	if s.store.CheckIfExists(ctx, i) {
		ans, err := s.store.UpdateTask(ctx, i)
		if err != nil {
			log.Printf("Error in SERVICES.UpdateTask: %v", err)
			return false, err
		}
		return ans, nil
	}
	return false, errors.New("ID not found")
}

func (s *Service) Delete_Task(ctx *gofr.Context, i int) (bool, error) {
	if s.store.CheckIfExists(ctx, i) {
		ans, err := s.store.DeleteTask(ctx, i)
		if err != nil {
			log.Printf("Error in SERVICES.DeleteTask: %v", err)
			return false, err
		}
		return ans, nil
	}
	return false, errors.New("ID not found")
}
