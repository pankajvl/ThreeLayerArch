package user

import (
	"ThreeLayerArch/models"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"testing"
)

func setupContext(t *testing.T) (*gofr.Context, *container.Mocks) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   t.Context(),
		Request:   nil,
		Container: mockContainer,
	}
	return ctx, mock
}

func TestCreateUser(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		input    string
		mockFunc func()
		expected int
		wantErr  bool
	}{
		{
			name:  "Success",
			input: "Alice",
			mockFunc: func() {
				mockSQL.SQL.ExpectExec("INSERT INTO users (name) VALUES (?)").
					WithArgs("Alice").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: 1,
			wantErr:  false,
		},
		{
			name:  "Failure",
			input: "Bob",
			mockFunc: func() {
				mockSQL.SQL.ExpectExec("INSERT INTO users (name) VALUES (?)").
					WithArgs("Bob").
					WillReturnError(sql.ErrConnDone)
			},
			expected: 0,
			wantErr:  true,
		},
	}

	s := &UserStore{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := s.Create(ctx, tt.input)
			assert.Equal(t, tt.expected, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		input    int
		mockFunc func()
		expected *models.User
		wantErr  bool
	}{
		{
			name:  "Success",
			input: 1,
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
					WithArgs(1).
					WillReturnRows(mockSQL.SQL.NewRows([]string{"id", "name"}).AddRow(1, "Alice"))
			},
			expected: &models.User{UserID: 1, Name: "Alice"},
			wantErr:  false,
		},
		{
			name:  "No Rows",
			input: 2,
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name:  "Error",
			input: 3,
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expected: nil,
			wantErr:  true,
		},
	}

	s := &UserStore{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := s.GetByID(ctx, tt.input)
			assert.Equal(t, tt.expected, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestViewUsers(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		mockFunc func()
		expected []models.User
		wantErr  bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				rows := mockSQL.SQL.NewRows([]string{"id", "name"}).
					AddRow(1, "Alice").
					AddRow(2, "Bob")
				mockSQL.SQL.ExpectQuery("SELECT * FROM users").
					WillReturnRows(rows)
			},
			expected: []models.User{
				{UserID: 1, Name: "Alice"},
				{UserID: 2, Name: "Bob"},
			},
			wantErr: false,
		},
		{
			name: "Query error",
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT * FROM users").
					WillReturnError(sql.ErrNoRows)
			},
			expected: []models.User{},
			wantErr:  true,
		},
		{
			name: "Row scan error",
			mockFunc: func() {
				rows := mockSQL.SQL.NewRows([]string{"id", "name"}).
					AddRow("invalid_id", "Charlie") // id should be int
				mockSQL.SQL.ExpectQuery("SELECT * FROM users").
					WillReturnRows(rows)
			},
			expected: []models.User{},
			wantErr:  true,
		},
	}

	s := &UserStore{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := s.ViewUsers(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expected, got) // handles empty slice vs nil
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

		})
	}
}
