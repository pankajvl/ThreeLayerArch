package task_test

import (
	"database/sql"
	"testing"

	taskstore "ThreeLayerArch/store/task"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTaskByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	rows := sqlmock.NewRows([]string{"id", "description", "completed"}).
		AddRow(1, "Task A", true)

	mock.ExpectQuery(`SELECT \* FROM tasks WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(rows)

	task, err := store.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if task.Tid != 1 || task.Task != "Task A" || task.Completed != true {
		t.Errorf("unexpected task result: %+v", task)
	}
}

func TestAddTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	mock.ExpectExec(`INSERT INTO tasks \(description, completed\) VALUES \(\?, \?\)`).
		WithArgs("New Task", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ok, err := store.AddTask("New Task")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !ok {
		t.Errorf("expected true, got false")
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	mock.ExpectExec(`UPDATE tasks SET completed = TRUE WHERE id = \?`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := store.UpdateTask(1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !ok {
		t.Errorf("expected true, got false")
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	mock.ExpectExec(`DELETE FROM tasks WHERE id = \?`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := store.DeleteTask(1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !ok {
		t.Errorf("expected true, got false")
	}
}

func TestCheckIfExists(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery(`SELECT id FROM tasks WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(rows)

	ok := store.CheckIfExists(1)
	if !ok {
		t.Errorf("expected true, got false")
	}
}

func TestViewTask(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	rows := sqlmock.NewRows([]string{"id", "description", "completed"}).
		AddRow(1, "Task A", true).
		AddRow(2, "Task B", false)

	mock.ExpectQuery(`SELECT \* FROM tasks`).
		WillReturnRows(rows)

	tasks, err := store.ViewTask()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Tid != 1 || tasks[0].Task != "Task A" {
		t.Errorf("unexpected first task: %+v", tasks[0])
	}
	if tasks[1].Tid != 2 || tasks[1].Task != "Task B" {
		t.Errorf("unexpected second task: %+v", tasks[1])
	}
}

func TestGetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	mock.ExpectQuery(`SELECT \* FROM tasks WHERE id = \?`).
		WithArgs(10).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = store.GetByID(10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckIfExists_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := taskstore.New(db)

	mock.ExpectQuery(`SELECT id FROM tasks WHERE id = \?`).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	ok := store.CheckIfExists(99)
	if ok {
		t.Errorf("expected false, got true")
	}
}
