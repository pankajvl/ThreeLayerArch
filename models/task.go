package models

type Tasks struct {
	Tid       int    `json:"id"`
	Task      string `json:"title"`
	Completed bool   `json:"completed"`
}

type Users struct {
	userID int
	name   string
}
