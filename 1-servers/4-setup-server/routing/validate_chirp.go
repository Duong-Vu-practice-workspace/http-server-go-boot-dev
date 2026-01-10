package routing

import (
	"errors"
	"net/http"
	"strings"
)

type ValidateChirp struct {
	Body string `json:"body"`
}
type ValidateChirpError struct {
	Error string `json:"error"`
}

var filterWordSet = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

const cleanedString = "****"

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request, chirp ValidateChirp) (string, error) {
	if len(chirp.Body) > 140 {
		msg := "chirp is too long"
		RespondWithError(w, http.StatusBadRequest, msg)
		return "", errors.New(msg)
	}
	words := strings.Split(chirp.Body, " ")
	result := strings.Join(filterWords(words), " ")
	return result, nil
}
func filterWords(words []string) []string {
	for i := range words {
		_, ok := filterWordSet[strings.ToLower(words[i])]
		if ok {
			words[i] = cleanedString
		}
	}
	return words
}
