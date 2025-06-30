package task

import (
	Model "ThreeLayerArch/models"
	"errors"
	"testing"
)

// mockStore implements TaskStore for testing purposes
type mockStore struct{}

func (m *mockStore) AddTask(task string) (bool, error) {
	if task == "fail" {
		return false, errors.New("failed to add task")
	}
	return true, nil
}

func (m *mockStore) ViewTask() ([]Model.Tasks, error) {
	return []Model.Tasks{{Tid: 1, Task: "Test Task", Completed: false}}, nil
}

func (m *mockStore) GetByID(id int) (Model.Tasks, error) {
	if id == 1 {
		return Model.Tasks{Tid: 1, Task: "Test Task", Completed: false}, nil
	}
	return Model.Tasks{}, errors.New("not found")
}

func (m *mockStore) UpdateTask(id int) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, errors.New("update failed")
}

func (m *mockStore) DeleteTask(id int) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, errors.New("delete failed")
}

func (m *mockStore) CheckIfExists(i int) bool {
	return i == 1
}

func TestAddTask(t *testing.T) {
	s := New(&mockStore{})

	// Test success
	ok, err := s.Add_Task("New Task")
	if !ok || err != nil {
		t.Errorf("Expected success, got err: %v", err)
	}

	// Test failure
	ok, err = s.Add_Task("")
	if err == nil {
		t.Errorf("Expected error for empty task")
	}
}

func TestViewTask(t *testing.T) {
	s := New(&mockStore{})

	tasks, err := s.View_Task()
	if err != nil || len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %v, err: %v", len(tasks), err)
	}
}

func TestGetByID(t *testing.T) {
	s := New(&mockStore{})

	// Exists
	task, err := s.Get_By_ID(1)
	if err != nil || task.Tid != 1 {
		t.Errorf("Expected valid task, got: %v, err: %v", task, err)
	}

	// Does not exist
	_, err = s.Get_By_ID(99)
	if err == nil {
		t.Errorf("Expected error for missing ID")
	}
}

func TestUpdateTask(t *testing.T) {
	s := New(&mockStore{})

	// Exists
	ok, err := s.Update_Task(1)
	if !ok || err != nil {
		t.Errorf("Expected success, got err: %v", err)
	}

	// Does not exist
	ok, err = s.Update_Task(99)
	if err == nil {
		t.Errorf("Expected error for invalid ID")
	}
}

func TestDeleteTask(t *testing.T) {
	s := New(&mockStore{})

	// Exists
	ok, err := s.Delete_Task(1)
	if !ok || err != nil {
		t.Errorf("Expected success, got err: %v", err)
	}

	// Does not exist
	ok, err = s.Delete_Task(99)
	if err == nil {
		t.Errorf("Expected error for invalid ID")
	}
}
