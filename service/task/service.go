package task

import (
	Model "ThreeLayerArch/models"
	"errors"
	"log"
)

type Service struct {
	store TaskStore
}

func New(store TaskStore) *Service {
	return &Service{store: store}
}

func (s *Service) Add_Task(task string) (bool, error) {
	if task == "" {
		return false, errors.New("task cannot be empty")
	}
	return s.store.AddTask(task)
}

func (s *Service) View_Task() ([]Model.Tasks, error) {

	return s.store.ViewTask()
}

func (s *Service) Get_By_ID(i int) (Model.Tasks, error) {

	if s.store.CheckIfExists(i) {
		ans, err := s.store.GetByID(i)
		if err != nil {
			log.Printf("Error in SERVICES.GetByID: %v", err)
			return Model.Tasks{}, err
		}
		return ans, nil
	}
	return Model.Tasks{}, errors.New("ID not found")
}

func (s *Service) Update_Task(i int) (bool, error) {
	if s.store.CheckIfExists(i) {
		ans, err := s.store.UpdateTask(i)
		if err != nil {
			log.Printf("Error in SERVICES.UpdateTask: %v", err)
			return false, err
		}
		return ans, nil
	}
	return false, errors.New("ID not found")
}

func (s *Service) Delete_Task(i int) (bool, error) {
	if s.store.CheckIfExists(i) {
		ans, err := s.store.DeleteTask(i)
		if err != nil {
			log.Printf("Error in SERVICES.DeleteTask: %v", err)
			return false, err
		}
		return ans, nil
	}
	return false, errors.New("ID not found")
}
