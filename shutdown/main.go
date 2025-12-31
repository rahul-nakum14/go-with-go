package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"syscall"
	"context"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Simulate a long-running request
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "Request completed hell yeah")
}

func main() {
	// Create a new HTTP server
	server := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/hello", handler)

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server is starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting the server:", err)
		}
	}()

	// Create a channel to listen for OS signals (e.g., Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal (like SIGINT)
	<-stop
	fmt.Println("\nShutting down gracefully...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown error:", err)
	} else {
		fmt.Println("Server shut down gracefully")
	}
}

/**
	Not a Graceful Shutdown Example
**/
// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func main(){

// 	mux :=  http.NewServeMux();
// 	mux.HandleFunc("/hello", func (w http.ResponseWriter, r *http.Request){
// 		fmt.Fprintf(w, "Hello, World!")
// 	})

// 	mux.HandleFunc("/",func (w http.ResponseWriter, r *http.Request){
// 		time.Sleep(10* time.Second)
// 		w.Write([]byte("Request is Completed"))
// 	})

// 	server := &http.Server{
// 		Addr: ":8080",
// 		Handler: mux,
// 	}

// 	if err:= server.ListenAndServe(); err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// 	// err:= http.ListenAndServe(":8080", mux)
// 	// if err != nil {
// 	// 	fmt.Println("Error starting server:", err)
// 	// }
// }