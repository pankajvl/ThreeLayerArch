package user_test

import (
	"ThreeLayerArch/handler/user"
	"ThreeLayerArch/models"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserService struct {
	CreateUserFunc  func(name string) (*models.User, error)
	GetUserByIDFunc func(id int) (*models.User, error)
	ViewUsersFunc   func() ([]models.User, error)
}

func (m *mockUserService) CreateUser(name string) (*models.User, error) {
	return m.CreateUserFunc(name)
}

func (m *mockUserService) GetUserByID(id int) (*models.User, error) {
	return m.GetUserByIDFunc(id)
}

func (m *mockUserService) View_Users() ([]models.User, error) {
	return m.ViewUsersFunc()
}

func TestCreateUser_Success(t *testing.T) {
	service := &mockUserService{
		CreateUserFunc: func(name string) (*models.User, error) {
			return &models.User{UserID: 1, Name: name}, nil
		},
	}

	handler := user.UserHandler{Service: service}

	body := []byte(`{"name": "John"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestCreateUser_InvalidBody(t *testing.T) {
	service := &mockUserService{}
	handler := user.UserHandler{Service: service}

	body := []byte(`invalid-json`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	service := &mockUserService{
		GetUserByIDFunc: func(id int) (*models.User, error) {
			return &models.User{UserID: id, Name: "John"}, nil
		},
	}

	handler := user.UserHandler{Service: service}

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	handler.GetUserByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	service := &mockUserService{
		GetUserByIDFunc: func(id int) (*models.User, error) {
			return nil, nil
		},
	}

	handler := user.UserHandler{Service: service}

	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	req.SetPathValue("id", "99")
	rec := httptest.NewRecorder()

	handler.GetUserByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestViewUsers_Success(t *testing.T) {
	service := &mockUserService{
		ViewUsersFunc: func() ([]models.User, error) {
			return []models.User{
				{UserID: 1, Name: "Alice"},
				{UserID: 2, Name: "Bob"},
			}, nil
		},
	}

	handler := user.UserHandler{Service: service}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	handler.ViewUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if body == "" || !contains(body, "Alice") || !contains(body, "Bob") {
		t.Errorf("unexpected response body: %s", body)
	}
}

func TestViewUsers_Error(t *testing.T) {
	service := &mockUserService{
		ViewUsersFunc: func() ([]models.User, error) {
			return nil, errors.New("db error")
		},
	}

	handler := user.UserHandler{Service: service}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	handler.ViewUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 even on log error, got %d", rec.Code)
	}
}

// Utility function
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
