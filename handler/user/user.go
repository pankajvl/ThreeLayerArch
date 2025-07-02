package user

import (
	"ThreeLayerArch/models"
	"gofr.dev/pkg/gofr"
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

func (h *UserHandler) CreateUser(ctx *gofr.Context) (any, error) {
	var req models.UserRequest

	if err := ctx.Bind(&req); err != nil || req.Name == "" {
		return nil, err
	}

	user, err := h.Service.CreateUser(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return user, nil
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
func (h *UserHandler) GetUserByID(ctx *gofr.Context) (any, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, err
	}

	user, err := h.Service.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return "Not Found", nil
	}
	return user, nil
}

// ViewUsers godoc
// @Summary      Get all users
// @Description  Returns a list of all users
// @Tags         users
// @Produce      plain
// @Success      200  {string}  string  "List of users"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /user [get]
func (h *UserHandler) ViewUsers(ctx *gofr.Context) (any, error) {
	ans, err := h.Service.View_Users(ctx)

	if err != nil {
		return nil, err
	}
	return ans, nil
}
