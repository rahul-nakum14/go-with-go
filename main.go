packaged main

import (
	"log"
	"net/http"
)

func demoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Hello ` + name + `"}`))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {

	/**
		Old Way - No server Mux

		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello World"}`))
		})

		http.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")

			name := r.URL.Query().Get("name")

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
	**/
	mux := http.NewServeMux()
	mux.HandleFunc("/demo", demoHandler)
	mux.HandleFunc("/health", healthHandler)
	log.Println("Server started on :8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
	// err := server.ListenAndServe
	// if err != nil {
	// 	log.Fatal(err)
	// }
}


