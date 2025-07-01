//package task
//
//import (
//	Models "ThreeLayerArch/models"
//	"bytes"
//	"errors"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"strconv"
//	"testing"
//)
//
//// MockService implements TaskServices
//type MockService struct {
//	tasks map[int]Models.Tasks
//}
//
//func (m *MockService) Add_Task(task string) (bool, error) {
//	if task == "fail" {
//		return false, errors.New("failed to add task")
//	}
//	return true, nil
//}
//
//func (m *MockService) View_Task() ([]Models.Tasks, error) {
//	return []Models.Tasks{{Tid: 1, Task: "Test Task", Completed: false}}, nil
//}
//
//func (m *MockService) Get_By_ID(i int) (Models.Tasks, error) {
//	if i == 1 {
//		return Models.Tasks{Tid: 1, Task: "Test Task", Completed: false}, nil
//	}
//	return Models.Tasks{}, errors.New("not found")
//}
//
//func (m *MockService) Update_Task(i int) (bool, error) {
//	if i == 1 {
//		return true, nil
//	}
//	return false, errors.New("update error")
//}
//
//func (m *MockService) Delete_Task(i int) (bool, error) {
//	if i == 1 {
//		return true, nil
//	}
//	return false, errors.New("delete error")
//}
//
//func TestAddTask(t *testing.T) {
//	service := &MockService{}
//	handler := New(service)
//
//	body := []byte(`{"task":"Test Task"}`)
//	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
//	w := httptest.NewRecorder()
//
//	handler.Addtask(w, req)
//
//	res := w.Result()
//	if res.StatusCode != http.StatusCreated {
//		t.Errorf("expected 201 Created, got %d", res.StatusCode)
//	}
//}
//
//func TestViewTask(t *testing.T) {
//	service := &MockService{}
//	handler := New(service)
//
//	req := httptest.NewRequest(http.MethodGet, "/task", nil)
//	w := httptest.NewRecorder()
//
//	handler.Viewtask(w, req)
//
//	res := w.Result()
//	body, _ := io.ReadAll(res.Body)
//	if !bytes.Contains(body, []byte("Test Task")) {
//		t.Errorf("expected task in response, got: %s", string(body))
//	}
//}
//
//func TestGetTask(t *testing.T) {
//	service := &MockService{}
//	handler := New(service)
//
//	req := httptest.NewRequest(http.MethodGet, "/task/1", nil)
//	req.SetPathValue("id", strconv.Itoa(1))
//	w := httptest.NewRecorder()
//
//	handler.Gettask(w, req)
//
//	res := w.Result()
//	if res.StatusCode != http.StatusOK {
//		t.Errorf("expected 200 OK, got %d", res.StatusCode)
//	}
//}
//
//func TestUpdateTask(t *testing.T) {
//	service := &MockService{}
//	handler := New(service)
//
//	req := httptest.NewRequest(http.MethodPut, "/task/1", nil)
//	req.SetPathValue("id", "1")
//	w := httptest.NewRecorder()
//
//	handler.Updatetask(w, req)
//
//	res := w.Result()
//	if res.StatusCode != http.StatusOK {
//		t.Errorf("expected 200 OK, got %d", res.StatusCode)
//	}
//}
//
//func TestDeleteTask(t *testing.T) {
//	service := &MockService{}
//	handler := New(service)
//
//	req := httptest.NewRequest(http.MethodDelete, "/task/1", nil)
//	req.SetPathValue("id", "1")
//	w := httptest.NewRecorder()
//
//	handler.Deletetask(w, req)
//
//	res := w.Result()
//	if res.StatusCode != http.StatusOK {
//		t.Errorf("expected 200 OK, got %d", res.StatusCode)
//	}
//}

package task_test

import (
	"ThreeLayerArch/handler/task"
	"ThreeLayerArch/models"
	"bytes"
	"errors"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		requestBody    string
		mockInput      string
		mockReturn     bool
		mockError      error
		expectedStatus int
		expectedBody   string
		skipMock       bool // For cases like invalid JSON where mock call isn't made
	}{
		{
			name:           "Success - task added",
			requestBody:    `{"task":"Test Task"}`,
			mockInput:      "Test Task",
			mockReturn:     true,
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
		{
			name:           "Error - Add_Task fails",
			requestBody:    `{"task":"Fail Task"}`,
			mockInput:      "Fail Task",
			mockReturn:     false,
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create task",
		},
		{
			name:           "Error - Invalid JSON body",
			requestBody:    `{invalid-json}`,
			skipMock:       true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := task.NewMockTaskService(ctrl)

			if !tt.skipMock {
				mockService.EXPECT().
					Add_Task(tt.mockInput).
					Return(tt.mockReturn, tt.mockError)
			}

			handler := task.New(mockService)

			req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte(tt.requestBody)))
			w := httptest.NewRecorder()

			handler.Addtask(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedBody != "" && !bytes.Contains(body, []byte(tt.expectedBody)) {
				t.Errorf("expected body to contain %q, got: %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestViewTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		mockReturn     []models.Tasks
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - returns task list",
			mockReturn: []models.Tasks{
				{Tid: 1, Task: "Test Task", Completed: false},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Test Task",
		},
		{
			name:           "Error - View_Task fails",
			mockReturn:     nil,
			mockError:      errors.New("db failure"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to retrieve tasks",
		},
		{
			name:           "Empty task list",
			mockReturn:     []models.Tasks{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "[]", // JSON array of 0 tasks
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := task.NewMockTaskService(ctrl)
			mockService.EXPECT().
				View_Task().
				Return(tt.mockReturn, tt.mockError)

			handler := task.New(mockService)

			req := httptest.NewRequest(http.MethodGet, "/task", nil)
			w := httptest.NewRecorder()

			handler.Viewtask(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if !bytes.Contains(body, []byte(tt.expectedBody)) {
				t.Errorf("expected body to contain %q, got: %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		urlParam       string
		mockID         int
		mockReturn     models.Tasks
		mockError      error
		expectedStatus int
		expectedBody   string
		expectMock     bool
	}{
		{
			name:           "Success - task found",
			urlParam:       "1",
			mockID:         1,
			mockReturn:     models.Tasks{Tid: 1, Task: "Test Task", Completed: false},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Test Task",
			expectMock:     true,
		},
		{
			name:           "Error - task not found",
			urlParam:       "99",
			mockID:         99,
			mockReturn:     models.Tasks{},
			mockError:      errors.New("not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "task not found",
			expectMock:     true,
		},
		{
			name:           "Error - invalid ID format",
			urlParam:       "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid task ID",
			expectMock:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := task.NewMockTaskService(ctrl)
			if tt.expectMock {
				mockService.EXPECT().
					Get_By_ID(tt.mockID).
					Return(tt.mockReturn, tt.mockError)
			}

			handler := task.New(mockService)

			req := httptest.NewRequest(http.MethodGet, "/task/"+tt.urlParam, nil)
			req.SetPathValue("id", tt.urlParam)
			w := httptest.NewRecorder()

			handler.Gettask(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			if !bytes.Contains(body, []byte(tt.expectedBody)) {
				t.Errorf("expected body to contain %q, got: %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		urlParam       string
		mockID         int
		mockReturn     bool
		mockError      error
		expectedStatus int
		expectedBody   string
		expectMock     bool
	}{
		{
			name:           "Success - task updated",
			urlParam:       "1",
			mockID:         1,
			mockReturn:     true,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectMock:     true,
		},
		{
			name:           "Error - update failed",
			urlParam:       "2",
			mockID:         2,
			mockReturn:     false,
			mockError:      errors.New("update error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to update task",
			expectMock:     true,
		},
		{
			name:           "Error - invalid ID",
			urlParam:       "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid task ID",
			expectMock:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := task.NewMockTaskService(ctrl)
			if tt.expectMock {
				mockService.EXPECT().
					Update_Task(tt.mockID).
					Return(tt.mockReturn, tt.mockError)
			}

			handler := task.New(mockService)

			req := httptest.NewRequest(http.MethodPut, "/task/"+tt.urlParam, nil)
			req.SetPathValue("id", tt.urlParam)
			w := httptest.NewRecorder()

			handler.Updatetask(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			if tt.expectedBody != "" && !bytes.Contains(body, []byte(tt.expectedBody)) {
				t.Errorf("expected body to contain %q, got: %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		urlParam       string
		mockID         int
		mockReturn     bool
		mockError      error
		expectedStatus int
		expectedBody   string
		expectMock     bool
	}{
		{
			name:           "Success - task deleted",
			urlParam:       "1",
			mockID:         1,
			mockReturn:     true,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectMock:     true,
		},
		{
			name:           "Error - deletion failed",
			urlParam:       "99",
			mockID:         99,
			mockReturn:     false,
			mockError:      errors.New("delete error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to delete task",
			expectMock:     true,
		},
		{
			name:           "Error - invalid ID format",
			urlParam:       "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid task ID",
			expectMock:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := task.NewMockTaskService(ctrl)

			if tt.expectMock {
				mockService.EXPECT().
					Delete_Task(tt.mockID).
					Return(tt.mockReturn, tt.mockError)
			}

			handler := task.New(mockService)

			req := httptest.NewRequest(http.MethodDelete, "/task/"+tt.urlParam, nil)
			req.SetPathValue("id", tt.urlParam)
			w := httptest.NewRecorder()

			handler.Deletetask(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			if tt.expectedBody != "" && !bytes.Contains(body, []byte(tt.expectedBody)) {
				t.Errorf("expected body to contain %q, got: %s", tt.expectedBody, string(body))
			}
		})
	}
}
