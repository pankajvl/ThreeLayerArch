package task

import (
	Models "ThreeLayerArch/models"
	"database/sql"
	"log"
)

type Store struct {
	db *sql.DB
}

// New creates a new task store
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AddTask(task string) (bool, error) {

	_, err := s.db.Exec("Insert into tasks (description,completed) values (?,?)", task, false)
	if err != nil {
		log.Printf("Error in STORE.AddTask: %v", err)
		return false, err
	}
	return true, nil
}

func (s *Store) ViewTask() ([]Models.Tasks, error) {

	var tID int
	var task string
	var completed bool

	var answers []Models.Tasks

	row, err := s.db.Query("select * from tasks")
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

func (s *Store) GetByID(id int) (Models.Tasks, error) {

	var tID int
	var task string
	var completed bool

	err := s.db.QueryRow("select * from tasks where id=?", id).Scan(&tID, &task, &completed)
	if err != nil {
		log.Printf("Error in STORE.GetByID: %v", err)
		return Models.Tasks{}, err
	}
	return Models.Tasks{tID, task, completed}, nil
}

func (s *Store) UpdateTask(id int) (bool, error) {

	_, err := s.db.Exec("UPDATE tasks SET completed= true WHERE id=?", id)
	if err != nil {
		log.Printf("Error in STORE.UpdateTask: %v", err)
		return false, err
	}
	return true, nil
}

func (s *Store) DeleteTask(id int) (bool, error) {

	_, err := s.db.Exec("delete from tasks where id=?", id)
	if err != nil {
		log.Printf("Error in STORE.DeleteTask: %v", err)
		return false, err
	}
	return true, nil
}

func (s *Store) CheckIfExists(i int) bool {
	ans := s.db.QueryRow("select id from tasks where id=?", i)

	var id int
	err := ans.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return true
}
