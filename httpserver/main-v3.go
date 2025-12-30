package main

import (
	"log"
	"net/http"
)

func homeHandle(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello, World!"))
}

func demoHandle(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is the demo endpoint"))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
func main(){
	mux := http.NewServeMux();
	mux.HandleFunc("/", homeHandle)
	mux.Handle("/demo", loggerMiddleware(http.HandlerFunc(demoHandle)))

	log.Println("Server started on :8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// err := http.ListenAndServe(":8080", mux)
	// http.HandleFunc("/", homeHandle)
	// http.HandleFunc("/demo", demoHandle)
	// log.Println("Server started.....")
	// err := http.ListenAndServe(":8080", nil)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}