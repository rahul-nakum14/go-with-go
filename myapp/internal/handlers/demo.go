package handlers

import "net/http"

func DemoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the demo endpoint"))
}