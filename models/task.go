package models

type Tasks struct {
	Tid       int
	Task      string
	Completed bool
	//userID    int
}

type Users struct {
	userID int
	name   string
}
