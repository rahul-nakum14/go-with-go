package utils


func parseJson(r *http.Request, payload interface{}) error {
	
	if r.Body == nil {
		return fmt.Errorf("empty body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func writeJson(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload) 
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	return writeJson(w, status, map[string]string{"error": err.Error()})
}