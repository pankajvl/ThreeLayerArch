package user

import (
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

func setupTestContext(t *testing.T, method, path, body string, pathVars map[string]string) *gofr.Context {
	mockContainer, _ := container.NewMockContainer(t)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if pathVars != nil {
		req = mux.SetURLVars(req, pathVars)
	}
	ctx := &gofr.Context{
		Context:   context.Background(),
		Request:   gofrHttp.NewRequest(req),
		Container: mockContainer,
	}
	return ctx
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		request    string
		userName   string
		expResp    gofrResponse
		mockReturn *models.User
		mockError  error
		ifMock     bool
	}{
		{
			name:     "Success",
			request:  `{"name":"Alice"}`,
			userName: "Alice",
			expResp: gofrResponse{
				result: &models.User{Name: "Alice"},
				err:    nil,
			},
			mockReturn: &models.User{Name: "Alice"},
			mockError:  nil,
			ifMock:     true,
		},
		{
			name:    "Invalid Input",
			request: `{}`, // name is required
			expResp: gofrResponse{
				result: nil,
				err:    errors.New("missing or empty name"),
			},
			ifMock: false,
		},

		{
			name:     "Error from service",
			request:  `{"name":"Bob"}`,
			userName: "Bob",
			expResp: gofrResponse{
				result: nil,
				err:    errors.New("DB error"),
			},
			mockReturn: nil,
			mockError:  errors.New("DB error"),
			ifMock:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := NewMockUserService(ctrl)
			h := &UserHandler{Service: mockSvc}

			ctx := setupTestContext(t, http.MethodPost, "/user", tt.request, nil)

			if tt.ifMock {
				mockSvc.EXPECT().CreateUser(ctx, tt.userName).Return(tt.mockReturn, tt.mockError)
			}

			resp, err := h.CreateUser(ctx)
			assert.Equal(t, tt.expResp.result, resp)
			if tt.name == "Invalid Input" {
				assert.Nil(t, resp)
				assert.Error(t, err)
				return
			}

		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name       string
		idParam    string
		expResp    gofrResponse
		mockReturn *models.User
		mockError  error
		ifMock     bool
	}{
		{
			name:    "Success",
			idParam: "1",
			expResp: gofrResponse{
				result: &models.User{UserID: 1, Name: "Alice"},
				err:    nil,
			},
			mockReturn: &models.User{UserID: 1, Name: "Alice"},
			mockError:  nil,
			ifMock:     true,
		},
		{
			name:    "Invalid ID",
			idParam: "abc",
			expResp: gofrResponse{
				result: nil,
				err:    strconv.ErrSyntax,
			},
			ifMock: false,
		},
		{
			name:    "Not Found",
			idParam: "2",
			expResp: gofrResponse{
				result: "Not Found",
				err:    nil,
			},
			mockReturn: nil,
			mockError:  nil,
			ifMock:     true,
		},
		{
			name:    "Service error",
			idParam: "3",
			expResp: gofrResponse{
				result: nil,
				err:    errors.New("DB error"),
			},
			mockReturn: nil,
			mockError:  errors.New("DB error"),
			ifMock:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSvc := NewMockUserService(ctrl)
			h := &UserHandler{Service: mockSvc}
			ctx := setupTestContext(t, http.MethodGet, "/user/{id}", "", map[string]string{"id": tt.idParam})

			if tt.ifMock {
				id, _ := strconv.Atoi(tt.idParam)
				mockSvc.EXPECT().GetUserByID(ctx, id).Return(tt.mockReturn, tt.mockError)
			}

			resp, err := h.GetUserByID(ctx)
			assert.Equal(t, tt.expResp.result, resp)
			if tt.name == "Invalid ID" {
				assert.Contains(t, err.Error(), "invalid syntax")
				assert.Nil(t, resp)
				return
			}

		})
	}
}

func TestViewUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := NewMockUserService(ctrl)
	h := &UserHandler{Service: mockSvc}

	expectedUsers := []models.User{
		{UserID: 1, Name: "Alice"},
		{UserID: 2, Name: "Bob"},
	}
	ctx := setupTestContext(t, http.MethodGet, "/user", "", nil)

	mockSvc.EXPECT().View_Users(ctx).Return(expectedUsers, nil)

	resp, err := h.ViewUsers(ctx)
	assert.Equal(t, expectedUsers, resp)
	assert.NoError(t, err)
}
