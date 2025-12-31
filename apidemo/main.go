package main

import (
	"context"
	"learn-go/apidemo/internal/config"
	"learn-go/apidemo/internal/http/handlers/student"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	cfg, _ := config.LoadConfig()
	if cfg != nil {

		mux := http.NewServeMux()

		mux.HandleFunc("POST /api/new/student", student.Create)

		logger.Println("Server started on: ", cfg.HTTPServer.Address)
		server := &http.Server{
			Addr:    cfg.HTTPServer.Address,
			Handler: mux,
		}

		done := make(chan os.Signal, 1)

		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			if err := server.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}()
		<-done

		log.Printf("Shutting Down Server Gracefully...")

		if err := server.Close(); err != nil {
			log.Fatalf("Could not close server: %s", err)
		}
		log.Printf("Server Stopped")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+s", err)
		}
		log.Printf("Server Exited Properly")
	}
}
