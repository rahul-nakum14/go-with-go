package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello World"}`))
	})

		http.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")

		name:= r.URL.Query().Get("name");

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}
		if name == "" {
			name = "World"
		}
			w.Write([]byte(`{"message":"Hello World"}`))

		// w.Write([]byte(`{"name:helo world"}`))
		// w.Write([]byte("Hello " + name))
	})

	log.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {

		log.Fatal(err)
	}
}
