package main

import (
	"log"
	"net/http"
	"github.com/rahulnakum/sqlcrud/internal/config"
	"github.com/rahulnakum/sqlcrud/internal/db"
	"github.com/rahulnakum/sqlcrud/internal/handler"
	"github.com/rahulnakum/sqlcrud/internal/repository"
	"github.com/rahulnakum/sqlcrud/internal/service"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserRepo(database, cfg.DBDriver)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	http.HandleFunc("/users", handler.Create)
	http.HandleFunc("/user", handler.Get)
	http.HandleFunc("/user/update", handler.Update)
	http.HandleFunc("/user/delete", handler.Delete)

	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", nil)
}
