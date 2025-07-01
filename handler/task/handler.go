package task

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handler handles HTTP requests for tasks
type Handler struct {
	service TaskService
}

// New creates a new task handler
func New(service TaskService) *Handler {
	return &Handler{service: service}
}

// Addtask godoc
// @Summary      Add a new task
// @Description  Create a task with the provided task description
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body  object  true  "Task content"
// @Success      201   {string}  string  "Created"
// @Failure      400   {string}  string  "Invalid request body"
// @Failure      500   {string}  string  "Failed to create task"
// @Router       /task [post]
func (h *Handler) Addtask(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Task string `json:"task"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ok, err := h.service.Add_Task(reqBody.Task)
	if err != nil || !ok {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Viewtask godoc
// @Summary      View all tasks
// @Description  Retrieve a list of all tasks
// @Tags         tasks
// @Produce      json
// @Success      200  {array}  models.Tasks
// @Failure      500  {string}  string  "Failed to retrieve tasks"
// @Router       /task [get]
func (h *Handler) Viewtask(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.View_Task()
	if err != nil {
		http.Error(w, "failed to retrieve tasks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

// Gettask godoc
// @Summary      Get a task by ID
// @Description  Retrieve a task by its ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  models.Tasks
// @Failure      400  {string}  string  "Invalid task ID"
// @Failure      404  {string}  string  "Task not found"
// @Router       /task/{id} [get]
func (h *Handler) Gettask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.Get_By_ID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// Updatetask godoc
// @Summary      Update a task
// @Description  Mark a task as completed or update its status
// @Tags         tasks
// @Produce      json
// @Param        id   path  int  true  "Task ID"
// @Success      200  {string}  string  "Updated"
// @Failure      400  {string}  string  "Invalid task ID"
// @Failure      500  {string}  string  "Failed to update task"
// @Router       /task/{id} [put]
func (h *Handler) Updatetask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	ok, err := h.service.Update_Task(id)
	if err != nil || !ok {
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Deletetask godoc
// @Summary      Delete a task
// @Description  Remove a task by its ID
// @Tags         tasks
// @Produce      json
// @Param        id   path  int  true  "Task ID"
// @Success      200  {string}  string  "Deleted"
// @Failure      400  {string}  string  "Invalid task ID"
// @Failure      500  {string}  string  "Failed to delete task"
// @Router       /task/{id} [delete]
func (h *Handler) Deletetask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	ok, err := h.service.Delete_Task(id)
	if err != nil || !ok {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
