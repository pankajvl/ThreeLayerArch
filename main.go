package main

import (
	userhandler "ThreeLayerArch/handler/user"
	usersvc "ThreeLayerArch/service/user"
	userstore "ThreeLayerArch/store/user"
	"log"
	"net/http"
	"time"

	"ThreeLayerArch/datasource"
	handler "ThreeLayerArch/handler/task"
	service "ThreeLayerArch/service/task"
	store "ThreeLayerArch/store/task"
)

func main() {
	db, err := datasource.New("root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Println(err)
		return
	}

	taskStore := store.New(db)
	taskSvc := service.New(taskStore)
	taskHandler := handler.New(taskSvc)

	http.HandleFunc("GET /task", taskHandler.Viewtask)
	http.HandleFunc("GET /task/{id}", taskHandler.Gettask)
	http.HandleFunc("POST /task", taskHandler.Addtask)
	http.HandleFunc("PUT /task/{id}", taskHandler.Updatetask)
	http.HandleFunc("DELETE /task/{id}", taskHandler.Deletetask)

	userStore := &userstore.UserStore{DB: db}
	userService := &usersvc.Service{Store: userStore}
	userHandler := &userhandler.UserHandler{Service: userService}

	http.HandleFunc("POST /user", userHandler.CreateUser)
	http.HandleFunc("GET /user/{id}", userHandler.GetUserByID)
	http.HandleFunc("GET /user", userHandler.ViewUsers)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {

		log.Fatal(err)
	}

}
