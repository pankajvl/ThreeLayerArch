package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createUser = `CREATE TABLE IF NOT EXISTS users (
       id INT AUTO_INCREMENT PRIMARY KEY,
       name VARCHAR(255)
     );`

func create_user_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// write your migrations here
			_, err := d.SQL.Exec(createUser)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
