package datasource

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func New(creds string) (*sql.DB, error) {
	return sql.Open("mysql", creds)
}
