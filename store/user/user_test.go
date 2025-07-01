package user_test

import (
	userstore "ThreeLayerArch/store/user"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := &userstore.UserStore{DB: db}

	mock.ExpectExec(`INSERT INTO users \(name\) VALUES \(\?\)`).
		WithArgs("Charlie").
		WillReturnResult(sqlmock.NewResult(3, 1))

	id, err := store.Create("Charlie")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if id != 3 {
		t.Errorf("expected ID 3, got %d", id)
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := &userstore.UserStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(5, "Diana")

	mock.ExpectQuery(`SELECT id, name FROM users WHERE id = \?`).
		WithArgs(5).
		WillReturnRows(rows)

	user, err := store.GetByID(5)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if user == nil || user.UserID != 5 || user.Name != "Diana" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := &userstore.UserStore{DB: db}

	mock.ExpectQuery(`SELECT id, name FROM users WHERE id = \?`).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	user, err := store.GetByID(99)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}

func TestViewUsers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Close()

	store := &userstore.UserStore{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(7, "Eve").
		AddRow(8, "Frank")

	mock.ExpectQuery(`SELECT \* FROM users`).
		WillReturnRows(rows)

	users, err := store.ViewUsers()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	if users[0].UserID != 7 || users[0].Name != "Eve" {
		t.Errorf("unexpected first user: %+v", users[0])
	}
	if users[1].UserID != 8 || users[1].Name != "Frank" {
		t.Errorf("unexpected second user: %+v", users[1])
	}
}
