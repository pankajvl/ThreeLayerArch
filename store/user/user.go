package userstore

import (
	"ThreeLayerArch/models"
	"database/sql"
	"log"
)

type UserStore struct {
	DB *sql.DB
}

func (s *UserStore) Create(name string) (int, error) {
	res, err := s.DB.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (s *UserStore) GetByID(id int) (*models.User, error) {
	var user models.User
	err := s.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id).
		Scan(&user.UserID, &user.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) ViewUsers() ([]models.User, error) {

	var uID int
	var Name string

	var answers []models.User

	row, err := s.DB.Query("select * from users")
	if err != nil {
		log.Printf("Error in STORE.View: %v", err)
		return []models.User{}, err
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&uID, &Name)
		if err != nil {
			log.Printf("Error in STORE.View: %v", err)
			return []models.User{}, err
		}
		answers = append(answers, models.User{uID, Name})

	}
	return answers, nil
}
