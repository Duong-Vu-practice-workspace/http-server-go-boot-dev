package chirp

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ValidateChirp struct {
	Body string `json:"body"`
}
type ValidateChirpError struct {
	Error string `json:"error"`
}
type ValidateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

var filterWordSet = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

const cleanedString = "****"

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := ValidateChirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := strings.ToLower(chirp.Body)
	words := strings.Split(data, " ")
	response := ValidateChirpResponse{CleanedBody: strings.Join(filterWords(words), " ")}
	_ = json.NewEncoder(w).Encode(response)
}
func filterWords(words []string) []string {
	for i := range words {
		_, ok := filterWordSet[words[i]]
		if ok {
			words[i] = cleanedString
		}
	}
	return words
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, ValidateChirpError{Error: msg})
}
