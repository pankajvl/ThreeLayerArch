package task

import (
	Models "ThreeLayerArch/models"
	"database/sql"
	"gofr.dev/pkg/gofr"
	"log"
)

type Store struct {
	db *sql.DB
}

// New creates a new task store
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AddTask(ctx *gofr.Context, task string) (bool, error) {
	_, err := ctx.SQL.ExecContext(ctx, "INSERT INTO tasks (description, completed) VALUES (?, ?)", task, false)
	if err != nil {
		log.Printf("Error in STORE.AddTask: %v", err)
		return false, err
	}
	return true, nil
}

func (s *Store) ViewTask(ctx *gofr.Context) ([]Models.Tasks, error) {
	var tID int
	var task string
	var completed bool
	var answers []Models.Tasks

	row, err := ctx.SQL.QueryContext(ctx, "SELECT * FROM tasks")
	if err != nil {
		log.Printf("Error in STORE.View: %v", err)
		return []Models.Tasks{}, err
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&tID, &task, &completed)
		if err != nil {
			log.Printf("Error in STORE.View: %v", err)
			return []Models.Tasks{}, err
		}
		answers = append(answers, Models.Tasks{tID, task, completed})
	}
	return answers, nil
}

func (s *Store) GetByID(ctx *gofr.Context, id int) (Models.Tasks, error) {
	var tID int
	var task string
	var completed bool

	err := ctx.SQL.QueryRowContext(ctx, "SELECT * FROM tasks WHERE id = ?", id).Scan(&tID, &task, &completed)
	if err != nil {
		log.Printf("Error in STORE.GetByID: %v", err)
		return Models.Tasks{}, err
	}
	return Models.Tasks{tID, task, completed}, nil
}

func (s *Store) UpdateTask(ctx *gofr.Context, id int) (bool, error) {
	_, err := ctx.SQL.ExecContext(ctx, "UPDATE tasks SET completed = TRUE WHERE id = ?", id)
	if err != nil {
		log.Printf("Error in STORE.UpdateTask: %v", err)
		return false, err
	}
	return true, nil
}

func (s *Store) DeleteTask(ctx *gofr.Context, id int) (bool, error) {
	_, err := ctx.SQL.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Printf("Error in STORE.DeleteTask: %v", err)
		return false, err
	}
	return true, nil
}

func (s *Store) CheckIfExists(ctx *gofr.Context, i int) bool {
	ans, err := ctx.SQL.QueryContext(ctx, "SELECT id FROM tasks WHERE id = ?", i)

	var id int
	err = ans.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return true
}
