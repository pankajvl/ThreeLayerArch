package user

import (
	"ThreeLayerArch/models"
	"database/sql"
	"gofr.dev/pkg/gofr"
	"log"
)

type UserStore struct {
	DB *sql.DB
}

func (s *UserStore) Create(ctx *gofr.Context, name string) (int, error) {
	res, err := ctx.SQL.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", name)

	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (s *UserStore) GetByID(ctx *gofr.Context, id int) (*models.User, error) {
	var user models.User
	err := ctx.SQL.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = ?", id).
		Scan(&user.UserID, &user.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) ViewUsers(ctx *gofr.Context) ([]models.User, error) {

	var uID int
	var Name string

	var answers []models.User

	row, err := ctx.SQL.QueryContext(ctx, "SELECT * FROM users")
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
