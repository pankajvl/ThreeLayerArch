package user

import (
	"ThreeLayerArch/models"
	"bytes"
	"errors"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		path      string
		body      []byte
		setupMock func(*MockUserService)
		handlerFn func(h *UserHandler, w http.ResponseWriter, r *http.Request)
		wantCode  int
		wantBody  []byte
		setPath   func(r *http.Request)
	}{
		{
			name:   "CreateUser_Success",
			method: http.MethodPost,
			path:   "/users",
			body:   []byte(`{"name": "John"}`),
			setupMock: func(m *MockUserService) {
				m.EXPECT().
					CreateUser("John").
					Return(&models.User{UserID: 1, Name: "John"}, nil)
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.CreateUser(w, r)
			},
			wantCode: http.StatusCreated,
		},
		{
			name:   "CreateUser_InvalidBody",
			method: http.MethodPost,
			path:   "/users",
			body:   []byte(`invalid-json`),
			setupMock: func(m *MockUserService) {
				// No expectation, parsing fails before call
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.CreateUser(w, r)
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "CreateUser_Failure_FromService",
			method: http.MethodPost,
			path:   "/users",
			body:   []byte(`{"name": "John"}`),
			setupMock: func(m *MockUserService) {
				m.EXPECT().
					CreateUser("John").
					Return(nil, errors.New("insert failed"))
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.CreateUser(w, r)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:   "GetUserByID_Success",
			method: http.MethodGet,
			path:   "/users/1",
			setupMock: func(m *MockUserService) {
				m.EXPECT().
					GetUserByID(1).
					Return(&models.User{UserID: 1, Name: "John"}, nil)
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.GetUserByID(w, r)
			},
			wantCode: http.StatusOK,
			setPath: func(r *http.Request) {
				r.SetPathValue("id", "1")
			},
		},
		{
			name:   "GetUserByID_NotFound",
			method: http.MethodGet,
			path:   "/users/99",
			setupMock: func(m *MockUserService) {
				m.EXPECT().
					GetUserByID(99).
					Return(nil, nil)
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.GetUserByID(w, r)
			},
			wantCode: http.StatusNotFound,
			setPath: func(r *http.Request) {
				r.SetPathValue("id", "99")
			},
		},
		{
			name:   "GetUserByID_Error_FromService",
			method: http.MethodGet,
			path:   "/users/3",
			setupMock: func(m *MockUserService) {
				m.EXPECT().
					GetUserByID(3).
					Return(nil, errors.New("db error"))
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.GetUserByID(w, r)
			},
			wantCode: http.StatusInternalServerError,
			setPath: func(r *http.Request) {
				r.SetPathValue("id", "3")
			},
		},
		{
			name:   "GetUserByID_InvalidID",
			method: http.MethodGet,
			path:   "/users/abc",
			setupMock: func(m *MockUserService) {
				// No call expected
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.GetUserByID(w, r)
			},
			wantCode: http.StatusBadRequest,
			setPath: func(r *http.Request) {
				r.SetPathValue("id", "abc")
			},
		},
		{
			name:   "ViewUsers_Success",
			method: http.MethodGet,
			path:   "/users",
			setupMock: func(m *MockUserService) {
				m.EXPECT().
					View_Users().
					Return([]models.User{
						{UserID: 1, Name: "Alice"},
						{UserID: 2, Name: "Bob"},
					}, nil)
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.ViewUsers(w, r)
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "ViewUsers_Error",
			method: http.MethodGet,
			path:   "/users",
			setupMock: func(m *MockUserService) {
				m.EXPECT().
					View_Users().
					Return(nil, errors.New("db error"))
			},
			handlerFn: func(h *UserHandler, w http.ResponseWriter, r *http.Request) {
				h.ViewUsers(w, r)
			},
			wantCode: http.StatusInternalServerError, // Assuming your handler logs and returns 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := NewMockUserService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			handler := &UserHandler{Service: mockSvc}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(tt.body))
			if tt.setPath != nil {
				tt.setPath(req)
			}

			rec := httptest.NewRecorder()

			tt.handlerFn(handler, rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("expected status %d, got %d", tt.wantCode, rec.Code)
			}

			if tt.wantBody != nil && !bytes.Contains(rec.Body.Bytes(), tt.wantBody) {
				t.Errorf("expected body to contain %q, got %q", tt.wantBody, rec.Body.String())
			}
		})
	}
}
