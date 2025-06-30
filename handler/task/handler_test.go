package task

import (
	Models "ThreeLayerArch/models"
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// MockService implements TaskServices
type MockService struct {
	tasks map[int]Models.Tasks
}

func (m *MockService) Add_Task(task string) (bool, error) {
	if task == "fail" {
		return false, errors.New("failed to add task")
	}
	return true, nil
}

func (m *MockService) View_Task() ([]Models.Tasks, error) {
	return []Models.Tasks{{Tid: 1, Task: "Test Task", Completed: false}}, nil
}

func (m *MockService) Get_By_ID(i int) (Models.Tasks, error) {
	if i == 1 {
		return Models.Tasks{Tid: 1, Task: "Test Task", Completed: false}, nil
	}
	return Models.Tasks{}, errors.New("not found")
}

func (m *MockService) Update_Task(i int) (bool, error) {
	if i == 1 {
		return true, nil
	}
	return false, errors.New("update error")
}

func (m *MockService) Delete_Task(i int) (bool, error) {
	if i == 1 {
		return true, nil
	}
	return false, errors.New("delete error")
}

func TestAddTask(t *testing.T) {
	service := &MockService{}
	handler := New(service)

	body := []byte(`{"task":"Test Task"}`)
	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Addtask(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", res.StatusCode)
	}
}

func TestViewTask(t *testing.T) {
	service := &MockService{}
	handler := New(service)

	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	w := httptest.NewRecorder()

	handler.Viewtask(w, req)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	if !bytes.Contains(body, []byte("Test Task")) {
		t.Errorf("expected task in response, got: %s", string(body))
	}
}

func TestGetTask(t *testing.T) {
	service := &MockService{}
	handler := New(service)

	req := httptest.NewRequest(http.MethodGet, "/task/1", nil)
	req.SetPathValue("id", strconv.Itoa(1))
	w := httptest.NewRecorder()

	handler.Gettask(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}

func TestUpdateTask(t *testing.T) {
	service := &MockService{}
	handler := New(service)

	req := httptest.NewRequest(http.MethodPut, "/task/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.Updatetask(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}

func TestDeleteTask(t *testing.T) {
	service := &MockService{}
	handler := New(service)

	req := httptest.NewRequest(http.MethodDelete, "/task/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.Deletetask(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}
