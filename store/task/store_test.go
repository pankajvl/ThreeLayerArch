package task

import (
	Models "ThreeLayerArch/models"
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

func TestViewTask(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		mockFunc func()
		expected []Models.Tasks
		wantErr  bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				rows := mockSQL.SQL.NewRows([]string{"id", "description", "completed"}).
					AddRow(1, "Task 1", false).
					AddRow(2, "Task 2", true)
				mockSQL.SQL.ExpectQuery("SELECT * FROM tasks").
					WillReturnRows(rows)
			},
			expected: []Models.Tasks{
				{Tid: 1, Task: "Task 1", Completed: false},
				{Tid: 2, Task: "Task 2", Completed: true},
			},
			wantErr: false,
		},
		{
			name: "Query error",
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT * FROM tasks").
					WillReturnError(sql.ErrConnDone)
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Scan error",
			mockFunc: func() {
				rows := mockSQL.SQL.NewRows([]string{"id", "description", "completed"}).
					AddRow("invalid_id", "Task 1", false)
				mockSQL.SQL.ExpectQuery("SELECT * FROM tasks").
					WillReturnRows(rows)
			},
			expected: nil,
			wantErr:  true,
		},
	}

	s := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := s.ViewTask(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})

	}
}

func TestGetByID(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		id       int
		mockFunc func()
		expected Models.Tasks
		wantErr  bool
	}{
		{
			name: "Success",
			id:   1,
			mockFunc: func() {
				rows := mockSQL.SQL.NewRows([]string{"id", "description", "completed"}).
					AddRow(1, "Single Task", false)
				mockSQL.SQL.ExpectQuery("SELECT * FROM tasks WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: Models.Tasks{Tid: 1, Task: "Single Task", Completed: false},
			wantErr:  false,
		},
		{
			name: "Query error",
			id:   2,
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT * FROM tasks WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrConnDone)
			},
			expected: Models.Tasks{},
			wantErr:  true,
		},
	}

	s := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := s.GetByID(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, Models.Tasks{}, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		id       int
		mockFunc func()
		expected bool
		wantErr  bool
	}{
		{
			name: "Success",
			id:   1,
			mockFunc: func() {
				mockSQL.SQL.ExpectExec("UPDATE tasks SET completed = TRUE WHERE id = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "Exec error",
			id:   2,
			mockFunc: func() {
				mockSQL.SQL.ExpectExec("UPDATE tasks SET completed = TRUE WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrConnDone)
			},
			expected: false,
			wantErr:  true,
		},
	}

	s := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := s.UpdateTask(ctx, tt.id)
			assert.Equal(t, tt.expected, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		id       int
		mockFunc func()
		expected bool
		wantErr  bool
	}{
		{
			name: "Success",
			id:   1,
			mockFunc: func() {
				mockSQL.SQL.ExpectExec("DELETE FROM tasks WHERE id = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "Exec error",
			id:   2,
			mockFunc: func() {
				mockSQL.SQL.ExpectExec("DELETE FROM tasks WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrConnDone)
			},
			expected: false,
			wantErr:  true,
		},
	}

	s := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := s.DeleteTask(ctx, tt.id)
			assert.Equal(t, tt.expected, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckIfExists(t *testing.T) {
	ctx, mockSQL := setupContext(t)

	tests := []struct {
		name     string
		id       int
		mockFunc func()
		expected bool
	}{
		{
			name: "Exists",
			id:   1,
			mockFunc: func() {
				rows := mockSQL.SQL.NewRows([]string{"id"}).AddRow(1)
				mockSQL.SQL.ExpectQuery("SELECT id FROM tasks WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: true,
		},
		{
			name: "Does not exist",
			id:   2,
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT id FROM tasks WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expected: false,
		},
		{
			name: "Query error",
			id:   3,
			mockFunc: func() {
				mockSQL.SQL.ExpectQuery("SELECT id FROM tasks WHERE id = ?").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expected: false,
		},
	}

	s := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got := s.CheckIfExists(ctx, tt.id)
			assert.Equal(t, tt.expected, got)
		})
	}
}
