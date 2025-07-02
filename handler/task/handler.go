package task

import (
	"gofr.dev/pkg/gofr"
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
func (h *Handler) Addtask(ctx *gofr.Context) (any, error) {
	var reqBody struct {
		Task string `json:"task"`
	}
	err := ctx.Bind(&reqBody)
	if err != nil {
		return nil, err
	}

	ok, err := h.service.Add_Task(ctx, reqBody.Task)
	if err != nil || !ok {
		return nil, err
	}

	return "inserted successfully", nil
}

// Viewtask godoc
// @Summary      View all tasks
// @Description  Retrieve a list of all tasks
// @Tags         tasks
// @Produce      json
// @Success      200  {array}  models.Tasks
// @Failure      500  {string}  string  "Failed to retrieve tasks"
// @Router       /task [get]
func (h *Handler) Viewtask(ctx *gofr.Context) (any, error) {
	tasks, err := h.service.View_Task(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
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
func (h *Handler) Gettask(ctx *gofr.Context) (any, error) {
	idStr, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	task, err := h.service.Get_By_ID(ctx, idStr)
	if err != nil {
		return nil, err
	}

	return task, nil
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
func (h *Handler) Updatetask(ctx *gofr.Context) (any, error) {
	idStr, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	ok, err := h.service.Update_Task(ctx, idStr)
	if err != nil || !ok {
		return nil, err
	}

	return "updated successfully", nil
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
func (h *Handler) Deletetask(ctx *gofr.Context) (any, error) {
	idStr, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	ok, err := h.service.Delete_Task(ctx, idStr)
	if err != nil || !ok {
		return nil, err
	}

	return "deleted successfully", nil
}
