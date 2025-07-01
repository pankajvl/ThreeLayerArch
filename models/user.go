package models

type User struct {
	UserID int
	Name   string
}

// UserRequest is used for creating a user.
type UserRequest struct {
	Name string `json:"name" example:"John Doe"`
}
