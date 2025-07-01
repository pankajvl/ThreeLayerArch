package user

import (
	"ThreeLayerArch/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// UserHandler handles user-related endpoints.
type UserHandler struct {
	Service UserService
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Accepts a user name and creates a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserRequest  true  "User to create"
// @Success      201   {object}  models.User
// @Failure      400   {string}  string  "Invalid request"
// @Failure      500   {string}  string  "Failed to create user"
// @Router       /user [post]

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.Service.CreateUser(req.Name)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieves a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      404  {string}  string  "User not found"
// @Failure      500  {string}  string  "Error fetching user"
// @Router       /user/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(id)

	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// ViewUsers godoc
// @Summary      Get all users
// @Description  Returns a list of all users
// @Tags         users
// @Produce      plain
// @Success      200  {string}  string  "List of users"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /user [get]
func (h *UserHandler) ViewUsers(w http.ResponseWriter, _ *http.Request) {
	ans, err := h.Service.View_Users()

	if err != nil {
		log.Printf("Error in HANDLER.ViewUsers: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, v := range ans {
		fmt.Fprintf(w, "ID: %d, Name: %s\n", v.UserID, v.Name)
	}
}
