package main

import (
	"ThreeLayerArch/datasource"
	_ "ThreeLayerArch/docs"
	handler "ThreeLayerArch/handler/task"
	userhandler "ThreeLayerArch/handler/user"
	service "ThreeLayerArch/service/task"
	usersvc "ThreeLayerArch/service/user"
	store "ThreeLayerArch/store/task"
	userstore "ThreeLayerArch/store/user"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"time"
)

// @title           Task Management API
// @version         1.0
// @description     API for managing tasks and users.
// @host            localhost:8080
// @BasePath        /

// @contact.name   Pankaj Venkat
// @contact.email  pankaj@example.com

func main() {
	db, err := datasource.New("root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Println(err)
		return
	}

	taskStore := store.New(db)
	taskSvc := service.New(taskStore)
	taskHandler := handler.New(taskSvc)

	// @Summary      Get all tasks
	// @Description  Returns a list of all tasks
	// @Tags         tasks
	// @Produce      json
	// @Success      200  {array}  models.Task
	// @Router       /task [get]
	http.HandleFunc("GET /task", taskHandler.Viewtask)

	// @Summary      Get task by ID
	// @Description  Returns a task by its ID
	// @Tags         tasks
	// @Produce      json
	// @Param        id   path      int  true  "Task ID"
	// @Success      200  {object}  models.Task
	// @Router       /task/{id} [get]
	http.HandleFunc("GET /task/{id}", taskHandler.Gettask)

	// @Summary      Create new task
	// @Description  Adds a new task
	// @Tags         tasks
	// @Accept       json
	// @Produce      json
	// @Param        task  body      models.Task  true  "Task to add"
	// @Success      201   {object}  models.Task
	// @Router       /task [post]
	http.HandleFunc("POST /task", taskHandler.Addtask)

	// @Summary      Update task
	// @Description  Updates a task by ID
	// @Tags         tasks
	// @Accept       json
	// @Produce      json
	// @Param        id    path      int          true  "Task ID"
	// @Param        task  body      models.Task  true  "Updated Task"
	// @Success      200   {object}  models.Task
	// @Router       /task/{id} [put]
	http.HandleFunc("PUT /task/{id}", taskHandler.Updatetask)

	// @Summary      Delete task
	// @Description  Deletes a task by ID
	// @Tags         tasks
	// @Param        id   path  int  true  "Task ID"
	// @Success      204  "No Content"
	// @Router       /task/{id} [delete]
	http.HandleFunc("DELETE /task/{id}", taskHandler.Deletetask)

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
	http.HandleFunc("POST /user", userHandler.CreateUser)

	// @Summary      Get user by ID
	// @Description  Returns a user by ID
	// @Tags         users
	// @Produce      json
	// @Param        id   path      int  true  "User ID"
	// @Success      200  {object}  models.User
	// @Router       /user/{id} [get]
	http.HandleFunc("GET /user/{id}", userHandler.GetUserByID)

	// @Summary      View all users
	// @Description  Returns a list of users
	// @Tags         users
	// @Produce      json
	// @Success      200  {array}  models.User
	// @Router       /user [get]
	http.HandleFunc("GET /user", userHandler.ViewUsers)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Swagger UI at http://localhost:8080/swagger/index.html")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
