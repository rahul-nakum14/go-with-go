// // package main

// // import (
// // 	"log"
// // 	"net/http"
// // )

// // type Server struct {
// // 	logger *log.Logger
// // }


// // func homeHandler(w http.ResponseWriter, r *http.Request){
// // 		w.WriteHeader(http.StatusOK)
// // 		w.Write([]byte("Hello, Go HTTP Server"))
// // 	}

// // func healthHandler(w http.ResponseWriter, r *http.Request) {
// // 	w.WriteHeader(http.StatusOK)
// // 	w.Write([]byte("OK"))
// // }

// // func main(){
// // 	logger := log.New(os.Stdout, "HTTP: ", log.LstdFlags)
// // 	server := &Server{
// // 		logger: logger,
// // 	}

// // 	mux := http.NewServeMux()
// // 	mux.HandleFunc("/health", server.healthHandler)
// // 	mux.HandleFunc("/users", server.usersHandler)

// // 	logger.Println("Server started on :8080")

// // 	err := http.ListenAndServe(":8080", mux)
// // 	if err != nil {
// // 		logger.Fatal(err)
// // 	}

// // 	http.HandleFunc("/",homeHandler)
// // 	http.HandleFunc("/health", healthHandler)

// // 	log.Println("Server started on :8080")

// // 	err := http.ListenAndServe(":8080", nil)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // }
// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// )

// type Server struct {
// 	logger *log.Logger
// }

// func main() {
// 	logger := log.New(os.Stdout, "HTTP: ", log.LstdFlags)

// 	server := &Server{
// 		logger: logger,
// 	}

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/health", server.healthHandler)
// 	mux.HandleFunc("/users", server.usersHandler)

// 	logger.Println("Server started on :8080")

// 	err := http.ListenAndServe(":8080", mux)
// 	if err != nil {
// 		logger.Fatal(err)
// 	}
// }

// func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
// 	s.logger.Println("Health check called")

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }

// func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
// 	s.logger.Println("Users endpoint called")

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Users list"))
// }
