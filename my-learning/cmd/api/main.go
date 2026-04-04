package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rahul-nakum14/go-user-crud/internal/config"
	"github.com/rahul-nakum14/go-user-crud/internal/db"
	"github.com/rahul-nakum14/go-user-crud/internal/handler"
	"github.com/rahul-nakum14/go-user-crud/internal/repository"
	"github.com/rahul-nakum14/go-user-crud/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	dbConn, err := db.ConnectPostgres(cfg.DbUrl)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// userService := &service.UserService{UserRepo: userRepo}
	// userHandler := &handler.UserHandler{UserService: userService}

	router := mux.NewRouter()
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	log.Printf("Server is running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}