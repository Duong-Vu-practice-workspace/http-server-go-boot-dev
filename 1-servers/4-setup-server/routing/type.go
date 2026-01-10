package routing

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"example.com/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	Db             *sql.DB
	Queries        *database.Queries
	Platform       string
	JwtSecret      string
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, ValidateChirpError{Error: msg})
}
