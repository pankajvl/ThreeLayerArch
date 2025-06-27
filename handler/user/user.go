package userhandler

import (
	"ThreeLayerArch/service/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Service *usersvc.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
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

func (h *UserHandler) ViewUsers(w http.ResponseWriter, r *http.Request) {
	ans, err := h.Service.View_Users()
	if err != nil {
		log.Printf("Error in HANDLER.Viewtask: %v", err)
		return
	}
	for _, v := range ans {
		fmt.Fprintf(w, "ID: %d, Name: %s\n", v.UserID, v.Name)
	}
}
