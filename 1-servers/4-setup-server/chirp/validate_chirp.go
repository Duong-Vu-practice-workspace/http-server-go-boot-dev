package chirp

import (
	"encoding/json"
	"net/http"
)

type ValidateChirp struct {
	Body string `json:"body"`
}
type ValidateChirpError struct {
	Error string `json:"error"`
}
type ValidateChirpResponse struct {
	Valid bool `json:"valid"`
}

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := ValidateChirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		response := ValidateChirpError{Error: "Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	if len(chirp.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := ValidateChirpError{Error: "Chirp is too long"}
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ValidateChirpResponse{Valid: true}
	_ = json.NewEncoder(w).Encode(response)
}
