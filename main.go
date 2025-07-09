package main

import (
	"ThreeLayerArch/datasource"
	_ "ThreeLayerArch/docs"
	handler "ThreeLayerArch/handler/task"
	userhandler "ThreeLayerArch/handler/user"
	"ThreeLayerArch/migrations/migrations"
	service "ThreeLayerArch/service/task"
	usersvc "ThreeLayerArch/service/user"
	store "ThreeLayerArch/store/task"
	userstore "ThreeLayerArch/store/user"
	"github.com/swaggo/http-swagger"
	"gofr.dev/pkg/gofr"
	"log"
	"net/http"
	"os"
)

// @title           Task Management API
// @version         1.0
// @description     API for managing tasks and users.
// @host            localhost:8080
// @BasePath        /

// @contact.name   Pankaj Venkat
// @contact.email  pankaj@example.com

func main() {
	db, err := datasource.New(os.Getenv("DB_CONN_STRING"))

	if err != nil {
		log.Println(err)
		return
	}
	app := gofr.New()

	app.Migrate(migrations.All())

	taskStore := store.New(db)
	taskSvc := service.New(taskStore)
	taskHandler := handler.New(taskSvc)
	app.GET("/", func(ctx *gofr.Context) (any, error) {
		return "Hello World", nil
	})

	// @Summary      Get all tasks
	// @Description  Returns a list of all tasks
	// @Tags         tasks
	// @Produce      json
	// @Success      200  {array}  models.Task
	// @Router       /task [get]
	app.GET("/task", taskHandler.Viewtask)

	// @Summary      Get task by ID
	// @Description  Returns a task by its ID
	// @Tags         tasks
	// @Produce      json
	// @Param        id   path      int  true  "Task ID"
	// @Success      200  {object}  models.Task
	// @Router       /task/{id} [get]
	app.GET("/task/{id}", taskHandler.Gettask)

	// @Summary      Create new task
	// @Description  Adds a new task
	// @Tags         tasks
	// @Accept       json
	// @Produce      json
	// @Param        task  body      models.Task  true  "Task to add"
	// @Success      201   {object}  models.Task
	// @Router       /task [post]
	app.POST("/task", taskHandler.Addtask)

	// @Summary      Update task
	// @Description  Updates a task by ID
	// @Tags         tasks
	// @Accept       json
	// @Produce      json
	// @Param        id    path      int          true  "Task ID"
	// @Param        task  body      models.Task  true  "Updated Task"
	// @Success      200   {object}  models.Task
	// @Router       /task/{id} [put]
	app.PUT("/task/{id}", taskHandler.Updatetask)

	// @Summary      Delete task
	// @Description  Deletes a task by ID
	// @Tags         tasks
	// @Param        id   path  int  true  "Task ID"
	// @Success      204  "No Content"
	// @Router       /task/{id} [delete]
	app.DELETE("/task/{id}", taskHandler.Deletetask)

	userStore := &userstore.UserStore{DB: db}
	userService := &usersvc.Service{Store: userStore}
	userHandler := &userhandler.UserHandler{Service: userService}

	// @Summary      Create user
	// @Description  Creates a new user
	// @Tags         users
	// @Accept       json
	// @Produce      json
	// @Param        user  body      models.User  true  "User to create"
	// @Success      201   {object}  models.User
	// @Router       /user [post]
	app.POST("/user", userHandler.CreateUser)

	// @Summary      Get user by ID
	// @Description  Returns a user by ID
	// @Tags         users
	// @Produce      json
	// @Param        id   path      int  true  "User ID"
	// @Success      200  {object}  models.User
	// @Router       /user/{id} [get]
	app.GET("/user/{id}", userHandler.GetUserByID)

	// @Summary      View all users
	// @Description  Returns a list of users
	// @Tags         users
	// @Produce      json
	// @Success      200  {array}  models.User
	// @Router       /user [get]
	app.GET("/user", userHandler.ViewUsers)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Swagger UI at http://localhost:9000/swagger/index.html")

	app.Run()
}
