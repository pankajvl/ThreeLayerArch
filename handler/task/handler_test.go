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
	task "ThreeLayerArch/handler/task"
	"ThreeLayerArch/models"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type gofrResponse struct {
	result any
	err    error
}

func newMockCtx(t *testing.T) *gofr.Context {
	mockContainer, _ := container.NewMockContainer(t)
	return &gofr.Context{
		Context:   context.Background(),
		Request:   nil,
		Container: mockContainer,
	}
}

func TestAddtask(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		expectErr  error
		expectOK   bool
		mockNeeded bool
	}{
		{"valid input", `{"task":"Learn Go"}`, nil, true, true},
		{"invalid JSON", `{"task":123}`, errors.New("invalid character"), false, false},
		{"service error", `{"task":"Learn Go"}`, errors.New("service error"), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := task.NewMockTaskService(ctrl)
			h := task.New(mockSvc)
			ctx := newMockCtx(t)

			req := httptest.NewRequest(http.MethodPost, "/task", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = gofrHttp.NewRequest(req)

			if tt.mockNeeded {
				mockSvc.EXPECT().Add_Task(ctx, "Learn Go").Return(tt.expectOK, tt.expectErr)
			}

			res, err := h.Addtask(ctx)
			if tt.expectErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, "inserted successfully", res)
			}
		})
	}
}

func TestViewtask(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		err    error
	}{
		{"success", []models.Tasks{{Tid: 1, Task: "Go", Completed: false}}, nil},
		{"error case", nil, errors.New("db error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := task.NewMockTaskService(ctrl)
			h := task.New(mockSvc)
			ctx := newMockCtx(t)

			mockSvc.EXPECT().View_Task(ctx).Return(tt.expect, tt.err)

			res, err := h.Viewtask(ctx)
			assert.Equal(t, tt.expect, res)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGettask(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		expect any
		err    error
		mock   bool
	}{
		{"valid ID", "1", models.Tasks{Tid: 1, Task: "Go"}, nil, true},
		{"invalid ID", "abc", nil, strconv.ErrSyntax, false},
		{"service error", "1", nil, errors.New("not found"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := task.NewMockTaskService(ctrl)
			h := task.New(mockSvc)
			ctx := newMockCtx(t)

			req := httptest.NewRequest(http.MethodGet, "/task/{id}", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			ctx.Request = gofrHttp.NewRequest(req)

			id, err := strconv.Atoi(tt.id)
			if tt.mock && err == nil {
				mockSvc.EXPECT().Get_By_ID(ctx, id).Return(models.Tasks{}, errors.New("some error"))
			}

			res, err := h.Gettask(ctx)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expect, res)
			}
		})
	}
}

func TestUpdatetask(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		success bool
		err     error
		mock    bool
	}{
		{"valid", "1", true, nil, true},
		{"invalid id", "xyz", false, strconv.ErrSyntax, false},
		{"service error", "1", false, errors.New("fail"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := task.NewMockTaskService(ctrl)
			h := task.New(mockSvc)
			ctx := newMockCtx(t)

			req := httptest.NewRequest(http.MethodPut, "/task/{id}", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			ctx.Request = gofrHttp.NewRequest(req)

			id, err := strconv.Atoi(tt.id)
			if tt.mock && err == nil {
				mockSvc.EXPECT().Update_Task(ctx, id).Return(tt.success, tt.err)
			}

			res, err := h.Updatetask(ctx)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, "updated successfully", res)
			}
		})
	}
}

func TestDeletetask(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		success bool
		err     error
		mock    bool
	}{
		{"valid", "1", true, nil, true},
		{"invalid id", "abc", false, strconv.ErrSyntax, false},
		{"service error", "1", false, errors.New("fail"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := task.NewMockTaskService(ctrl)
			h := task.New(mockSvc)
			ctx := newMockCtx(t)

			req := httptest.NewRequest(http.MethodDelete, "/task/{id}", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			ctx.Request = gofrHttp.NewRequest(req)

			id, err := strconv.Atoi(tt.id)
			if tt.mock && err == nil {
				mockSvc.EXPECT().Delete_Task(ctx, id).Return(tt.success, tt.err)
			}

			res, err := h.Deletetask(ctx)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, "deleted successfully", res)
			}
		})
	}
}
