package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createTask = `CREATE TABLE IF NOT EXISTS tasks (
         id INT AUTO_INCREMENT PRIMARY KEY,
         description VARCHAR(255) NOT NULL,
         completed BOOLEAN NOT NULL DEFAULT FALSE
     );`

func create_tasks_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// write your migrations here
			_, err := d.SQL.Exec(createTask)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
